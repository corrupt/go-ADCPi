package adcpi

import (
	"errors"
	"fmt"
)

type ADCPi struct {
	i2c1    *i2cconf
	i2c2    *i2cconf
	channel int
}

func NewADC(address1, address2 Address, bus int) (*ADCPi, error) {
	i2c1, err := newI2c(address1, bus)
	if err != nil {
		return nil, err
	}
	i2c2, err := newI2c(address2, bus)
	if err != nil {
		return nil, err
	}
	adc := &ADCPi{
		i2c1:    i2c1,
		i2c2:    i2c2,
		channel: 1,
	}
	return adc, nil
}

func (adcpi *ADCPi) SetChannel(channel int) error {
	switch channel {
	case 1:
		adcpi.channel = channel
		adcpi.i2c1.SetChannel(Ch1)
		return nil
	case 2:
		adcpi.channel = channel
		adcpi.i2c1.SetChannel(Ch2)
		return nil
	case 3:
		adcpi.channel = channel
		adcpi.i2c1.SetChannel(Ch3)
		return nil
	case 4:
		adcpi.channel = channel
		adcpi.i2c1.SetChannel(Ch4)
		return nil
	case 5:
		adcpi.channel = channel
		adcpi.i2c2.SetChannel(Ch1)
		return nil
	case 6:
		adcpi.channel = channel
		adcpi.i2c2.SetChannel(Ch2)
		return nil
	case 7:
		adcpi.channel = channel
		adcpi.i2c2.SetChannel(Ch3)
		return nil
	case 8:
		adcpi.channel = channel
		adcpi.i2c2.SetChannel(Ch4)
		return nil
	default:
		return errors.New("Channel can only be between 1 and 8")
	}
}

func (adcpi *ADCPi) SetSamplerate(rate Samplerate) {
	adcpi.i2c1.SetSamplerate(rate)
	adcpi.i2c2.SetSamplerate(rate)
}

func (adcpi *ADCPi) SetGain(gain Gain) {
	adcpi.i2c1.SetGain(gain)
	adcpi.i2c2.SetGain(gain)
}

func (adcpi *ADCPi) ReadRaw(channel int) (int32, bool, error) {
	var i2c *i2cconf
	var buf []byte

	err := adcpi.SetChannel(channel)
	if err != nil {
		return -1, false, err
	}

	switch channel {
	case 1, 2, 3, 4:
		i2c = adcpi.i2c1
	case 5, 6, 7, 8:
		i2c = adcpi.i2c2
	}
	switch i2c.samplerate {
	case SR12, SR14, SR16:
		buf = make([]byte, 3)
	case SR18:
		buf = make([]byte, 4)
	}

	for {
		////fmt.Printf("Channel: %d, Configbyte: %X\n", channel, i2c.configbyte)
		i, err := i2c.smb.Read_i2c_block_data(i2c.configbyte, buf)
		if err != nil {
			return -1, false, err
		}
		if i != len(buf) {
			return -1, false, fmt.Errorf("Read errors, received incorrect number of bytes. Expected: %d, received: %d\n", len(buf), i)
		}
		updated := valueUpdated(i2c.samplerate, buf)
		val := interpretValue(i2c.samplerate, buf)
		if updated {
			return val, true, nil
		} else {
			continue
			//return val, false, nil
		}
	}
}

func (adcpi ADCPi) ReadVoltage(channel int) (float32, int32, error) {
	var lsb, pga float32

	raw, nw, err := adcpi.ReadRaw(channel)
	if err != nil {
		return -1.0, -1, err
	}
	switch channel {
	case 1, 2, 3, 4:
		lsb = adcpi.i2c1.lsb
		pga = adcpi.i2c1.pga
	case 5, 6, 7, 8:
		lsb = adcpi.i2c2.lsb
		pga = adcpi.i2c2.pga
	}

	voltage := float32(float32(raw) * (lsb / pga) * 2.471)
	if voltage < 0 {
		voltage = 0
	}
	return voltage, raw, nil
}
