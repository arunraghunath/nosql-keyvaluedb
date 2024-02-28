package main

import "encoding/binary"

const metaPgnum = 0

type meta struct {
	freespacePgnum pagenum
}

func newMeta() *meta {
	return &meta{}
}

func (m *meta) serialize(buf []byte) {
	pos := 0
	binary.BigEndian.PutUint64(buf[pos:], uint64(m.freespacePgnum))
	pos += 1
}

func (m *meta) deserialize(buf []byte) {
	pos := 0
	m.freespacePgnum = pagenum(binary.BigEndian.Uint64(buf[pos:]))
}
