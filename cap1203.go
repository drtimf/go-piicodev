// Core Electronics PiicoDev Capacitive Touch Sensor CAP1203
package piicodev

import "fmt"

const CAP1203Address = 0x28

const (
	CAP1203MainControlReg = 0x00
	CAP1203MainControlBitInt = 0

	CAP1203GeneralStatusReg = 0x02
	CAP1203InputStatusReg = 0x03

	CAP1203Input1DeltaCountReg = 0x10
	CAP1203Input2DeltaCountReg = 0x11
	CAP1203Input3DeltaCountReg = 0x12

	CAP1203SensitivityControlReg = 0x1F
	CAP1203SensitivityControlBitDeltaSense = 4

	CAP1203MultipleTouchConfigReg = 0x2A

	CAP1203ProdIDReg = 0xFD
	CAP1203ProdIDValue = 0x6D
)

type CAP1203 struct {
	i2c *I2C
}

func NewCAP1203(addr uint8, bus int) (c *CAP1203, err error) {
	c = &CAP1203{}
	if c.i2c, err = OpenI2C(addr, bus); err != nil {
		return
	}

	var prodID byte
	if prodID, err = c.i2c.ReadRegU8(CAP1203ProdIDReg); err != nil {
		return
	}

	if prodID != CAP1203ProdIDValue {
		err = fmt.Errorf("CAP1203 product ID of 0x%x is not 0x%x", prodID, CAP1203ProdIDValue)
		return
	}

	err = c.SetSensitivity(3)
	return
}

func (c *CAP1203) GetSensitivity() (sensitivity int, err error) {
	sensitivity, err = c.i2c.ReadRegBits(CAP1203SensitivityControlReg, CAP1203SensitivityControlBitDeltaSense, 3)
	return
}

func (c *CAP1203) SetSensitivity(sensitivity int) (err error) {
	err = c.i2c.WriteRegBits(CAP1203SensitivityControlReg, CAP1203SensitivityControlBitDeltaSense, 3, sensitivity)
	return
}

func (c *CAP1203) GetMultipleTouchEnabled() (enabled bool, err error) {
	enabled, err = c.i2c.ReadRegBit(CAP1203MultipleTouchConfigReg, 7)
	return
}

func (c *CAP1203) SetMultipleTouchEnabled(enabled bool) (err error) {
	err = c.i2c.WriteRegBit(CAP1203MultipleTouchConfigReg, 7, enabled)
	return
}

func bitToBool(v byte, pos uint) bool {
	if ((v >> pos) & 0x01) == 1 {
		return true
	} else {
		return false
	}
}

// Clears the interrupt flag and the sensor input status flags
func (c *CAP1203) clearInterrupt() (err error) {
	c.i2c.WriteRegBit(CAP1203MainControlReg, CAP1203MainControlBitInt, false)
	return
}

func (c *CAP1203) IsTouched() (touched bool, err error) {
	touched, err = c.i2c.ReadRegBit(CAP1203GeneralStatusReg, 0)
	return
}

func (c *CAP1203) Read() (status1, status2, status3 bool, err error) {

	// Read touch sensors
	var s byte
	if s, err = c.i2c.ReadRegU8(CAP1203InputStatusReg); err != nil {
		return
	}

	c.clearInterrupt()

	status1 = bitToBool(s, 0)
	status2 = bitToBool(s, 1)
	status3 = bitToBool(s, 2)
	return
}

func (c *CAP1203) ReadDeltaCounts() (count1, count2, count3 int, err error) {
	// Raw sensor values
	var c1, c2, c3 byte
	if c1, err = c.i2c.ReadRegU8(CAP1203Input1DeltaCountReg); err != nil {
		return
	}

	if c2, err = c.i2c.ReadRegU8(CAP1203Input2DeltaCountReg); err != nil {
		return
	}

	if c3, err = c.i2c.ReadRegU8(CAP1203Input3DeltaCountReg); err != nil {
		return
	}

	count1 = int(c1)
	count2 = int(c2)
	count3 = int(c3)
	return
}

func (c *CAP1203) Close() {
	c.i2c.Close()
}

