package main

import (
	"encoding/binary"
	"fmt"
)

type freespace struct {
	lastpage   pagenum
	freedPages []pagenum
}

func newFreespace() *freespace {
	return &freespace{lastpage: metaPgnum, freedPages: []pagenum{}}
}

func (fs *freespace) nextPage() pagenum {
	if len(fs.freedPages) != 0 {
		retpg := fs.freedPages[0]
		fs.freedPages = fs.freedPages[1:]
		return retpg
	}
	fs.lastpage += 1
	fmt.Println("New page is", fs.lastpage)
	return fs.lastpage
}

func (fs *freespace) serialize(buf []byte) {
	pos := 0
	binary.BigEndian.PutUint64(buf[pos:], uint64(fs.lastpage))
	pos += 8
	binary.BigEndian.PutUint64(buf[pos:], uint64(len(fs.freedPages)))
	pos += 8
	for _, page := range fs.freedPages {
		binary.BigEndian.PutUint64(buf[pos:], uint64(page))
		pos += pagenumsize
	}
}

func (fs *freespace) deserialize(buf []byte) {
	pos := 0
	fs.lastpage = pagenum(binary.BigEndian.Uint64(buf[pos:]))
	pos += 8
	freedPagesCount := binary.BigEndian.Uint64(buf[pos:])
	pos += 8
	for i := 0; i < int(freedPagesCount); i++ {
		fs.freedPages = append(fs.freedPages, pagenum(binary.BigEndian.Uint64(buf[pos:])))
		pos += pagenumsize
	}
}
