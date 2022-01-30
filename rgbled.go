// Core Electronics PiicoDev RGB LED
package piicodev

const RGBLEDAddress = 0x08

const (
	RGBLEDDeviceIDReg = 0x00
	RGBLEDFirmwareVersionReg = 0x01
	RGBLEDControlReg = 0x03
	RGBLEDClearReg = 0x04
	RGBLEDBrightnessReg = 0x06
	RGBLEDValuesReg = 0x07
)

type RGBLED struct {
	i2c  *I2C
	leds []byte
}

func NewRGBLED(addr uint8, bus int) (l *RGBLED, err error) {
	l = &RGBLED{ leds: make([]byte, 9, 9) }
	if l.i2c, err = OpenI2C(addr, bus); err != nil {
		return
	}

	l.Clear()
	err = l.Show()
	return
}

func (l *RGBLED) GetDeviceID() (id byte, err error) {
	id, err = l.i2c.ReadRegU8(RGBLEDDeviceIDReg)
	return
}

func (l *RGBLED) GetFirmwareVersion() (ver uint16, err error) {
	ver, err = l.i2c.ReadRegU16LE(RGBLEDFirmwareVersionReg)
	return
}

func (l *RGBLED) EnablePowerLED(state bool) (err error) {
	v := byte(0)
	if state {
		v = 1
	}

	err = l.i2c.WriteRegU8(RGBLEDControlReg, v)
	return
}

func (l *RGBLED) SetBrightness(b byte) (err error) {
	err = l.i2c.WriteRegU8(RGBLEDBrightnessReg, b)
	return
}

func (l *RGBLED) Clear() (err error) {
	if err = l.i2c.WriteRegU8(RGBLEDClearReg, 1); err != nil {
		return
	}
	l.ClearPixels()
	return
}

func (l *RGBLED) ClearPixels() {
	l.FillPixels(0, 0, 0)
}

func (l *RGBLED) FillPixels(red, green, blue byte) {
	for i := 0; i < 3; i++ {
		l.SetPixel(i, red, green, blue)
	}
}

func (l *RGBLED) SetPixel(num int, red, green, blue byte) {
	l.leds[num*3] = red
	l.leds[(num*3)+1] = green
	l.leds[(num*3)+2] = blue
}

func (l *RGBLED) Show() (err error) {
	err = l.i2c.WriteReg(RGBLEDValuesReg, l.leds)
	return
}

func (l *RGBLED) Close() {
	l.i2c.Close()
}

