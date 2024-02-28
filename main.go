package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Printf("OS page size is %T", os.Getpagesize())
	al, _ := newAccess("nosql.db", os.Getpagesize())
	fmt.Println("Creating a page")
	p := al.newBlankPage()
	p.num = al.nextPage()
	copy(p.data, "data-3")
	al.writePage(p)
	al.writeFreespace()
	al.close()

	al, _ = newAccess("nosql.db", os.Getpagesize())
	fmt.Println("Creating a page")
	p = al.newBlankPage()
	p.num = al.nextPage()
	copy(p.data, "data-4")
	al.writePage(p)
	al.writeFreespace()
	al.close()
}
