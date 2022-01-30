// Core Electronics PiicoDev RGB LED
// The Python and firmware implementations: https://github.com/CoreElectronics/CE-PiicoDev-RGB-LED-MicroPython-Module
package piicodev

const RGBLEDAddress = 0x08

const (
	RGBLEDDeviceIDReg        = 0x00
	RGBLEDFirmwareVersionReg = 0x01
	RGBLEDControlReg         = 0x03
	RGBLEDClearReg           = 0x04
	RGBLEDBrightnessReg      = 0x06
	RGBLEDValuesReg          = 0x07
)

type RGBLED struct {
	i2c  *I2C
	leds []byte
}

// NewRGBLED creates a new RGB LED instances
func NewRGBLED(addr uint8, bus int) (l *RGBLED, err error) {
	l = &RGBLED{leds: make([]byte, 9, 9)}
	if l.i2c, err = OpenI2C(addr, bus); err != nil {
		return
	}

	l.Clear()
	err = l.Show()
	return
}

// GetDeviceID gets the device ID of the RGB LED device: REVISIT - currently failing
func (l *RGBLED) GetDeviceID() (id byte, err error) {
	id, err = l.i2c.ReadRegU8(RGBLEDDeviceIDReg)
	return
}

// GetFirmwareVersion gets the firmware version of the RGB LED device: REVISIT - currently failing
func (l *RGBLED) GetFirmwareVersion() (ver uint16, err error) {
	ver, err = l.i2c.ReadRegU16LE(RGBLEDFirmwareVersionReg)
	return
}

// EnablePowerLED sets the state of the green power LED to on or off depending on the state passed in
func (l *RGBLED) EnablePowerLED(state bool) (err error) {
	v := byte(0)
	if state {
		v = 1
	}

	err = l.i2c.WriteRegU8(RGBLEDControlReg, v)
	return
}

// SetBrightness sets the brightness of all LEDs to a level from 0 being least bright to 255 being most bright
func (l *RGBLED) SetBrightness(b byte) (err error) {
	err = l.i2c.WriteRegU8(RGBLEDBrightnessReg, b)
	return
}

// Clear turns off all the LEDs
func (l *RGBLED) Clear() (err error) {
	if err = l.i2c.WriteRegU8(RGBLEDClearReg, 1); err != nil {
		return
	}
	l.ClearPixels()
	return
}

// ClearPixels sets all the pixel values to off.
// Note that Show() needs to be called to update the LEDs.
func (l *RGBLED) ClearPixels() {
	l.FillPixels(0, 0, 0)
}

// FillPixels sets all the pixel colors to red, green and blue levels which are each from 0 being least bright to 255 being most bright for that color.
// Note that Show() needs to be called to update the LEDs.
func (l *RGBLED) FillPixels(red, green, blue byte) {
	for i := 0; i < 3; i++ {
		l.SetPixel(i, red, green, blue)
	}
}

// SetPixel sets an individual pixel color to a red, green and blue level which are each from 0 being least bright to 255 being most bright for that color.
// Note that Show() needs to be called to update the LEDs.
func (l *RGBLED) SetPixel(num int, red, green, blue byte) {
	l.leds[num*3] = red
	l.leds[(num*3)+1] = green
	l.leds[(num*3)+2] = blue
}

// Show sets the three LEDs to the colors described by the Pixel functions
func (l *RGBLED) Show() (err error) {
	err = l.i2c.WriteReg(RGBLEDValuesReg, l.leds)
	return
}

// Close cleans up the connection for the RGB LED instances
func (l *RGBLED) Close() {
	l.i2c.Close()
}
