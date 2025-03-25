package utils

import (
	"CBCTF/internel/log"
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/gopacket/gopacket"
	"github.com/gopacket/gopacket/layers"
	"github.com/gopacket/gopacket/pcap"
	"net/http"
	"os"
)

const (
	Unknown  = "unknown"
	Request  = "request"
	Response = "response"
)

type Connection struct {
	SrcIP   string
	DstIP   string
	SrcPort uint16
	DstPort uint16
	Payload []byte
}

func (conn Connection) ParsePayload() (interface{}, string) {
	reader := bytes.NewReader(conn.Payload)

	req, err := http.ReadRequest(bufio.NewReader(reader))
	if err == nil {
		return req, Request
	}

	reader = bytes.NewReader(conn.Payload)
	resp, err := http.ReadResponse(bufio.NewReader(reader), nil)
	if err == nil {
		return resp, Response
	}
	return conn.Payload, Unknown
}

func ReadPcap(path string) ([]Connection, bool, string) {
	if _, err := os.Stat(path); err != nil {
		log.Logger.Warningf("Failed to get file: %s", err)
		if errors.Is(err, os.ErrNotExist) {
			return []Connection{}, false, "PcapNotFound"
		}
		return []Connection{}, false, "UnknownError"
	}
	handle, err := pcap.OpenOffline(path)
	if err != nil {
		log.Logger.Warningf("Failed to read .pcap: %s", err)
		return []Connection{}, false, "ReadPcapError"
	}
	defer handle.Close()
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	var connections []Connection
	tmp := make(map[string]*Connection)
	for packet := range packetSource.Packets() {
		network := packet.NetworkLayer()
		if network == nil {
			continue
		}
		ipv4, ok := network.(*layers.IPv4)
		if !ok {
			continue
		}
		transport := packet.TransportLayer()
		if transport == nil {
			continue
		}
		tcp, ok := transport.(*layers.TCP)
		if ok {
			connID := fmt.Sprintf("%s:%d-%s:%d", ipv4.SrcIP.String(), tcp.SrcPort, ipv4.DstIP.String(), tcp.DstPort)
			if _, exists := tmp[connID]; !exists {
				tmp[connID] = &Connection{
					SrcIP:   ipv4.SrcIP.String(),
					DstIP:   ipv4.DstIP.String(),
					SrcPort: uint16(tcp.SrcPort),
					DstPort: uint16(tcp.DstPort),
				}
			}
			tmp[connID].Payload = append(tmp[connID].Payload, tcp.Payload...)
			if tcp.FIN {
				connections = append(connections, *tmp[connID])
				delete(tmp, connID)
			}
		}
	}
	return connections, true, "Success"
}
