package main

import (
	"fmt"
	"os"

	"github.com/soyum2222/editPE"
)

func main() {
    // open a exe file
	file, err := os.ReadFile("..\\hello.exe")
	if err != nil {
		panic(err)
	}

	p := editPE.PE{}
	p.Parse(file)

	fmt.Println(len(p.ImageSectionHeaders))

	// add a new section size is 100 byte
	p.AddSection(".new", 100)

	// save to disk
	f, err := os.Create("..\\hello2.exe")
	if err != nil {
		panic(err)
	}

	f.Write(p.Raw)
}
