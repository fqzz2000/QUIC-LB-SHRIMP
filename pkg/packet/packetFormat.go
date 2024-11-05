package packet

type PacketType uint8

const (
	Initial   PacketType = 0x00
	ZeroRTT   PacketType = 0x01
	HandShake PacketType = 0x02
	Retry     PacketType = 0x03
	OneRTT    PacketType = 0x04
)

type QuicHeader interface {
	GetCID() ([]byte, error)
	GetPacketType() (PacketType, error)
	GetHeaderForm() (uint8, error)
}

type LongHeader struct {
	HeaderForm     uint8
	LongPacketType PacketType
	TypeSpecific   uint8
	Version        uint32
	DCIDLength     uint8
	DCID           []byte
	SCIDLength     uint8
	SCID           []byte
}

type ShortHeader struct {
	HeaderForm         uint8
	ReservedBits       uint8
	KeyPhase           uint8
	PacketNumberLength uint8
	DCID               []byte
	PacketNumber       uint64
}

func (lh *LongHeader) GetCID() ([]byte, error) {
	return lh.DCID, nil
}

func (lh *LongHeader) GetPacketType() (PacketType, error) {
	return lh.LongPacketType, nil
}

func (lh *LongHeader) GetHeaderForm() (uint8, error) {
	return lh.HeaderForm, nil
}

func (sh *ShortHeader) GetCID() ([]byte, error) {
	return sh.DCID, nil
}

func (sh *ShortHeader) GetPacketType() (PacketType, error) {
	return OneRTT, nil
}

func (sh *ShortHeader) GetHeaderForm() (uint8, error) {
	return sh.HeaderForm, nil
}
