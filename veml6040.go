// Core Electronics PiicoDev Colour Sensor VEML6040
// Spec sheet: https://www.vishay.com/docs/84331/designingveml6040.pdf
package piicodev

import (
	"math"
	"time"
)

const (
	VEML6040Address = 0x10

	VEML6040ConfigReg = 0x00
	VEML6040RedReg    = 0x08
	VEML6040GreenReg  = 0x09
	VEML6040BlueReg   = 0x0A
	VEML6040WhiteReg  = 0x0B

	VEML6040DefaultSettings = 0x00 // initialise gain:1x, integration 40ms, Green Sensitivity 0.25168, Max. Detectable Lux 16496, No Trig, Auto mode, enabled.
	VEML6040Shutdown        = 0x01
)

type VEML6040 struct {
	i2c *I2C
}

// NewVEML6040 creates a new VEML6040 instances
func NewVEML6040(addr uint8, bus int) (c *VEML6040, err error) {
	c = &VEML6040{}
	if c.i2c, err = OpenI2C(addr, bus); err != nil {
		return
	}

	if err = c.i2c.WriteRegU8(VEML6040ConfigReg, VEML6040Shutdown); err != nil {
		return
	}

	if err = c.i2c.WriteRegU8(VEML6040ConfigReg, VEML6040DefaultSettings); err != nil {
		return
	}

	// Need to wait for initialization
	time.Sleep(50 * time.Millisecond)
	return
}

// ReadRGB gets the red, green and blue values from the sensor
func (c *VEML6040) ReadRGBW() (red, green, blue, white uint16, err error) {
	if red, err = c.i2c.ReadRegU16LE(VEML6040RedReg); err != nil {
		return
	}

	if green, err = c.i2c.ReadRegU16LE(VEML6040GreenReg); err != nil {
		return
	}

	if blue, err = c.i2c.ReadRegU16LE(VEML6040BlueReg); err != nil {
		return
	}

	if white, err = c.i2c.ReadRegU16LE(VEML6040WhiteReg); err != nil {
		return
	}

	return
}

// Close cleans up the connection for the VEML6040 instances
func (c *VEML6040) Close() {
	c.i2c.Close()
}

// CalculateCCT calculates the correlated colour temperature (CCT) from RGB values
func CalculateCCT(red, green, blue uint16) (cct float64) {
	// Generate the XYZ colours based on the matrix in https://www.vishay.com/docs/84331/designingveml6040.pdf
	colourX := (-0.023249 * float64(red)) + (0.291014 * float64(green)) + (-0.364880 * float64(blue))
	colourY := (-0.042799 * float64(red)) + (0.272148 * float64(green)) + (-0.279591 * float64(blue))
	colourZ := (-0.155901 * float64(red)) + (0.251534 * float64(green)) + (-0.076240 * float64(blue))

	/*
		// Generate the XYZ colours based on the standard matrix
		colourX := (0.048403 * float64(red)) + (0.183633 * float64(green)) + (-0.253589 * float64(blue))
		colourY := (0.022916 * float64(red)) + (0.176388 * float64(green)) + (-0.183205 * float64(blue))
		colourZ := (-0.077436 * float64(red)) + (0.124541 * float64(green)) + (0.032081 * float64(blue))
	*/

	colourTotal := colourX + colourY + colourZ
	if colourTotal == 0 {
		return
	}

	x := colourX / colourTotal
	y := colourY / colourTotal

	// Use McCamy formula
	n := (x - 0.3320) / (0.1858 - y)
	cct = 449.0*(n*n*n) + 3525.0*(n*n) + 6823.3*n + 5520.33
	return
}

// CalculateHSV calculates the hue, saturation and value (HSV) from red, green and blue (RGB) values.
// The red, green and blue values are 0-65535.
// The hue is 0-360 and the saturation and value are 0-100.
func CalculateHSV(red, green, blue uint16) (hue, saturation, value float64) {
	r := float64(red) / 65535
	g := float64(green) / 65535
	b := float64(blue) / 65535

	min := math.Min(r, math.Min(g, b))
	max := math.Max(r, math.Max(g, b))
	del := max - min

	value = max * 100

	if del == 0 {
		// Grey with no chroma
		hue = 0
		saturation = 0
	} else {
		// Chromatic data
		saturation = (del / max) * 100

		if r == max {
			hue = (g - b) / del // between yellow and magenta
		} else if g == max {
			hue = 2.0 + (b-r)/del // between cyan and yellow
		} else if b == max {
			hue = 4.0 + (r-g)/del // between magenta and cyan
		}

		hue *= 60.0 // degrees

		if hue < 0.0 {
			hue += 360.0
		}
	}

	return
}
