package adcpi

import (
	"fmt"
	"os"
)

func setBit(_byte *byte, bit uint, value bool) {

	//fmt.Printf("byte: %x, bit: %d, value: %t\n", *_byte, bit, value)

	if value {
		*_byte |= (1 << bit)
	} else {
		*_byte &^= (1 << bit)
	}
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func isBitSet(_byte byte, bit uint) (check bool) {
	check = (_byte & (1 << bit)) != 0
	return check
}
