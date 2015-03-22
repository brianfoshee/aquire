package main

import (
_	"bytes"
_	"encoding/binary"
	"fmt"
_	"strconv"
	"time"
	"github.com/davecheney/i2c"
)

	// Info is 0x49
	// Reading is 0x52

func main() {
	i, err := i2c.New(0x63, 1)
	if err != nil {
		fmt.Println("Error opening device")
	}
	defer i.Close()

	n, err := i.WriteByte(0x52)
	//status := []byte{0x53, 0x54, 0x41, 0x54, 0x55, 0x53}
	//sleep := []byte{0x53, 0x4C, 0x45, 0x45, 0x50}
	//n, err := i.Write(sleep)
	if err != nil {
		fmt.Println("Error writing byte")
	}
	fmt.Println("Wrote bytes", n)
	time.Sleep(time.Second * 1)

	buf := make([]byte, 20)
	r, err := i.Read(buf)
	if err != nil {
		fmt.Println("Error reading bytes")
	}
	fmt.Printf("Read %v bytes: %v\n", r, buf)

	fmt.Println(string(buf))

	//b := bytes.NewBuffer(buf)
	//var o int
	//binary.Read(b, binary.LittleEndian, &o)	
	//fmt.Println("Slice: ", o)
	//for i := 0; i < len(o); i++ {
	//	fmt.Println("Read ", strconv.Itoa(o[i]))
	//}
}
