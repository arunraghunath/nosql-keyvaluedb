package main

import (
	"fmt"
	"os"
)

type pagenum uint

type page struct {
	num  pagenum
	data []byte
}

type accesslayer struct {
	file     *os.File
	pagesize int
	*meta
	*freespace
}

func newAccess(path string, size int) (*accesslayer, error) {
	al := &accesslayer{meta: newMeta(), pagesize: pagesize, freespace: newFreespace()}
	if _, err := os.Stat(path); err == nil {
		al.file, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			return nil, err
		}
		err = al.readMeta()
		fmt.Println("meta info is", al.meta.freespacePgnum)
		if err != nil {
			return nil, err
		}
		err = al.readFreespace()
		if err != nil {
			return nil, err
		}
	} else {
		al.file, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			_ = al.close()
			return nil, err
		}
		al.meta.freespacePgnum = al.nextPage()
		_ = al.writeFreespace()
		_ = al.writeMeta()

	}
	return al, nil
}

func (al *accesslayer) close() error {
	if al.file != nil {
		err := al.file.Close()
		if err != nil {
			return fmt.Errorf("could not close the file : %s", err)
		}
		al.file = nil
	}
	return nil
}

func (al *accesslayer) newBlankPage() *page {
	return &page{data: make([]byte, al.pagesize, al.pagesize)}
}

func (al *accesslayer) writePage(p *page) error {
	offset := p.num * pagenum(al.pagesize)
	_, err := al.file.WriteAt(p.data, int64(offset))
	return err
}

func (al *accesslayer) readPage(pgnum pagenum) (*page, error) {
	offset := pgnum * pagenum(al.pagesize)
	p := al.newBlankPage()
	n, err := al.file.ReadAt(p.data, int64(offset))
	fmt.Println("Read content count", n)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (al *accesslayer) writeFreespace() error {
	p := al.newBlankPage()
	al.freespace.serialize(p.data)
	p.num = al.meta.freespacePgnum
	err := al.writePage(p)
	return err
}

func (al *accesslayer) readFreespace() error {
	p, err := al.readPage(al.meta.freespacePgnum)
	if err != nil {
		return err
	}
	al.freespace.deserialize(p.data)
	return nil
}

func (al *accesslayer) writeMeta() error {
	p := al.newBlankPage()
	p.num = metaPgnum
	al.meta.serialize(p.data)
	err := al.writePage(p)
	return err
}

func (al *accesslayer) readMeta() error {
	p, err := al.readPage(metaPgnum)
	if err != nil {
		return err
	}
	al.meta.deserialize(p.data)
	return nil
}
