// Core Electronics PiicoDev Distance Sensor VL53L1X
package piicodev

import (
	"fmt"
	"time"
)

const (
	VL53L1XAddress = 0x29
)

var (
	_VL51L1X_DEFAULT_CONFIGURATION = []byte{
		0x00, // 0x2d : set bit 2 and 5 to 1 for fast plus mode (1MHz I2C), else don't touch
		0x00, // 0x2e : bit 0 if I2C pulled up at 1.8V, else set bit 0 to 1 (pull up at AVDD)
		0x00, // 0x2f : bit 0 if GPIO pulled up at 1.8V, else set bit 0 to 1 (pull up at AVDD)
		0x01, // 0x30 : set bit 4 to 0 for active high interrupt and 1 for active low (bits 3:0 must be 0x1), use SetInterruptPolarity()
		0x02, // 0x31 : bit 1 = interrupt depending on the polarity, use CheckForDataReady()
		0x00, // 0x32 : not user-modifiable (NUM)
		0x02, // 0x33 : NUM
		0x08, // 0x34 : NUM
		0x00, // 0x35 : NUM
		0x08, // 0x36 : NUM
		0x10, // 0x37 : NUM
		0x01, // 0x38 : NUM
		0x01, // 0x39 : NUM
		0x00, // 0x3a : NUM
		0x00, // 0x3b : NUM
		0x00, // 0x3c : NUM
		0x00, // 0x3d : NUM
		0xff, // 0x3e : NUM
		0x00, // 0x3f : NUM
		0x0F, // 0x40 : NUM
		0x00, // 0x41 : NUM
		0x00, // 0x42 : NUM
		0x00, // 0x43 : NUM
		0x00, // 0x44 : NUM
		0x00, // 0x45 : NUM
		0x20, // 0x46 : interrupt configuration 0->level low detection, 1-> level high, 2-> Out of window, 3->In window, 0x20-> New sample ready , TBC
		0x0b, // 0x47 : NUM
		0x00, // 0x48 : NUM
		0x00, // 0x49 : NUM
		0x02, // 0x4a : NUM
		0x0a, // 0x4b : NUM
		0x21, // 0x4c : NUM
		0x00, // 0x4d : NUM
		0x00, // 0x4e : NUM
		0x05, // 0x4f : NUM
		0x00, // 0x50 : NUM
		0x00, // 0x51 : NUM
		0x00, // 0x52 : NUM
		0x00, // 0x53 : NUM
		0xc8, // 0x54 : NUM
		0x00, // 0x55 : NUM
		0x00, // 0x56 : NUM
		0x38, // 0x57 : NUM
		0xff, // 0x58 : NUM
		0x01, // 0x59 : NUM
		0x00, // 0x5a : NUM
		0x08, // 0x5b : NUM
		0x00, // 0x5c : NUM
		0x00, // 0x5d : NUM
		0x01, // 0x5e : NUM
		0xdb, // 0x5f : NUM
		0x0f, // 0x60 : NUM
		0x01, // 0x61 : NUM
		0xf1, // 0x62 : NUM
		0x0d, // 0x63 : NUM
		0x01, // 0x64 : Sigma threshold MSB (mm in 14.2 format for MSB+LSB), use SetSigmaThreshold(), default value 90 mm
		0x68, // 0x65 : Sigma threshold LSB
		0x00, // 0x66 : Min count Rate MSB (MCPS in 9.7 format for MSB+LSB), use SetSignalThreshold()
		0x80, // 0x67 : Min count Rate LSB
		0x08, // 0x68 : NUM
		0xb8, // 0x69 : NUM
		0x00, // 0x6a : NUM
		0x00, // 0x6b : NUM
		0x00, // 0x6c : Intermeasurement period MSB, 32 bits register, use SetIntermeasurementInMs()
		0x00, // 0x6d : Intermeasurement period
		0x0f, // 0x6e : Intermeasurement period
		0x89, // 0x6f : Intermeasurement period LSB
		0x00, // 0x70 : NUM
		0x00, // 0x71 : NUM
		0x00, // 0x72 : distance threshold high MSB (in mm, MSB+LSB), use SetD:tanceThreshold()
		0x00, // 0x73 : distance threshold high LSB
		0x00, // 0x74 : distance threshold low MSB ( in mm, MSB+LSB), use SetD:tanceThreshold()
		0x00, // 0x75 : distance threshold low LSB
		0x00, // 0x76 : NUM
		0x01, // 0x77 : NUM
		0x0f, // 0x78 : NUM
		0x0d, // 0x79 : NUM
		0x0e, // 0x7a : NUM
		0x0e, // 0x7b : NUM
		0x00, // 0x7c : NUM
		0x00, // 0x7d : NUM
		0x02, // 0x7e : NUM
		0xc7, // 0x7f : ROI center, use SetROI()
		0xff, // 0x80 : XY ROI (X=Width, Y=Height), use SetROI()
		0x9B, // 0x81 : NUM
		0x00, // 0x82 : NUM
		0x00, // 0x83 : NUM
		0x00, // 0x84 : NUM
		0x01, // 0x85 : NUM
		0x01, // 0x86 : clear interrupt, use ClearInterrupt()
		0x40, // 0x87 : start ranging, use StartRanging() or StopRanging(), If you want an automatic start after VL53L1X_init() call, put 0x40 in location 0x87
	}
)

type VL53L1X struct {
	i2c *I2C
}

func NewVL53L1X(addr uint8, bus int) (d *VL53L1X, err error) {
	d = &VL53L1X{}
	if d.i2c, err = OpenI2C(addr, bus); err != nil {
		return
	}

	if err = d.Reset(); err != nil {
		return
	}

	var modelID uint16
	if modelID, err = d.ReadModelID(); err != nil {
		return
	}

	if modelID != 0xEACC {
		err = fmt.Errorf("model ID of VL53L1X device is 0x%X and not 0xEACC", modelID)
		return
	}

	// Write the default configuration
	if err = d.i2c.WriteReg16(0x2D, _VL51L1X_DEFAULT_CONFIGURATION); err != nil {
		return
	}

	time.Sleep(100 * time.Millisecond)

	// the API triggers this change in VL53L1_init_and_start_range() once a
	// measurement is started; assumes MM1 and MM2 are disabled
	var v uint16
	if v, err = d.i2c.ReadReg16U16BE(0x0022); err != nil {
		return
	}

	if err = d.i2c.WriteReg16U16BE(0x001E, v*4); err != nil {
		return
	}

	time.Sleep(200 * time.Millisecond)

	return
}

func (d *VL53L1X) Reset() (err error) {
	var i byte
	for i = 0; i <= 1; i++ {
		if err = d.i2c.WriteReg16U8(0x0000, i); err != nil {
			return
		}
		time.Sleep(100 * time.Millisecond)
	}

	return
}

func (d *VL53L1X) ReadModelID() (modelID uint16, err error) {
	if modelID, err = d.i2c.ReadReg16U16BE(0x010F); err != nil {
		return
	}
	return
}

func (d *VL53L1X) Read() (rng uint16, err error) {
	var data []byte
	if data, err = d.i2c.ReadReg16(0x0089, 17); err != nil {
		return
	}

	/*
	   fmt.Println(data)
	   		range_status = data[0]
	           report_status = data[1]
	           stream_count = data[2]
	           dss_actual_effective_spads_sd0 = (data[3]<<8) + data[4]
	           peak_signal_count_rate_mcps_sd0 = (data[5]<<8) + data[6]
	           ambient_count_rate_mcps_sd0 = (data[7]<<8) + data[8]
	           sigma_sd0 = (data[9]<<8) + data[10]
	           phase_sd0 = (data[11]<<8) + data[12]
	           final_crosstalk_corrected_range_mm_sd0 = (data[13]<<8) + data[14]
	           peak_signal_count_rate_crosstalk_corrected_mcps_sd0 = (data[15]<<8) + data[16]

	           status = None
	           if range_status in (17, 2, 1, 3):
	               status = "HardwareFail"
	           elif range_status == 13:
	               status = "MinRangeFail"
	           elif range_status == 18:
	               status = "SynchronizationInt"
	           elif range_status == 5:
	               status = "OutOfBoundsFail"
	           elif range_status == 4:
	               status = "SignalFail"
	           elif range_status == 6:
	               status = "SignalFail"
	           elif range_status == 7:
	               status = "WrapTargetFail"
	           elif range_status == 12:
	               status = "XtalkSignalFail"
	           elif range_status == 8:
	               status = "RangeValidMinRangeClipped"
	           elif range_status == 9:
	               if stream_count == 0:
	                   status = "RangeValidNoWrapCheckFail"
	               else:
	                   status = "OK"
	*/
	rng = uint16(data[13])<<8 | uint16(data[14])
	return
}

func (d *VL53L1X) Close() {
	d.i2c.Close()
}
