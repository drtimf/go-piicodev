// Core Electronics PiicoDev Ambient Light Sensor VEML6030
package piicodev

const (
	VEML6030Address = 0x10
	_ALS_CONF = 0x00
	_REG_ALS = 0x04
	_DEFAULT_SETTINGS = 0x00    // initialise gain:1x, integration 100ms, persistence 1, disable interrupt
)

type VEML6030 struct {
	i2c *I2C
}

func NewVEML6030(addr uint8, bus int) (l *VEML6030, err error) {
	l = &VEML6030{}
	if l.i2c, err = OpenI2C(addr, bus); err != nil {
		return
	}

	err = l.i2c.WriteRegU8(_ALS_CONF, _DEFAULT_SETTINGS)
	return
}


func (l *VEML6030) Read() (light float64, err error) {
	var rawLight uint16
	if rawLight, err = l.i2c.ReadRegU16LE(_REG_ALS); err != nil {
		return
	}

	light = float64(rawLight) * 0.0288
	return
}

func (l *VEML6030) Close() {
	l.i2c.Close()
}
