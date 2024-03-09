// Core Electronics PiicoDev Potentiometer
package piicodev

import "fmt"

type Potentiometer struct {
	i2c              *I2C
	potType          uint16
	minimum, maximum float64
}

const (
	PotentiometerAddress = 0x35

	_DEVICE_ID_POT   = 379
	_DEVICE_ID_SLIDE = 411

	_POT_REG_WHOAMI      = 0x01
	_POT_REG_FIRM_MAJ    = 0x02
	_POT_REG_FIRM_MIN    = 0x03
	_POT_REG_I2C_ADDRESS = 0x04
	_POT_REG_POT         = 0x05
	_POT_REG_LED         = 0x07
	_POT_REG_SELF_TEST   = 0x09
)

func NewPotentiometer(addr uint8, bus int) (s *Potentiometer, err error) {
	s = &Potentiometer{
		minimum: 0.0,
		maximum: 100.0,
	}

	if s.i2c, err = OpenI2C(addr, bus); err != nil {
		return
	}

	if s.potType, err = s.i2c.ReadRegU16BE(_POT_REG_WHOAMI); err != nil {
		return
	}

	if s.potType != _DEVICE_ID_POT && s.potType != _DEVICE_ID_SLIDE {
		err = fmt.Errorf("the potentiometer type is %d rather than %d (pot) or %d (slide)", s.potType, _DEVICE_ID_POT, _DEVICE_ID_SLIDE)
		return
	}

	return
}

func (s *Potentiometer) GetType() int {
	switch s.potType {
	case _DEVICE_ID_POT:
		return 1
	case _DEVICE_ID_SLIDE:
		return 2
	}

	return 0
}

// Returns the firmware version
func (s *Potentiometer) ReadFirmwareVersion() (major uint8, minor uint8, err error) {
	if major, err = s.i2c.ReadRegU8(_POT_REG_FIRM_MAJ); err != nil {
		return
	}

	minor, err = s.i2c.ReadRegU8(_POT_REG_FIRM_MIN)

	return
}

func (s *Potentiometer) SelfTest() (v uint8, err error) {
	v, err = s.i2c.ReadRegU8(_POT_REG_SELF_TEST)
	return
}

func (s *Potentiometer) GetLED() (v uint8, err error) {
	v, err = s.i2c.ReadRegU8(_POT_REG_LED)
	return
}

func (s *Potentiometer) SetLED(enable bool) (err error) {
	var e byte = 0
	if enable {
		e = 1
	}
	err = s.i2c.WriteRegU8(_POT_REG_LED|(1<<7), e)
	return
}

// Returns a value from 0 to 1023 as the raw value from the potentiometer
func (s *Potentiometer) ReadRawValue() (v uint16, err error) {
	v, err = s.i2c.ReadRegU16BE(_POT_REG_POT)
	return
}

func (s *Potentiometer) SetMinimum(minimum float64) {
	s.minimum = minimum
}

func (s *Potentiometer) SetMaximum(maximum float64) {
	s.maximum = maximum
}

// Returns a value from s.minimum to s.maximum from the potentiometer
func (s *Potentiometer) ReadValue() (v float64, err error) {
	var rv uint16
	if rv, err = s.i2c.ReadRegU16BE(_POT_REG_POT); err != nil {
		return
	}

	v = s.minimum + (((s.maximum - s.minimum) * float64(rv)) / 1023.0)
	return
}

// Close closes the handle to the device
func (s *Potentiometer) Close() {
	s.i2c.Close()
}
