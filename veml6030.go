// Core Electronics PiicoDev Ambient Light Sensor VEML6030
// Spec sheet: https://www.vishay.com/docs/84366/veml6030.pdf
package piicodev

type VEML6030Gain uint16

const (
	VEML6030GainOneEighth  VEML6030Gain = 2
	VEML6030GainOneQuarter VEML6030Gain = 3
	VEML6030GainOne        VEML6030Gain = 0
	VEML6030GainTwo        VEML6030Gain = 1
)

type VEML6030IntegrationTime uint16

const (
	VEML6030IntegrationTime25  VEML6030IntegrationTime = 12
	VEML6030IntegrationTime50  VEML6030IntegrationTime = 8
	VEML6030IntegrationTime100 VEML6030IntegrationTime = 0
	VEML6030IntegrationTime200 VEML6030IntegrationTime = 1
	VEML6030IntegrationTime400 VEML6030IntegrationTime = 2
	VEML6030IntegrationTime800 VEML6030IntegrationTime = 3
)

const (
	VEML6030Address = 0x10

	// Registers
	_SETTING_REG            = 0x00
	_POWER_SAVE_REG         = 0x03
	_AMBIENT_LIGHT_DATA_REG = 0x04

	_ENABLE   = 1
	_DISABLE  = 0
	_SHUTDOWN = 1
	_POWER    = 0

	// Bit masks
	_SD_MASK          = 0xFFFE
	_GAIN_MASK        = 0xE7FF
	_INTEG_MASK       = 0xFC3F
	_POW_SAVE_EN_MASK = 0x06

	// Bit positions
	_NO_SHIFT  = 0
	_GAIN_POS  = 11
	_INTEG_POS = 6
)

var (
	// Values for gain (VEML6030Gain) and integration time (VEML6030IntegrationTime) from raw register values
	_gains            = []float64{1.0, 2.0, 1.0 / 8.0, 1.0 / 4.0}
	_integrationTimes = []uint16{100, 200, 400, 800, 0, 0, 0, 0, 50, 0, 0, 0, 25, 0, 0, 0}

	// Lux scale factors indexed by gain (2, 1, 1/4, 1/8) for each integration time
	// (https://www.vishay.com/docs/84367/designingveml6030.pdf)
	_eightHundredIt = []float64{.0036, .0072, .0288, .0576}
	_fourHundredIt  = []float64{.0072, .0144, .0576, .1152}
	_twoHundredIt   = []float64{.0144, .0288, .1152, .2304}
	_oneHundredIt   = []float64{.0288, .0576, .2304, .4608}
	_fiftyIt        = []float64{.0576, .1152, .4608, .9216}
	_twentyFiveIt   = []float64{.1152, .2304, .9216, 1.8432}
)

type VEML6030 struct {
	i2c *I2C
}

func NewVEML6030(addr uint8, bus int) (l *VEML6030, err error) {
	l = &VEML6030{}
	if l.i2c, err = OpenI2C(addr, bus); err != nil {
		return
	}

	err = l.PowerOn()
	return
}

func (l *VEML6030) Shutdown() (err error) {
	err = l.updateRegister(_SETTING_REG, _SD_MASK, _SHUTDOWN, _NO_SHIFT)
	return
}

func (l *VEML6030) PowerOn() (err error) {
	err = l.updateRegister(_SETTING_REG, _SD_MASK, _POWER, _NO_SHIFT)
	return
}

func (l *VEML6030) EnablePowerSave() (err error) {
	err = l.updateRegister(_POWER_SAVE_REG, _POW_SAVE_EN_MASK, _ENABLE, _NO_SHIFT)
	return
}

func (l *VEML6030) DisablePowerSave() (err error) {
	err = l.updateRegister(_POWER_SAVE_REG, _POW_SAVE_EN_MASK, _DISABLE, _NO_SHIFT)
	return
}

// GetPowerSave reads the current power save state
func (l *VEML6030) GetPowerSave() (enabled bool, err error) {
	var regVal uint16
	if regVal, err = l.i2c.ReadRegU16LE(_POWER_SAVE_REG); err != nil {
		return
	}

	if (regVal & (^uint16(_POW_SAVE_EN_MASK))) == 0 {
		enabled = false
	} else {
		enabled = true
	}

	return
}

// GetGainRaw reads the raw gain setting from: Configuration register bits 12:11
func (l *VEML6030) GetGainRaw() (gain VEML6030Gain, err error) {
	var g uint16
	if g, err = l.i2c.ReadRegU16LE(_SETTING_REG); err != nil {
		return
	}

	gain = VEML6030Gain((g & (^uint16(_GAIN_MASK))) >> _GAIN_POS)
	return
}

// GetGain reads the sensor gain
func (l *VEML6030) GetGainValue() (gain float64, err error) {
	var rawGain VEML6030Gain
	if rawGain, err = l.GetGainRaw(); err != nil {
		return
	}

	gain = _gains[rawGain]
	return
}

// SetGain sets the gain
func (l *VEML6030) SetGain(gain VEML6030Gain) (err error) {
	err = l.updateRegister(_SETTING_REG, _GAIN_MASK, uint16(gain), _GAIN_POS)
	return
}

// GetIntegrationTimeRaw reads the raw integration time setting from: Configuration register bits 9:6
func (l *VEML6030) GetIntegrationTimeRaw() (integTime VEML6030IntegrationTime, err error) {
	var it uint16
	if it, err = l.i2c.ReadRegU16LE(_SETTING_REG); err != nil {
		return
	}

	integTime = VEML6030IntegrationTime((it & (^uint16(_INTEG_MASK))) >> _INTEG_POS)
	return
}

// GetIntegrationTime reads the integration time in ms
func (l *VEML6030) GetIntegrationTimeValue() (integTime uint16, err error) {
	var rawIntegTime VEML6030IntegrationTime
	if rawIntegTime, err = l.GetIntegrationTimeRaw(); err != nil {
		return
	}

	integTime = _integrationTimes[rawIntegTime]
	return
}

// SetIntegrationTime sets the integration time
func (l *VEML6030) SetIntegrationTime(integTime VEML6030IntegrationTime) (err error) {
	err = l.updateRegister(_SETTING_REG, _INTEG_MASK, uint16(integTime), _INTEG_POS)
	return
}

// calculateLux calulates a lux from a raw light reading based on gain and integration time
func (l *VEML6030) calculateLux(rawLight uint16) (luxValue float64, err error) {
	var gain VEML6030Gain
	if gain, err = l.GetGainRaw(); err != nil {
		return
	}

	var integTime VEML6030IntegrationTime
	if integTime, err = l.GetIntegrationTimeRaw(); err != nil {
		return
	}

	var gainIndex int = 0
	switch gain {
	case VEML6030GainTwo:
		gainIndex = 0
	case VEML6030GainOne:
		gainIndex = 1
	case VEML6030GainOneQuarter:
		gainIndex = 2
	case VEML6030GainOneEighth:
		gainIndex = 3
	}

	var luxScale float64
	switch integTime {
	case VEML6030IntegrationTime800:
		luxScale = _eightHundredIt[gainIndex]
	case VEML6030IntegrationTime400:
		luxScale = _fourHundredIt[gainIndex]
	case VEML6030IntegrationTime200:
		luxScale = _twoHundredIt[gainIndex]
	case VEML6030IntegrationTime100:
		luxScale = _oneHundredIt[gainIndex]
	case VEML6030IntegrationTime50:
		luxScale = _fiftyIt[gainIndex]
	case VEML6030IntegrationTime25:
		luxScale = _twentyFiveIt[gainIndex]
	}

	luxValue = float64(rawLight) * luxScale
	return
}

// Read samples the raw light level and converts to a lux value
func (l *VEML6030) Read() (light float64, err error) {
	var rawLight uint16
	if rawLight, err = l.i2c.ReadRegU16LE(_AMBIENT_LIGHT_DATA_REG); err != nil {
		return
	}

	light, err = l.calculateLux(rawLight)
	return
}

// updateRegister updates the appropriate bits in a register
func (l *VEML6030) updateRegister(reg byte, mask uint16, bits uint16, startPos uint8) (err error) {
	var v uint16
	if v, err = l.i2c.ReadRegU16LE(reg); err != nil {
		return
	}

	v &= mask
	v |= (bits << startPos)

	err = l.i2c.WriteRegU16LE(reg, v)
	return
}

// Close closes the handle to the device
func (l *VEML6030) Close() {
	l.i2c.Close()
}
