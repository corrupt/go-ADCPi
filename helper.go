package adcpi

import (
	"fmt"
	"os"
	"syscall"
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

func valueUpdated(rate Samplerate, buf []byte) bool {
	var config byte
	switch rate {
	case SR12, SR14, SR16:
		config = buf[2]
	case SR18:
		config = buf[3]
	}
	return !isBitSet(config, 7)
}

func ioctl(fd, cmd, arg uintptr) (err error) {
	_, _, e1 := syscall.Syscall6(syscall.SYS_IOCTL, fd, cmd, arg, 0, 0, 0)
	if e1 != 0 {
		err = e1
	}
	return
}

func interpretValue(rate Samplerate, buf []byte) int32 {
	var lower, middle, upper byte
	var value int32 = 0

	switch rate {
	case SR12, SR14, SR16:
		upper = buf[0]
		lower = buf[1]
		//config = buf[2]
	case SR18:
		upper = buf[0]
		middle = buf[1]
		lower = buf[2]
		//config = buf[3]
	}
	switch rate {
	case SR12:
		//fmt.Println("SR12")
		value = (int32(upper&0x0F) << 8) | int32(lower)
		if isBitSet(upper, 3) {
			//fmt.Printf("Value: %d, Mask: %x\n", value, 0x000007FF)
			value &= 0x000007FF //clear the sign bit
			//fmt.Printf("Value: %d, Mask: %x\n", value, 0x000007FF)
			value *= -1
		}
		return value
	case SR14:
		//fmt.Println("SR14")
		value = (int32(upper&0x3F) << 8) | int32(lower)
		if isBitSet(upper, 5) {
			value &= 0x00001FFF
			value *= -1
		}
		return value
	case SR16:
		//fmt.Println("SR16")
		value = (int32(upper) << 8) | int32(lower)
		if isBitSet(upper, 7) {
			value &= 0x00007FFF
			value *= -1
		}
		return value
	case SR18:
		//fmt.Println("SR18")
		value = (int32(upper&0x03) << 16) | (int32(middle) << 8) | int32(lower)
		if isBitSet(upper, 1) {
			value &= 0x0001FFFF
			value *= -1
		}
		return value
	}
	return 0
}
