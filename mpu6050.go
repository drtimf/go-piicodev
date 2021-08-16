// Core Electronics PiicoDev Motion Sensor MPU-6050
package piicodev

import (
	"time"
)

type MPU6050AccelRange uint16

const (
	MPU6050AccelRange2G  MPU6050AccelRange = 0x00
	MPU6050AccelRange4G  MPU6050AccelRange = 0x08
	MPU6050AccelRange8G  MPU6050AccelRange = 0x10
	MPU6050AccelRange16G MPU6050AccelRange = 0x18
)

type MPU6050GyroRange uint16

const (
	MPU6050GyroRange250Deg  MPU6050GyroRange = 0x00
	MPU6050GyroRange500Deg  MPU6050GyroRange = 0x08
	MPU6050GyroRange1000Deg MPU6050GyroRange = 0x10
	MPU6050GyroRange2000Deg MPU6050GyroRange = 0x18
)

const (
	MPU6050Address = 0x68

	GRAVITIY_MS2 = 9.80665

	// Scale Modifiers
	ACC_SCLR_2G  = 16384.0
	ACC_SCLR_4G  = 8192.0
	ACC_SCLR_8G  = 4096.0
	ACC_SCLR_16G = 2048.0

	GYR_SCLR_250DEG  = 131.0
	GYR_SCLR_500DEG  = 65.5
	GYR_SCLR_1000DEG = 32.8
	GYR_SCLR_2000DEG = 16.4

	// MPU-6050 Registers
	PWR_MGMT_1 = 0x6B
	PWR_MGMT_2 = 0x6C

	SELF_TEST_X = 0x0D
	SELF_TEST_Y = 0x0E
	SELF_TEST_Z = 0x0F
	SELF_TEST_A = 0x10

	ACCEL_XOUT0 = 0x3B
	ACCEL_XOUT1 = 0x3C
	ACCEL_YOUT0 = 0x3D
	ACCEL_YOUT1 = 0x3E
	ACCEL_ZOUT0 = 0x3F
	ACCEL_ZOUT1 = 0x40

	TEMP_OUT0 = 0x41
	TEMP_OUT1 = 0x42

	GYRO_XOUT0 = 0x43
	GYRO_XOUT1 = 0x44
	GYRO_YOUT0 = 0x45
	GYRO_YOUT1 = 0x46
	GYRO_ZOUT0 = 0x47
	GYRO_ZOUT1 = 0x48

	ACCEL_CONFIG = 0x1C
	GYRO_CONFIG  = 0x1B
)

type MPU6050 struct {
	i2c *I2C
}

func NewMPU6050(addr uint8, bus int) (t *MPU6050, err error) {
	t = &MPU6050{}
	if t.i2c, err = OpenI2C(addr, bus); err != nil {
		return
	}

	// Wake up the MPU-6050 since it starts in sleep mode
	for i := 0; i < 3; i++ {
		if err = t.i2c.WriteRegU8(PWR_MGMT_1, 0); err != nil {
			return
		}
		time.Sleep(5 * time.Millisecond)
	}

	return
}

func (t *MPU6050) ReadTemperature() (tempC float64, err error) {
	var rawTemp int16
	if rawTemp, err = t.i2c.ReadRegS16BE(TEMP_OUT0); err != nil {
		return
	}

	tempC = (float64(rawTemp) / 340) + 36.53
	return
}

func (t *MPU6050) SetAccelRange(r MPU6050AccelRange) (err error) {
	err = t.i2c.WriteRegU16BE(ACCEL_CONFIG, uint16(r))
	return
}

func (t *MPU6050) GetAccelRangeRaw() (r MPU6050AccelRange, err error) {
	var rawRange uint16
	if rawRange, err = t.i2c.ReadRegU16BE(ACCEL_CONFIG); err != nil {
		return
	}

	r = MPU6050AccelRange(rawRange)
	return
}

func (t *MPU6050) GetAccelRangeValue() (r int, err error) {
	var rawRange MPU6050AccelRange
	if rawRange, err = t.GetAccelRangeRaw(); err != nil {
		return
	}

	switch rawRange {
	case MPU6050AccelRange2G:
		r = 2
	case MPU6050AccelRange4G:
		r = 4
	case MPU6050AccelRange8G:
		r = 8
	case MPU6050AccelRange16G:
		r = 16
	}

	return
}

func (t *MPU6050) ReadAccelData() (x, y, z float64, err error) {
	var accelRange MPU6050AccelRange
	if accelRange, err = t.GetAccelRangeRaw(); err != nil {
		return
	}

	var scaler float64
	switch accelRange {
	case MPU6050AccelRange2G:
		scaler = ACC_SCLR_2G
	case MPU6050AccelRange4G:
		scaler = ACC_SCLR_4G
	case MPU6050AccelRange8G:
		scaler = ACC_SCLR_8G
	case MPU6050AccelRange16G:
		scaler = ACC_SCLR_16G
	}

	var aX, aY, aZ int16
	if aX, err = t.i2c.ReadRegS16BE(ACCEL_XOUT0); err != nil {
		return
	}

	if aY, err = t.i2c.ReadRegS16BE(ACCEL_YOUT0); err != nil {
		return
	}

	if aZ, err = t.i2c.ReadRegS16BE(ACCEL_ZOUT0); err != nil {
		return
	}

	x = (float64(aX) * GRAVITIY_MS2) / scaler
	y = (float64(aY) * GRAVITIY_MS2) / scaler
	z = (float64(aZ) * GRAVITIY_MS2) / scaler

	return
}

func (t *MPU6050) SetGyroRange(r MPU6050GyroRange) (err error) {
	err = t.i2c.WriteRegU16BE(GYRO_CONFIG, uint16(r))
	return
}

func (t *MPU6050) GetGyroRangeRaw() (r MPU6050GyroRange, err error) {
	var rawRange uint16
	if rawRange, err = t.i2c.ReadRegU16BE(GYRO_CONFIG); err != nil {
		return
	}

	r = MPU6050GyroRange(rawRange)
	return
}

func (t *MPU6050) GetGyroRangeValue() (r int, err error) {
	var rawRange MPU6050GyroRange
	if rawRange, err = t.GetGyroRangeRaw(); err != nil {
		return
	}

	switch rawRange {
	case MPU6050GyroRange250Deg:
		r = 250
	case MPU6050GyroRange500Deg:
		r = 500
	case MPU6050GyroRange1000Deg:
		r = 1000
	case MPU6050GyroRange2000Deg:
		r = 2000
	}

	return
}

func (t *MPU6050) ReadGyroData() (x, y, z float64, err error) {
	var gyroRange MPU6050GyroRange
	if gyroRange, err = t.GetGyroRangeRaw(); err != nil {
		return
	}

	var scaler float64
	switch gyroRange {
	case MPU6050GyroRange250Deg:
		scaler = GYR_SCLR_250DEG
	case MPU6050GyroRange500Deg:
		scaler = GYR_SCLR_500DEG
	case MPU6050GyroRange1000Deg:
		scaler = GYR_SCLR_1000DEG
	case MPU6050GyroRange2000Deg:
		scaler = GYR_SCLR_2000DEG
	}

	var gX, gY, gZ int16
	if gX, err = t.i2c.ReadRegS16BE(GYRO_XOUT0); err != nil {
		return
	}

	if gY, err = t.i2c.ReadRegS16BE(GYRO_YOUT0); err != nil {
		return
	}

	if gZ, err = t.i2c.ReadRegS16BE(GYRO_ZOUT0); err != nil {
		return
	}

	x = float64(gX) / scaler
	y = float64(gY) / scaler
	z = float64(gZ) / scaler

	return
}

func (t *MPU6050) Close() {
	t.i2c.Close()
}
