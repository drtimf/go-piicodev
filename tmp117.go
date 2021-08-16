// Core Electronics PiicoDev Precision Temperature Sensor TMP117
package piicodev

const TMP117Address = 0x48

type TMP117 struct {
	i2c *I2C
}

func NewTMP117(addr uint8, bus int) (t *TMP117, err error) {
	t = &TMP117{}
	t.i2c, err = OpenI2C(addr, bus)
	return
}

func (t *TMP117) ReadTempC() (tempC float64, err error) {
	var rawTemp uint16
	if rawTemp, err = t.i2c.ReadRegU16BE(0); err != nil {
		return
	}

	if rawTemp >= 0x8000 {
		tempC = -256.0 + float64(rawTemp-0x8000)*7.8125e-3
	} else {
		tempC = float64(rawTemp) * 7.8125e-3
	}

	return
}

func (t *TMP117) ReadTempF() (tempF float64, err error) {
	var tempC float64
	if tempC, err = t.ReadTempC(); err != nil {
		return
	}

	tempF = (tempC * 9.0 / 5.0) + 32.0
	return
}

func (t *TMP117) ReadTempK() (tempK float64, err error) {
	var tempC float64
	if tempC, err = t.ReadTempC(); err != nil {
		return
	}

	tempK = tempC + 273.15
	return
}

func (t *TMP117) Close() {
	t.i2c.Close()
}
