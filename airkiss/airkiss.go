package airkiss

import (
	"github.com/sigurn/crc8"
)

type AirKiss struct {
	src    []uint16
	essid  string
	passwd string
	random uint8
}

func MkInt16(H uint8, L uint8) uint16 {
	return (uint16(H) << 8) | (uint16(L) & 0xff)
}

func MkInt8(H uint8, L uint8) uint8 {
	return (H << 4) | (L & 0x0f)
}

func HInt8(data uint8) uint8 {
	return (data & 0xf0) >> 4
}

func LInt8(data uint8) uint8 {
	return data & 0xf
}

func New(essid string, passwd string, random uint8) *AirKiss {
	return &AirKiss{
		essid:  essid,
		passwd: passwd,
		random: random,
	}
}

func (airkiss *AirKiss) makeGuide() {
	guide := []uint16{1, 2, 3, 4}

	airkiss.src = append(airkiss.src, guide...)
}

func (airkiss *AirKiss) magicCode(dataLen uint8) {
	table := crc8.MakeTable(crc8.CRC8_MAXIM)
	crc := crc8.Checksum([]byte(airkiss.essid), table)

	code := make([]uint16, 4)
	code[0] = MkInt16(0, MkInt8(0, HInt8(dataLen)))
	code[1] = MkInt16(0, MkInt8(1, LInt8(dataLen)))
	code[2] = MkInt16(0, MkInt8(2, HInt8(crc)))
	code[3] = MkInt16(0, MkInt8(3, LInt8(crc)))

	airkiss.src = append(airkiss.src, code...)
}

func (airkiss *AirKiss) prefixCode(passwdLen uint8) {
	table := crc8.MakeTable(crc8.CRC8_MAXIM)
	crc := crc8.Checksum([]byte{passwdLen}, table)

	code := make([]uint16, 4)
	code[0] = MkInt16(0, MkInt8(4, HInt8(passwdLen)))
	code[1] = MkInt16(0, MkInt8(5, LInt8(passwdLen)))
	code[2] = MkInt16(0, MkInt8(6, HInt8(crc)))
	code[3] = MkInt16(0, MkInt8(7, LInt8(crc)))

	airkiss.src = append(airkiss.src, code...)
}

func (airkiss *AirKiss) sequenceCode(data []uint8, packetLen uint8) {
	var index uint8
	var dataIndex uint8

	for index = 0; index < packetLen; index++ {
		table := crc8.MakeTable(crc8.CRC8_MAXIM)
		newData := append([]byte{index}, data[dataIndex:dataIndex+4]...)
		crc := crc8.Checksum(newData, table)

		code := make([]uint16, 6)
		code[0] = MkInt16(0, (0x80 | (crc & 0x7f)))
		code[1] = MkInt16(0, (0x80 | index))
		for i, da := range data[dataIndex : dataIndex+4] {
			code[i+2] = MkInt16(1, da)
		}
		dataIndex += 4
		airkiss.src = append(airkiss.src, code...)
	}
}

func (airkiss *AirKiss) GreateCodePackage() []uint16 {

	sequenCode := append([]uint8(airkiss.passwd), airkiss.random)
	sequenCode = append(sequenCode, []uint8(airkiss.essid)...)

	airkiss.makeGuide()

	for i := 0; i < 5; i++ {
		airkiss.magicCode(uint8(len(sequenCode)))
	}

	airkiss.prefixCode(uint8(len(airkiss.passwd)))

	packetLen := uint8(len(sequenCode) / 4)
	if len(sequenCode)%4 > 0 {
		packetLen++
	}

	airkiss.sequenceCode(sequenCode, packetLen)

	return airkiss.src
}
