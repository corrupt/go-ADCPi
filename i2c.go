package adcpi

import (
	"errors"
	"fmt"
	"github.com/corrupt/go-smbus"
	"os"
)

const (
	configbyte_default     = 0x00
	conversionmode_default = Continuous
	samplerate_default     = SR12
	channel_default        = Ch1
	gain_default           = X1
)

type Samplerate uint8
type Channel uint8
type Address uint8
type Conversionmode bool
type Gain uint8

type i2cconf struct {
	smb            *smbus.SMBus
	address        Address
	channel        Channel
	conversionmode Conversionmode
	samplerate     Samplerate
	gain           Gain
	configbyte     byte
	pga            float32
	lsb            float32
}

const (
	Ch1 Channel = 1
	Ch2 Channel = 2
	Ch3 Channel = 3
	Ch4 Channel = 4
	Ch5 Channel = 5
	Ch6 Channel = 6
	Ch7 Channel = 7
	Ch8 Channel = 8
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
	Continuous Conversionmode = true
	OneShot    Conversionmode = false
)

const (
	X1 Gain = 1 << iota
	X2
	X4
	X8
)

func (i2c *i2cconf) SetSamplerate(rate Samplerate) error {
	switch rate {
	case SR12:
		setBit(&(i2c.configbyte), 2, false)
		setBit(&(i2c.configbyte), 3, false)
		i2c.samplerate = rate
		i2c.lsb = 0.0005
		return nil
	case SR14:
		setBit(&(i2c.configbyte), 2, true)
		setBit(&(i2c.configbyte), 3, false)
		i2c.samplerate = rate
		i2c.lsb = 0.000125
		return nil
	case SR16:
		setBit(&(i2c.configbyte), 2, false)
		setBit(&(i2c.configbyte), 3, true)
		i2c.samplerate = rate
		i2c.lsb = 0.00003125
		return nil
	case SR18:
		setBit(&(i2c.configbyte), 2, true)
		setBit(&(i2c.configbyte), 3, true)
		i2c.samplerate = rate
		i2c.lsb = 0.0000078125
		return nil
	default:
		return errors.New(fmt.Sprintf("Invalid conversion rate: %d. Allowed values: %d, %d, %d, %d.", rate, SR12, SR14, SR16, SR18))
	}
}

func (i2c *i2cconf) SetChannel(channel Channel) error {
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

func (i2c *i2cconf) SetConversionmode(mode Conversionmode) {
	setBit(&(i2c.configbyte), 4, bool(mode))
}

func (i2c *i2cconf) SetGain(gain Gain) error {
	switch gain {
	case X1:
		setBit(&(i2c.configbyte), 0, false)
		setBit(&(i2c.configbyte), 1, false)
		i2c.gain = gain
		i2c.pga = 0.5
		return nil
	case X2:
		setBit(&(i2c.configbyte), 0, true)
		setBit(&(i2c.configbyte), 1, false)
		i2c.gain = gain
		i2c.pga = 1.0
		return nil
	case X4:
		setBit(&(i2c.configbyte), 0, false)
		setBit(&(i2c.configbyte), 1, true)
		i2c.gain = gain
		i2c.pga = 2.0
		return nil
	case X8:
		setBit(&(i2c.configbyte), 0, true)
		setBit(&(i2c.configbyte), 1, true)
		i2c.gain = gain
		i2c.pga = 4.0
		return nil
	default:
		return errors.New(fmt.Sprintf("Invalid gain: %d. Allowed values: %d, %d, %d, %d", X1, X2, X4, X8))
	}
}

func newI2c(address Address, bus int) (*i2cconf, error) {
	smb, err := smbus.New(uint(bus), byte(address))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	i2c := &i2cconf{
		smb: smb,
	}
	i2c.SetChannel(channel_default)
	i2c.SetSamplerate(samplerate_default)
	i2c.SetConversionmode(conversionmode_default)
	i2c.SetGain(gain_default)
	i2c.SetSamplerate(samplerate_default)
	return i2c, nil
}
