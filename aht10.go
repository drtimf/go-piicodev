package piicodev

// AHT10 Temperature and Humidity Sensor
/* https://eleparts.co.kr/data/goods_attach/202306/good-pdf-12751003-1.pdf
 * https://cdn-learn.adafruit.com/assets/assets/000/091/676/original/AHT20-datasheet-2020-4-16.pdf?1591047915
 *
 * Temperature:
 * - temperature range........ -40C..+85C
 * - temperature resolution... 0.01C
 * - temperature accuracy..... +-0.3C
 * - response time............ 5..30sec
 * - measurement with high frequency leads to heating of the
 *   sensor, must be > 2 seconds apart to keep self-heating below 0.1C
 *
 * Humidity:
 * - relative humidity range........ 0%..100%
 * - relative humidity resolution... 0.024%
 * - relative humidity accuracy..... +-2%
 * - response time............ 5..30sec
 * - measurement with high frequency leads to heating of the
 *   sensor, must be > 2 seconds apart to keep self-heating below 0.1C
 * - long-term exposure for 60 hours outside the normal range
 *   (humidity > 80%) can lead to a temporary drift of the
 *   signal +3%, sensor slowly returns to the calibrated state at normal
 *   operating conditions
 */

import (
	"fmt"
	"time"
)

type AHT10 struct {
	i2c *I2C
}

const (
	AHT10Address = 0x38

	_AHT1X_REG_INIT              = 0xBE //initialization register
	_AHTXX_REG_STATUS            = 0x71 //read status byte register
	_AHTXX_REG_START_MEASUREMENT = 0xAC //start measurement register
	_AHTXX_REG_SOFT_RESET        = 0xBA //soft reset register

	// calibration register controls
	_AHT1X_INIT_CTRL_NORMAL_MODE = 0x00 //normal mode on/off       bit[6:5]
	_AHT1X_INIT_CTRL_CYCLE_MODE  = 0x20 //cycle mode on/off        bit[6:5]
	_AHT1X_INIT_CTRL_CMD_MODE    = 0x40 //command mode  on/off     bit[6:5]
	_AHTXX_INIT_CTRL_CAL_ON      = 0x08 //calibration coeff on/off bit[3]

	// status byte register controls
	_AHTXX_STATUS_CTRL_BUSY        = 0x80 //busy                      bit[7]
	_AHT1X_STATUS_CTRL_NORMAL_MODE = 0x00 //normal mode status        bit[6:5]
	_AHT1X_STATUS_CTRL_CYCLE_MODE  = 0x20 //cycle mode status         bit[6:5]
	_AHT1X_STATUS_CTRL_CMD_MODE    = 0x40 //command mode status       bit[6:5]
	_AHTXX_STATUS_CTRL_CRC         = 0x10 //CRC8 status               bit[4], no info in datasheet
	_AHTXX_STATUS_CTRL_CAL_ON      = 0x08 //Calibration coeff status  bit[3]
)

func NewAHT10(addr uint8, bus int) (s *AHT10, err error) {
	s = new(AHT10)

	if s.i2c, err = OpenI2C(addr, bus); err != nil {
		return
	}

	if err = s.SoftReset(); err != nil {
		return
	}

	time.Sleep(100 * time.Millisecond)
	if err = s.SetInitRegister(_AHTXX_INIT_CTRL_CAL_ON | _AHT1X_STATUS_CTRL_NORMAL_MODE); err != nil {
		return
	}

	time.Sleep(100 * time.Millisecond)
	return
}

func calculateCRC(data []uint8) (crc uint8) {
	crc = 0xFF
	for i := 0; i < len(data); i++ {
		crc ^= data[i]
		for bi := 8; bi > 0; bi-- {
			if (crc & 0x80) != 0 {
				crc = (crc << 1) ^ 0x31
			} else {
				crc = crc << 1
			}
		}
	}

	return
}

func (s *AHT10) ReadSensor() (temperature float64, humidity float64, err error) {
	// Start measurement
	if err = s.i2c.WriteReg(_AHTXX_REG_START_MEASUREMENT, []byte{0x33, 0x00}); err != nil {
		return
	}

	// Wait for measurement to complete (should be less than 80ms)
	var status uint8

	for i := 0; i < 20; i++ {
		time.Sleep(10 * time.Millisecond)

		if status, err = s.GetStatus(); err != nil {
			return
		}

		if (status & _AHTXX_STATUS_CTRL_BUSY) == 0 {
			break
		}
	}

	if (status & _AHTXX_STATUS_CTRL_BUSY) == _AHTXX_STATUS_CTRL_BUSY {
		err = fmt.Errorf("timeout waiting for AHT10 sensor measurement to complete")
		return
	}

	// Sensor measurement values: {status, RH, RH, RH+T, T, T, CRC}
	var data []byte
	if data, err = s.i2c.ReadReg(_AHTXX_REG_STATUS, 7); err != nil {
		return
	}

	crc := calculateCRC(data)
	if crc != 0 {
		err = fmt.Errorf("the calculated CRC of AHT10 sensor data is not zero: 0x%x", crc)
		return
	}

	// Temperature (C) = T/(2^20) * 200 - 50
	temperature = ((float64((uint32(data[3])&0x0F)<<16|uint32(data[4])<<8|uint32(data[5])) * 200) / 0x100000) - 50

	// Humidity (%) = RH/(2^20) * 100
	humidity = (float64(uint32(data[1])<<12|uint32(data[2])<<4|uint32(data[3])>>4) * 100) / 0x100000

	return
}

func (s *AHT10) SetInitRegister(value uint8) (err error) {
	err = s.i2c.WriteReg(_AHT1X_REG_INIT, []byte{value, 0})
	return
}

/*
 *
 * AHT1x status register controls:
 *      7    6    5    4   3    2   1   0
 *      BSY, MOD, MOD, xx, CAL, xx, xx, xx
 *      - BSY:
 *        - 1, sensor busy/measuring
 *        - 0, sensor idle/sleeping
 *      - MOD:
 *        - 00, normal mode
 *        - 01, cycle mode
 *        - 1x, comand mode
 *      - CAL:
 *        - 1, calibration on
 *        - 0, calibration off
 */
func (s *AHT10) GetStatus() (status uint8, err error) {
	status, err = s.i2c.ReadRegU8(_AHTXX_REG_STATUS)
	return
}

func (s *AHT10) SoftReset() (err error) {
	err = s.i2c.Write([]byte{_AHTXX_REG_SOFT_RESET})
	return
}

// Close closes the handle to the device
func (s *AHT10) Close() {
	s.i2c.Close()
}
