package packet

import "encoding/binary"

type PacketProcessor struct {
	DCIDLength uint8 // TODO: it suppose to be a map of connection unique id to length
}

type HeaderParser interface {
	ParsePacket(packet []byte) (*QuicHeader, error)

	ClassifyPacket(packet []byte) (PacketType, error)

	ExtractCID(packet []byte) ([]byte, error)
}

type Validator interface {
	ValidatePacket(packet []byte) error
}

func assert(condition bool) {
	if !condition {
		panic("assertion failed")
	}
}

// first implement parse packet
func (p *PacketProcessor) ParsePacket(packet []byte) (QuicHeader, error) {
	// get first byte of packet
	firstByte := packet[0]
	headerForm := firstByte >> 7

	assert(headerForm&0x1 == 1 || headerForm&0x1 == 0)

	if headerForm == 1 {
		// long header
		return p.parseLongHeader(packet)
	} else {
		// short header
		return p.parseShortHeader(packet)
	}
}

func (p *PacketProcessor) parseLongHeader(packet []byte) (*LongHeader, error) {
	header := &LongHeader{}
	header.HeaderForm = 1
	header.LongPacketType = PacketType((packet[0] >> 4) & 0x1)
	header.TypeSpecific = packet[0] & 0x0F
	header.Version = binary.BigEndian.Uint32(packet[1:5])
	header.DCIDLength = packet[5] // DCID length report length in byte
	header.DCID = packet[6 : 6+header.DCIDLength]
	header.SCIDLength = packet[6+header.DCIDLength] // SCID length report length in byte
	header.SCID = packet[6+header.DCIDLength+1 : 6+header.DCIDLength+1+header.SCIDLength]
	return header, nil
}

func (p *PacketProcessor) parseShortHeader(packet []byte) (*ShortHeader, error) {
	header := &ShortHeader{}
	header.HeaderForm = 0
	header.ReservedBits = (packet[0] >> 3) & 0x3
	header.KeyPhase = (packet[0] >> 2) & 0x1
	header.PacketNumberLength = packet[0] & 0x3
	header.DCID = packet[1 : 1+p.DCIDLength] // length of DCID is expected to known by LB
	return header, nil
}
