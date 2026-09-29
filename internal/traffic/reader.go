package traffic

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/gopacket/gopacket"
	"github.com/gopacket/gopacket/layers"
	"github.com/gopacket/gopacket/pcapgo"
)

// Read synchronously: PacketSource.Packets hides read errors and can leak its
// producer goroutine when a cancelled consumer stops reading.
func walkTrafficPackets(ctx context.Context, path string, visit func(gopacket.Packet, layers.LinkType) error) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	r := bufio.NewReader(f)
	magic, err := r.Peek(4)
	if err != nil {
		if err == io.EOF {
			// Capture sidecars can create a file before the first packet arrives.
			// Only a genuinely empty file is an empty capture, not a partial header.
			if len(magic) == 0 {
				return ctx.Err()
			}
			return io.ErrUnexpectedEOF
		}
		return err
	}
	var source gopacket.PacketDataSource
	var link layers.LinkType
	mixed := bytes.Equal(magic, []byte{10, 13, 13, 10})
	if mixed {
		source, err = pcapgo.NewNgReader(r, pcapgo.NgReaderOptions{WantMixedLinkType: true})
	} else {
		var reader *pcapgo.Reader
		reader, err = pcapgo.NewReader(r)
		if err == nil {
			source, link = reader, reader.LinkType()
		}
	}
	if err != nil {
		return err
	}
	for {
		if err = ctx.Err(); err != nil {
			return err
		}
		data, ci, readErr := source.ReadPacketData()
		if readErr == io.EOF {
			if ci.CaptureLength > 0 || len(data) > 0 {
				return fmt.Errorf("read %s: %w", path, io.ErrUnexpectedEOF)
			}
			return nil
		}
		if readErr != nil {
			return fmt.Errorf("read %s: %w", path, readErr)
		}
		if mixed {
			link = ci.AncillaryData[0].(layers.LinkType)
		}
		packet := gopacket.NewPacket(data, link, gopacket.Default)
		packet.Metadata().CaptureInfo = ci
		if err = visit(packet, link); err != nil {
			return err
		}
	}
}

func isOriginalTrafficCapture(name string) bool {
	name = strings.ToLower(name)
	return (strings.HasSuffix(name, ".pcap") || strings.HasSuffix(name, ".pcapng")) && !strings.Contains(name, ".enrich.")
}
