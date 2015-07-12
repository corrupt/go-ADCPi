package main

import (
	"errors"
	"fmt"
	"os"
	"syscall"
)

const (
	i2c_SLAVE              = 0x0703
	BUFSIZE                = 25
	address_default        = 0xD0 //0x68 and R/W=0 (bit is "read/not write, 1 = read, 0 = write)"
	configbyte_default     = 0x00
	bitrate_default        = 18
	conversionmode_default = Continuous
	samplerate_default     = SR12
	channel_default        = Ch1
	gain_default           = X1
)

type Samplerate uint8
type Channel uint8
type Address uint8
type RW bool
type Conversionmode bool
type Gain uint8

type I2C struct {
	rc             *os.File
	address        Address
	channel        Channel
	rwmode         RW
	conversionmode Conversionmode
	samplerate     Samplerate
	gain           Gain
	addressbyte    byte
	configbyte     byte
}

const (
	Ch1 Channel = 1
	Ch2 Channel = 2
	Ch3 Channel = 3
	Ch4 Channel = 4
)

const (
	SR12 Samplerate = 12
	SR14 Samplerate = 14
	SR16 Samplerate = 16
	SR18 Samplerate = 18
)

const (
	A68 Address = 0x68
	A69 Address = 0x69
	A6A Address = 0x6A
	A6B Address = 0x6B
	A6C Address = 0x6C
	A6D Address = 0x6D
	A6E Address = 0x6E
	A6F Address = 0x6F
)

const (
	Read  RW = true
	Write RW = false
)

const (
	Continuous Conversionmode = true
	OneShot    Conversionmode = false
)

const (
	X1 Gain = 1 << iota
	X2
	X4
	X8
)

func setBit(_byte *byte, bit uint, value bool) {

	//fmt.Printf("byte: %x, bit: %d, value: %t\n", *_byte, bit, value)

	if value {
		*_byte |= (1 << bit)
	} else {
		*_byte &^= (1 << bit)
	}
}

func (i2c *I2C) SetAddress(address Address) error {
	switch address {
	case A68:
		i2c.address = address
		setBit(&(i2c.addressbyte), 1, false)
		setBit(&(i2c.addressbyte), 2, false)
		setBit(&(i2c.addressbyte), 3, false)
		return nil
	case A69:
		i2c.address = address
		setBit(&(i2c.addressbyte), 1, true)
		setBit(&(i2c.addressbyte), 2, false)
		setBit(&(i2c.addressbyte), 3, false)
		return nil
	case A6A:
		i2c.address = address
		setBit(&(i2c.addressbyte), 1, false)
		setBit(&(i2c.addressbyte), 2, true)
		setBit(&(i2c.addressbyte), 3, false)
		return nil
	case A6B:
		i2c.address = address
		setBit(&(i2c.addressbyte), 1, true)
		setBit(&(i2c.addressbyte), 2, true)
		setBit(&(i2c.addressbyte), 3, false)
		return nil
	case A6C:
		i2c.address = address
		setBit(&(i2c.addressbyte), 1, false)
		setBit(&(i2c.addressbyte), 2, false)
		setBit(&(i2c.addressbyte), 3, true)
		return nil
	case A6D:
		i2c.address = address
		setBit(&(i2c.addressbyte), 1, true)
		setBit(&(i2c.addressbyte), 2, false)
		setBit(&(i2c.addressbyte), 3, true)
		return nil
	case A6E:
		i2c.address = address
		setBit(&(i2c.addressbyte), 1, false)
		setBit(&(i2c.addressbyte), 2, true)
		setBit(&(i2c.addressbyte), 3, true)
		return nil
	case A6F:
		i2c.address = address
		setBit(&(i2c.addressbyte), 1, true)
		setBit(&(i2c.addressbyte), 2, true)
		setBit(&(i2c.addressbyte), 3, true)
		return nil
	default:
		return errors.New(fmt.Sprintf("Invalid Address: %x. Allowed values: %x, %x, %x, %x, %x, %x, %x, %x", address, A68, A69, A6A, A6B, A6C, A6D, A6E, A6F))
	}
}

func (i2c *I2C) setRW(state RW) {
	setBit(&(i2c.configbyte), 0, bool(state))
}

func (i2c *I2C) SetSamplerate(rate Samplerate) error {
	switch rate {
	case SR12:
		i2c.samplerate = rate
		setBit(&(i2c.configbyte), 2, false)
		setBit(&(i2c.configbyte), 3, false)
		return nil
	case SR14:
		i2c.samplerate = rate
		setBit(&(i2c.configbyte), 2, true)
		setBit(&(i2c.configbyte), 3, false)
		return nil
	case SR16:
		i2c.samplerate = rate
		setBit(&(i2c.configbyte), 2, false)
		setBit(&(i2c.configbyte), 3, true)
		return nil
	case SR18:
		i2c.samplerate = rate
		setBit(&(i2c.configbyte), 2, true)
		setBit(&(i2c.configbyte), 3, true)
		return nil
	default:
		return errors.New(fmt.Sprintf("Invalid conversion rate: %d. Allowed values: %d, %d, %d, %d.", rate, SR12, SR14, SR16, SR18))
	}
}

func (i2c *I2C) SetChannel(channel Channel) error {
	switch channel {
	case Ch1:
		i2c.channel = channel
		setBit(&(i2c.configbyte), 5, false)
		setBit(&(i2c.configbyte), 6, false)
		return nil
	case Ch2:
		i2c.channel = channel
		setBit(&(i2c.configbyte), 5, true)
		setBit(&(i2c.configbyte), 6, false)
		return nil
	case Ch3:
		i2c.channel = channel
		setBit(&(i2c.configbyte), 5, false)
		setBit(&(i2c.configbyte), 6, true)
		return nil
	case Ch4:
		i2c.channel = channel
		setBit(&(i2c.configbyte), 5, true)
		setBit(&(i2c.configbyte), 6, true)
		return nil
	default:
		return errors.New(fmt.Sprintf("Invalid channel: %d. Allowed values: %d, %d, %d, %d", channel, Ch1, Ch2, Ch3, Ch4))
	}
}

func (i2c *I2C) SetConversionmode(mode Conversionmode) {
	setBit(&(i2c.configbyte), 4, bool(mode))
}

func (i2c *I2C) SetGain(gain Gain) error {
	switch gain {
	case X1:
		i2c.gain = gain
		setBit(&(i2c.configbyte), 0, false)
		setBit(&(i2c.configbyte), 1, false)
		return nil
	case X2:
		i2c.gain = gain
		setBit(&(i2c.configbyte), 0, true)
		setBit(&(i2c.configbyte), 1, false)
		return nil
	case X4:
		i2c.gain = gain
		setBit(&(i2c.configbyte), 0, false)
		setBit(&(i2c.configbyte), 1, true)
		return nil
	case X8:
		i2c.gain = gain
		setBit(&(i2c.configbyte), 0, true)
		setBit(&(i2c.configbyte), 1, true)
		return nil
	default:
		return errors.New(fmt.Sprintf("Invalid gain: %d. Allowed values: %d, %d, %d, %d", X1, X2, X4, X8))
	}
}

func New(address Address, bus int) (*I2C, error) {
	f, err := os.OpenFile(fmt.Sprintf("/dev/i2c-%d", bus), os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	if err := ioctl(f.Fd(), i2c_SLAVE, uintptr(address)); err != nil {
		return nil, err
	}
	i2c := &I2C{
		rc: f,
	}
	i2c.SetChannel(channel_default)
	i2c.SetSamplerate(samplerate_default)
	i2c.SetConversionmode(conversionmode_default)
	i2c.SetGain(gain_default)
	i2c.SetSamplerate(samplerate_default)
	err = i2c.SetAddress(address)
	if err != nil {
		return nil, err
	}
	return i2c, nil
}

func (i2c *I2C) WriteConfiguration() error {
	i2c.setRW(Write)
	_, err := i2c.rc.Write([]byte{i2c.addressbyte, i2c.configbyte})
	if err != nil {
		return err
	}
	i2c.setRW(Read)
	return nil
}

func (i2c *I2C) Read() (float32, error) {
	var buf []byte
	switch i2c.samplerate {
	case SR12, SR14, SR16:
		buf = make([]byte, 3)
	case SR18:
		buf = make([]byte, 4)
	}
	i2c.setRW(Read)
	_, err := i2c.rc.Write([]byte{i2c.addressbyte})
	if err != nil {
		return 0, err
	}
	i2c.rc.Read(buf)
	value, err := interpretvalue(i2c.samplerate, buf)
	if err != nil {
		return 0, err
	} else {
		return value, nil
	}
}

func ioctl(fd, cmd, arg uintptr) (err error) {
	_, _, e1 := syscall.Syscall6(syscall.SYS_IOCTL, fd, cmd, arg, 0, 0, 0)
	if e1 != 0 {
		err = e1
	}
	return
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

func interpretvalue(rate Samplerate, buf []byte) (float32, error) {
	var lower, middle, upper, config byte
	var value float32 = 0
	switch rate {
	case SR12, SR14, SR16:
		if len(buf) < 3 {
			return 0, errors.New("Invalid buffer length, must be 3")
		}
		upper = buf[0]
		lower = buf[1]
		config = buf[2]
	case SR18:
		if len(buf) < 4 {
			return 0, errors.New("Invalid buffer length, must be 4")
		}
		upper = buf[0]
		middle = buf[1]
		lower = buf[2]
		config = buf[3]
	}
	config = config
	switch rate {
	case SR12:
		value = float32(((upper & 0x0F) << 8) | lower)
		return value, nil
	case SR14:
		value = float32(((upper & 0x3F) << 8) | lower)
		return value, nil
	case SR16:
		value = float32((upper << 8) | lower)
		return value, nil
	case SR18:
		value = float32(((upper & 0x03) << 16) | (middle << 8) | lower)
		return value, nil
	}
	return 0, errors.New("Invalid samplerate")
}

func main() {
	i2c, err := New(A68, 1)
	checkErr(err)
	i2c.SetChannel(Ch2)
	for {
		val, err := i2c.Read()
		checkErr(err)
		fmt.Println(val)
	}
}
