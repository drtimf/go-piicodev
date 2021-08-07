// Core Electronics PiicoDev Pressure Seneor MS5637
// The original C implementation is https://github.com/TEConnectivity/MS5637_Generic_C_Driver
package piicodev

import (
	"time"
)

const (
	MS5637Address = 0x76

	// MS5637 device commands
	_SOFTRESET = 0x1E
	_MS5637_START_PRESSURE_ADC_CONVERSION = 0x40
	_MS5637_START_TEMPERATURE_ADC_CONVERSION = 0x50
	_MS5637_CONVERSION_OSR_MASK = 0x0F
	_ADC_READ = 0x00

	// MS5637 commands read eeprom
	_MS5637_PROM_ADDR_0 = 0xA0
	_MS5637_PROM_ADDR_1 = 0xA2
	_MS5637_PROM_ADDR_2 = 0xA4
	_MS5637_PROM_ADDR_3 = 0xA6
	_MS5637_PROM_ADDR_4 = 0xA8
	_MS5637_PROM_ADDR_5 = 0xAA
	_MS5637_PROM_ADDR_6 = 0xAC

	// MS5637 commands conversion time
	_MS5637_CONV_TIME_OSR_256 = 1    // 0.001
	_MS5637_CONV_TIME_OSR_512 = 2    // 0.002
	_MS5637_CONV_TIME_OSR_1024 = 3   // 0.003
	_MS5637_CONV_TIME_OSR_2048 = 5   // 0.005
	_MS5637_CONV_TIME_OSR_4096 = 9   // 0.009
	_MS5637_CONV_TIME_OSR_8192 = 17  // 0.017

	// MS5637 commands resolution 
	_RESOLUTION_OSR_256 = 0
	_RESOLUTION_OSR_512 = 1 
	_RESOLUTION_OSR_1024 = 2
	_RESOLUTION_OSR_2048 = 3
	_RESOLUTION_OSR_4096 = 4
	_RESOLUTION_OSR_8192 = 5

	// Coefficients indexes for temperature and pressure computation
	_MS5637_CRC_INDEX = 0
	_MS5637_PRESSURE_SENSITIVITY_INDEX = 1 
	_MS5637_PRESSURE_OFFSET_INDEX = 2
	_MS5637_TEMP_COEFF_OF_PRESSURE_SENSITIVITY_INDEX = 3
	_MS5637_TEMP_COEFF_OF_PRESSURE_OFFSET_INDEX = 4
	_MS5637_REFERENCE_TEMPERATURE_INDEX = 5
	_MS5637_TEMP_COEFF_OF_TEMPERATURE_INDEX = 6
)

type MS5637ADCParams struct {
	cmd            byte
	conversionTime time.Duration
}

type MS5637 struct {
	i2c           *I2C
	coeffs        []uint16
	tempParam     MS5637ADCParams
	pressureParam MS5637ADCParams
}

func NewMS5637(addr uint8, bus int) (p *MS5637, err error) {
	p = &MS5637{}
	if p.i2c, err = OpenI2C(addr, bus); err != nil {
		return
	}

	if err = p.i2c.WriteU8(_SOFTRESET); err != nil {
		return
	}

	time.Sleep(15 * time.Millisecond)

	p.coeffs, err = p.ReadEEPROMCoeffs()
	p.SetResolution(_RESOLUTION_OSR_8192)
	return
}

func (a *MS5637ADCParams) setResolution(res byte, cmdType byte) {
	times := []byte {
		_MS5637_CONV_TIME_OSR_256, _MS5637_CONV_TIME_OSR_512, _MS5637_CONV_TIME_OSR_1024,
		_MS5637_CONV_TIME_OSR_2048, _MS5637_CONV_TIME_OSR_4096, _MS5637_CONV_TIME_OSR_8192,
	}

	a.cmd = (byte(res) * 2) | cmdType
	a.conversionTime = time.Duration(times[(a.cmd & _MS5637_CONVERSION_OSR_MASK)/2])
}

func (p *MS5637) SetResolution(res byte) {
	p.tempParam.setResolution(res, _MS5637_START_TEMPERATURE_ADC_CONVERSION)
	p.pressureParam.setResolution(res, _MS5637_START_PRESSURE_ADC_CONVERSION)	
}

func (p *MS5637) ReadEEPROMCoeffs() (coeffs []uint16, err error) {
	coeffAddresses := []byte {
		_MS5637_PROM_ADDR_0, _MS5637_PROM_ADDR_1, _MS5637_PROM_ADDR_2, _MS5637_PROM_ADDR_3,
		_MS5637_PROM_ADDR_4, _MS5637_PROM_ADDR_5, _MS5637_PROM_ADDR_6,
	}

	coeffs = make([]uint16, 0, 7)

	for _, ca := range coeffAddresses {
		var v uint16
		if v, err = p.i2c.ReadRegU16BE(ca); err != nil {
			return
		}
		coeffs = append(coeffs, v)
	}
	return
}

func (p *MS5637) getCoeff(coeff int) int64 {
	return int64(p.coeffs[coeff])
}

func (p *MS5637) readADC(param *MS5637ADCParams) (val uint32, err error) {
	if err = p.i2c.WriteU8(param.cmd); err != nil {
		return
	}

	time.Sleep(param.conversionTime * time.Millisecond)

	val, err = p.i2c.ReadRegU24BE(_ADC_READ)
	return
}

func (p *MS5637) Read() (pressure float64, temperature float64, err error) {
	var adc_temperature, adc_pressure uint32
	var dT, temp, off, sens, pr, t2, off2, sens2 int64

	if adc_temperature, err = p.readADC(&(p.tempParam)); err != nil {
		return
	}

	if adc_pressure, err = p.readADC(&(p.pressureParam)); err != nil {
		return
	}

	// Difference between actual and reference temperature = D2 - Tref
	dT = int64(adc_temperature) - (p.getCoeff(_MS5637_REFERENCE_TEMPERATURE_INDEX) << 8)

	// Actual temperature = 2000 + dT * TEMPSENS
	temp = 2000 + ((dT * p.getCoeff(_MS5637_TEMP_COEFF_OF_TEMPERATURE_INDEX)) >> 23)

	// Second order temperature compensation
	if temp < 2000 {
		t2 = (3 * (dT  * dT)) >> 33
		off2 = 61 * (temp - 2000) * (temp - 2000) / 16 
		sens2 = 29 * (temp - 2000) * (temp - 2000) / 16 
		if temp < -1500 {
			off2 += 17 * ((temp + 1500) * (temp + 1500))
			sens2 += 9 * ((temp + 1500) * (temp + 1500))
		}
	} else {
		t2 = (5 * (dT  * dT)) >> 38
		off2 = 0
		sens2 = 0
	}

	// OFF = OFF_T1 + TCO * dT
	off = (p.getCoeff(_MS5637_PRESSURE_OFFSET_INDEX) << 17) + ((p.getCoeff(_MS5637_TEMP_COEFF_OF_PRESSURE_OFFSET_INDEX) * dT) >> 6) 
	off -= off2

	// Sensitivity at actual temperature = SENS_T1 + TCS * dT
	sens = (p.getCoeff(_MS5637_PRESSURE_SENSITIVITY_INDEX) << 16) + ((p.getCoeff(_MS5637_TEMP_COEFF_OF_PRESSURE_SENSITIVITY_INDEX) * dT) >> 7) 
	sens -= sens2

	// Temperature compensated pressure = D1 * SENS - OFF
	pr = (((int64(adc_pressure) * sens) >> 21) - off) >> 15 

	temperature = float64(temp - t2) / 100.0
	pressure = float64(pr) / 100.0

	return
}

func (p *MS5637) Close() {
	p.i2c.Close()
}
