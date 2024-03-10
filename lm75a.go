package piicodev

// LM75a Temperature sensor and thermal watchdog
// https://www.nxp.com/docs/en/data-sheet/LM75A.pdf

type LM75A struct {
	i2c *I2C
}

const (
	LM75AAddress = 0x4F

	_LM75A_REG_TEMPERATURE = 0x00
)

func NewLM75A(addr uint8, bus int) (s *LM75A, err error) {
	s = new(LM75A)

	if s.i2c, err = OpenI2C(addr, bus); err != nil {
		return
	}

	return
}

func (s *LM75A) ReadTemperature() (temperature float64, err error) {
	var t uint16
	if t, err = s.i2c.ReadRegU16BE(_LM75A_REG_TEMPERATURE); err != nil {
		return
	}

	temperature = float64(t) / 256.0
	return
}

// Close closes the handle to the device
func (s *LM75A) Close() {
	s.i2c.Close()
}
