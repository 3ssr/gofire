package gofire

import (
	"bytes"
	"encoding/binary"
	"gofire/iface"
)

const HeaderLength uint32 = 8

type FireMsg struct {
	ID            string
	PayloadLength uint32
	ActionId      uint32
	Payload       []byte
}

func NewMsg() iface.IMsg {
	m := &FireMsg{}
	return m
}

func (m *FireMsg) SetPayload(payload []byte) {
	m.Payload = payload
}

func (m *FireMsg) GetPayloadLen() uint32 {
	return m.PayloadLength
}

func (m *FireMsg) GetAction() uint32 {
	return m.ActionId
}

func (h *FireMsg) UnpackHead(data []byte) error {
	r := bytes.NewReader(data)
	err := binary.Read(r, binary.LittleEndian, &h.ID)
	if err != nil {
		return err
	}

	err = binary.Read(r, binary.LittleEndian, &h.PayloadLength)
	if err != nil {
		return err
	}

	err = binary.Read(r, binary.LittleEndian, &h.ActionId)
	if err != nil {
		return err
	}

	return nil
}

func (m *FireMsg) Pack() ([]byte, error) {
	buffer := bytes.NewBuffer([]byte{})
	err := binary.Write(buffer, binary.LittleEndian, &m.ID)
	if err != nil {
		return nil, err
	}

	err = binary.Write(buffer, binary.LittleEndian, &m.PayloadLength)
	if err != nil {
		return nil, err
	}

	err = binary.Write(buffer, binary.LittleEndian, &m.ActionId)
	if err != nil {
		return nil, err
	}

	err = binary.Write(buffer, binary.LittleEndian, &m.Payload)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
