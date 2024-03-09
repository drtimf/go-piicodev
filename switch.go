// Core Electronics PiicoDev Switch
package piicodev

import "fmt"

type Switch struct {
	i2c *I2C
}

const (
	SwitchAddress = 0x42

	_DEVICE_ID_SWITCH = 409

	_SWITCH_REG_WHOAMI                = 0x01
	_SWITCH_REG_FIRM_MAJ              = 0x02
	_SWITCH_REG_FIRM_MIN              = 0x03
	_SWITCH_REG_I2C_ADDRESS           = 0x04
	_SWITCH_REG_LED                   = 0x05
	_SWITCH_REG_IS_PRESSED            = 0x11
	_SWITCH_REG_WAS_PRESSED           = 0x12
	_SWITCH_REG_DOUBLE_PRESS_DETECTED = 0x13
	_SWITCH_REG_PRESS_COUNT           = 0x14
	_SWITCH_REG_DOUBLE_PRESS_DURATION = 0x21
	_SWITCH_REG_EMA_PARAMETER         = 0x22
	_SWITCH_REG_EMA_PERIOD            = 0x23
)

func NewSwitch(addr uint8, bus int) (s *Switch, err error) {
	s = new(Switch)

	if s.i2c, err = OpenI2C(addr, bus); err != nil {
		return
	}

	var switchID uint16
	if switchID, err = s.i2c.ReadRegU16BE(_SWITCH_REG_WHOAMI); err != nil {
		return
	}

	if switchID != _DEVICE_ID_SWITCH {
		err = fmt.Errorf("the switch ID is %d rather than %d", switchID, _SWITCH_REG_WHOAMI)
		return
	}

	return
}

// Returns the firmware version
func (s *Switch) ReadFirmwareVersion() (major uint8, minor uint8, err error) {
	if major, err = s.i2c.ReadRegU8(_SWITCH_REG_FIRM_MAJ); err != nil {
		return
	}

	minor, err = s.i2c.ReadRegU8(_SWITCH_REG_FIRM_MIN)

	return
}

func (s *Switch) GetLED() (v uint8, err error) {
	v, err = s.i2c.ReadRegU8(_SWITCH_REG_LED)
	return
}

func (s *Switch) SetLED(enable bool) (err error) {
	var e byte = 0
	if enable {
		e = 1
	}
	err = s.i2c.WriteRegU8(_SWITCH_REG_LED|(1<<7), e)
	return
}

func (s *Switch) GetPressCount() (count int, err error) {
	var c uint16
	if c, err = s.i2c.ReadRegU16BE(_SWITCH_REG_PRESS_COUNT); err != nil {
		return
	}

	count = int(c)
	return
}

func (s *Switch) IsPressed() (pressed bool, err error) {
	var p uint8
	if p, err = s.i2c.ReadRegU8(_SWITCH_REG_IS_PRESSED); err != nil {
		return
	}

	if p == 0 {
		pressed = true
	} else {
		pressed = false
	}

	return
}

func (s *Switch) WasPressed() (pressed bool, err error) {
	var p uint8
	if p, err = s.i2c.ReadRegU8(_SWITCH_REG_WAS_PRESSED); err != nil {
		return
	}

	if p == 1 {
		pressed = true
	} else {
		pressed = false
	}

	return
}

// Was pressed twice within the double-press duration
func (s *Switch) WasDoublePressed() (doublePressed bool, err error) {
	var p uint8
	if p, err = s.i2c.ReadRegU8(_SWITCH_REG_DOUBLE_PRESS_DETECTED); err != nil {
		return
	}

	if p == 1 {
		doublePressed = true
	} else {
		doublePressed = false
	}

	return
}

// If the button is pressed twice within this period (ms) a double-press will be registered (default duration is 300)
func (s *Switch) GetDoublePressDuration() (duration uint16, err error) {
	duration, err = s.i2c.ReadRegU16BE(_SWITCH_REG_DOUBLE_PRESS_DURATION)
	return
}

// If the button is pressed twice within this period (ms) a double-press will be registered
func (s *Switch) SetDoublePressDuration(duration uint16) (err error) {
	err = s.i2c.WriteRegU16BE(_SWITCH_REG_DOUBLE_PRESS_DURATION|(1<<7), duration)
	return
}

// Get the exponential moving average (EMA) parameters for the switch debounce (default parameter is 63 and period is 20)
func (s *Switch) GetDebounceEMAParameters() (parameter uint8, period uint8, err error) {
	if parameter, err = s.i2c.ReadRegU8(_SWITCH_REG_EMA_PARAMETER); err != nil {
		return
	}
	period, err = s.i2c.ReadRegU8(_SWITCH_REG_EMA_PERIOD)
	return
}

// Set the exponential moving average (EMA) parameters for the switch debounce
func (s *Switch) SetDebounceEMAParameters(parameter uint8, period uint8) (err error) {
	if err = s.i2c.WriteRegU8(_SWITCH_REG_EMA_PARAMETER|(1<<7), parameter); err != nil {
		return
	}
	err = s.i2c.WriteRegU8(_SWITCH_REG_EMA_PERIOD|(1<<7), period)
	return
}

// Close closes the handle to the device
func (s *Switch) Close() {
	s.i2c.Close()
}
