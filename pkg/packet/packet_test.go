package packet

import (
	"testing"
)

func TestParseLongHeader(t *testing.T) {
	tests := []struct {
		name     string
		packet   []byte
		expected *LongHeader
	}{
		{
			name: "Valid Initial Packet",
			packet: []byte{
				0xC0,                   // Header Form (1) + Fixed Bit (1) + Long Packet Type (00) + Type Specific (0000)
				0x00, 0x00, 0x00, 0x01, // Version
				0x08,                   // DCID Length
				0x01, 0x02, 0x03, 0x04, // DCID (8 bytes)
				0x05, 0x06, 0x07, 0x08,
				0x04,             // SCID Length
				0x0A, 0x0B, 0x0C, // SCID (4 bytes)
				0x0D,
			},
			expected: &LongHeader{
				HeaderForm:     1,
				LongPacketType: Initial,
				TypeSpecific:   0,
				Version:        1,
				DCIDLength:     8,
				DCID:           []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
				SCIDLength:     4,
				SCID:           []byte{0x0A, 0x0B, 0x0C, 0x0D},
			},
		},
	}

	processor := &PacketProcessor{DCIDLength: 8}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			header, err := processor.parseLongHeader(tt.packet)
			if err != nil {
				t.Errorf("parseLongHeader() error = %v", err)
				return
			}

			if header.HeaderForm != tt.expected.HeaderForm {
				t.Errorf("HeaderForm = %v, want %v", header.HeaderForm, tt.expected.HeaderForm)
			}
			if header.LongPacketType != tt.expected.LongPacketType {
				t.Errorf("LongPacketType = %v, want %v", header.LongPacketType, tt.expected.LongPacketType)
			}
			if header.Version != tt.expected.Version {
				t.Errorf("Version = %v, want %v", header.Version, tt.expected.Version)
			}
		})
	}
}

func TestParseShortHeader(t *testing.T) {
	tests := []struct {
		name     string
		packet   []byte
		expected *ShortHeader
	}{
		{
			name: "Valid Short Header Packet",
			packet: []byte{
				0x40,                   // Header Form (0) + Fixed Bit (1) + Reserved (00) + Key Phase (0) + Packet Number Length (00)
				0x01, 0x02, 0x03, 0x04, // DCID (8 bytes)
				0x05, 0x06, 0x07, 0x08,
			},
			expected: &ShortHeader{
				HeaderForm:         0,
				ReservedBits:       0,
				KeyPhase:           0,
				PacketNumberLength: 0,
				DCID:               []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
			},
		},
	}

	processor := &PacketProcessor{DCIDLength: 8}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			header, err := processor.parseShortHeader(tt.packet)
			if err != nil {
				t.Errorf("parseShortHeader() error = %v", err)
				return
			}

			if header.HeaderForm != tt.expected.HeaderForm {
				t.Errorf("HeaderForm = %v, want %v", header.HeaderForm, tt.expected.HeaderForm)
			}
			if header.ReservedBits != tt.expected.ReservedBits {
				t.Errorf("ReservedBits = %v, want %v", header.ReservedBits, tt.expected.ReservedBits)
			}
			if header.KeyPhase != tt.expected.KeyPhase {
				t.Errorf("KeyPhase = %v, want %v", header.KeyPhase, tt.expected.KeyPhase)
			}
		})
	}
}
