package traffic

import (
	"sort"

	"github.com/gopacket/gopacket/layers"
)

type tcpSegment struct {
	seq  uint32
	data []byte
}
type tcpStream struct {
	segments []tcpSegment
	evidence Evidence
	bytes    int
	syn      bool
	initial  uint32
}
type streamCollector struct {
	streams    map[string]*tcpStream
	buffered   int
	consume    func([]byte, Evidence)
	warn       func(string)
	prefixOnly bool
}

func newStreamCollector(consume func([]byte, Evidence), warn func(string)) *streamCollector {
	return &streamCollector{streams: make(map[string]*tcpStream), consume: consume, warn: warn}
}

func (s *streamCollector) add(tcp *layers.TCP, e Evidence) {
	key := e.SrcIP + ":" + e.SrcPort + ">" + e.DstIP + ":" + e.DstPort
	stream := s.streams[key]
	if stream != nil && tcp.SYN && (!stream.syn || stream.initial != tcp.Seq) {
		s.emit(stream)
		delete(s.streams, key)
		stream = nil
	}
	if stream == nil {
		if len(s.streams) >= maxStreams {
			s.warn("stream_count_limit")
			return
		}
		stream = &tcpStream{evidence: e, syn: tcp.SYN, initial: tcp.Seq}
		s.streams[key] = stream
	}
	if e.Time.Before(stream.evidence.Time) {
		stream.evidence.Time = e.Time
	}
	if e.Time.After(stream.evidence.EndTime) {
		stream.evidence.EndTime = e.Time
	}
	if len(tcp.Payload) > 0 {
		if s.prefixOnly {
			// PROXY headers are at the beginning of a connection. Do not retain
			// the tunneled attack traffic merely to discover its client address.
			remaining := 4096 - stream.bytes
			if remaining <= 0 {
				return
			}
			copyTCP := *tcp
			copyTCP.Payload = tcp.Payload[:min(len(tcp.Payload), remaining)]
			tcp = &copyTCP
		}
		if stream.bytes+len(tcp.Payload) > maxStreamBytes || s.buffered+len(tcp.Payload) > maxBufferedBytes || len(stream.segments) >= 8192 {
			s.warn("stream_buffer_limit")
		} else {
			seq := tcp.Seq
			if tcp.SYN {
				seq++
			}
			stream.segments = append(stream.segments, tcpSegment{seq: seq, data: append([]byte{}, tcp.Payload...)})
			stream.bytes += len(tcp.Payload)
			s.buffered += len(tcp.Payload)
		}
	}
	// Delay FIN/RST flushing until EOF to allow late out-of-order segments.
}

func (s *streamCollector) emit(stream *tcpStream) {
	s.buffered -= stream.bytes
	if len(stream.segments) == 0 {
		return
	}
	base := stream.segments[0].seq
	sort.SliceStable(stream.segments, func(i, j int) bool { return int32(stream.segments[i].seq-base) < int32(stream.segments[j].seq-base) })
	var data []byte
	next := stream.segments[0].seq
	if stream.syn && next != stream.initial+1 {
		s.warn("tcp_capture_gap")
	}
	for _, segment := range stream.segments {
		delta := int32(segment.seq - next)
		if delta > 0 {
			// Never join across a capture gap: doing so can fabricate flags.
			s.consume(data, stream.evidence)
			data = nil
			s.warn("tcp_capture_gap")
			next = segment.seq
		} else if delta < 0 {
			skip := int(-int64(delta))
			if skip >= len(segment.data) {
				continue
			}
			segment.data = segment.data[skip:]
		}
		data = append(data, segment.data...)
		next += uint32(len(segment.data))
	}
	if len(data) > 0 {
		s.consume(data, stream.evidence)
	}
}

func (s *streamCollector) flush() {
	keys := make([]string, 0, len(s.streams))
	for key := range s.streams {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		s.emit(s.streams[key])
		delete(s.streams, key)
	}
}
