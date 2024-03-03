// Core Electronics PiicoDev Air Quality Sensor EBS160
package piicodev

import (
	"encoding/binary"
	"fmt"
	"time"
)

type ENS160 struct {
	i2c    *I2C
	config byte
	status byte
	aqi    byte
	tvoc   uint16
	eco2   uint16
}

const (
	ENS160Address = 0x53

	_VAL_PART_ID           = 0x160
	_VAL_OPMODE_DEEP_SLEEP = 0x00
	_VAL_OPMODE_IDLE       = 0x01
	_VAL_OPMODE_STANDARD   = 0x02
	_VAL_OPMODE_RESET      = 0xF0

	_BIT_DEVICE_STATUS_NEWGPR        = 0
	_BIT_DEVICE_STATUS_NEWDAT        = 1
	_BIT_DEVICE_STATUS_VALIDITY_FLAG = 2
	_BIT_DEVICE_STATUS_STATER        = 6
	_BIT_DEVICE_STATUS_STATAS        = 7

	// Registers
	_REG_PART_ID       = 0x00
	_REG_OPMODE        = 0x10
	_REG_CONFIG        = 0x11
	_REG_COMMAND       = 0x12
	_REG_TEMP_IN       = 0x13
	_REG_RH_IN         = 0x15
	_REG_DEVICE_STATUS = 0x20
	_REG_DATA_AQI      = 0x21
	_REG_DATA_TVOC     = 0x22
	_REG_DATA_ECO2     = 0x24
	_REG_DATA_T        = 0x30
	_REG_DATA_RH       = 0x32
	_REG_DATA_MISR     = 0x38
	_REG_GPR_WRITE     = 0x40
	_REG_GPR_READ      = 0x48
)

func NewENS160(addr uint8, bus int) (s *ENS160, err error) {
	s = new(ENS160)
	if s.i2c, err = OpenI2C(addr, bus); err != nil {
		return
	}

	// REVISIT
	s.config = 0

	var part_id uint16
	if part_id, err = s.i2c.ReadRegU16LE(_REG_PART_ID); err != nil {
		return
	}

	if part_id != _VAL_PART_ID {
		err = fmt.Errorf("part ID read from ENS160 is %d rather than %d", part_id, _VAL_PART_ID)
		return
	}

	if err = s.i2c.WriteRegU8(_REG_OPMODE, _VAL_OPMODE_STANDARD); err != nil {
		return
	}

	time.Sleep(20 * time.Millisecond)

	if _, err = s.i2c.ReadRegU8(_REG_OPMODE); err != nil {
		return
	}

	time.Sleep(20 * time.Millisecond)

	if err = s.i2c.WriteRegU8(_REG_CONFIG, s.config); err != nil {
		return
	}

	return
}

func (s *ENS160) readData() (err error) {
	var status byte
	if status, err = s.i2c.ReadRegU8(_REG_DEVICE_STATUS); err != nil {
		return
	}

	if status&(1<<_BIT_DEVICE_STATUS_NEWDAT) != 0 {
		var data []byte
		if data, err = s.i2c.ReadReg(_REG_DEVICE_STATUS, 6); err != nil {
			return
		}

		s.status = data[0]
		s.aqi = data[1]
		s.tvoc = binary.LittleEndian.Uint16(data[2:4])
		s.eco2 = binary.LittleEndian.Uint16(data[4:6])
	}

	return
}

func (s *ENS160) GetTemperature() (temperature float64, err error) {
	var t uint16
	if t, err = s.i2c.ReadRegU16LE(_REG_TEMP_IN); err != nil {
		return
	}

	temperature = (float64(t) / 64.0) - 273.15
	return
}

func (s *ENS160) GetStatus() (status byte, err error) {
	if err = s.readData(); err != nil {
		return
	}

	status = s.status
	return
}

func (s *ENS160) GetOperation() (operation string, err error) {
	var status byte
	if status, err = s.GetStatus(); err != nil {
		return
	}

	switch (status >> _BIT_DEVICE_STATUS_VALIDITY_FLAG) & 0x03 {
	case 0:
		operation = "operating ok"
	case 1:
		operation = "warm-up"
	case 2:
		operation = "initial start-up"
	case 3:
		operation = "no valid output"
	}

	return
}

// Read air quality indices (AQIs)
func (s *ENS160) ReadAQI() (aqi byte, rating string, err error) {
	if err = s.readData(); err != nil {
		return
	}

	aqi = s.aqi
	switch aqi {
	case 1:
		rating = "excellent"
	case 2:
		rating = "good"
	case 3:
		rating = "moderate"
	case 4:
		rating = "poor"
	case 5:
		rating = "unhealthy"
	default:
		rating = "invalid"
	}
	return
}

// Read true volatile organic compounds (TrueVOC) including ethanol, toluene, as well as hydrogen and nitrogen dioxide
func (s *ENS160) ReadTVOC() (tvoc uint16, err error) {
	if err = s.readData(); err != nil {
		return
	}

	tvoc = s.tvoc
	return
}

// Read CO2-equivalents
func (s *ENS160) ReadECO2() (eco2 uint16, rating string, err error) {
	if err = s.readData(); err != nil {
		return
	}

	eco2 = s.eco2

	if eco2 > 1500 {
		rating = "unhealthy"
	} else if eco2 > 1000 {
		rating = "poor"
	} else if eco2 > 800 {
		rating = "fair"
	} else if eco2 > 600 {
		rating = "good"
	} else if eco2 >= 400 {
		rating = "excellent"
	} else {
		rating = "invalid"
	}

	return
}

// Close closes the handle to the device
func (s *ENS160) Close() {
	s.i2c.Close()
}

