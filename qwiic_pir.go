// SparkFun Electronics Qwiic PIR
// The firmware: https://github.com/sparkfun/Qwiic_PIR
// The Python implementation: https://github.com/sparkfun/Qwiic_PIR_Py
package piicodev

import "fmt"

const (
	QwiicPIRAddress = 0x12

	QwiicPIRDeviceIDReg             = 0x00
	QwiicPIRFirmwareVersionMajorReg = 0x01
	QwiicPIRFirmwareVersionMinorReg = 0x02
	QwiicPIREventStatusReg          = 0x03
	QwiicPIRInterruptConfigReg      = 0x04
	QwiicPIREventDebounceTimeReg    = 0x05
	QwiicPIRDetectedQueueStatusReg  = 0x07
	QwiicPIRDetectedQueueFrontReg   = 0x08
	QwiicPIRDetectedQueueBackReg    = 0x0C
	QwiicPIRRemovedQueueStatusReg   = 0x10
	QwiicPIRRemovedQueueFrontReg    = 0x11
	QwiicPIRRemovedQueueBackReg     = 0x15
	QwiicPIRI2cAddressReg           = 0x19

	QwiicPIRDeviceID = 0x72
)

type QwiicPIR struct {
	i2c *I2C
}

// NewQwiicPIR creates a new QwiicPIR instances
func NewQwiicPIR(addr uint8, bus int) (p *QwiicPIR, err error) {
	p = &QwiicPIR{}
	if p.i2c, err = OpenI2C(addr, bus); err != nil {
		return
	}

	var deviceID byte
	if deviceID, err = p.GetDeviceID(); err != nil {
		return
	}

	if deviceID != QwiicPIRDeviceID {
		err = fmt.Errorf("the device ID 0x%x read from the Qwiic PIR does not match 0x%x", deviceID, QwiicPIRDeviceID)
		return
	}

	return
}

// GetDeviceID gets the device ID of the QwiicPIR (should be 0x72)
func (p *QwiicPIR) GetDeviceID() (id byte, err error) {
	id, err = p.i2c.ReadRegU8(QwiicPIRDeviceIDReg)
	return
}

// GetFirmwareVersion gets the firmware version of the QwiicPIR (currently 1.1)
func (p *QwiicPIR) GetFirmwareVersion() (ver [2]byte, err error) {
	if ver[0], err = p.i2c.ReadRegU8(QwiicPIRFirmwareVersionMajorReg); err != nil {
		return
	}
	ver[1], err = p.i2c.ReadRegU8(QwiicPIRFirmwareVersionMinorReg)
	return
}

// GetRawReading detected flag returns true if the PIR has detected an object, otherwise false
func (p *QwiicPIR) GetRawReading() (detected bool, err error) {
	return p.i2c.ReadRegBit(QwiicPIREventStatusReg, 0)
}

// IsObjectDetected detected flag returns true if the PIR has detected an object after debouncing, otherwise false
func (p *QwiicPIR) IsObjectDetected() (detected bool, err error) {
	return p.i2c.ReadRegBit(QwiicPIREventStatusReg, 3)
}

// IstObjectRemoved removed flag returns true if the PIR has detected the removal of an object after debouncing, otherwise false
func (p *QwiicPIR) IstObjectRemoved() (removed bool, err error) {
	return p.i2c.ReadRegBit(QwiicPIREventStatusReg, 2)
}

// IsAvailable available flag returns true if the PIR has either a detected or a removal event
func (p *QwiicPIR) IsAvailable() (available bool, err error) {
	return p.i2c.ReadRegBit(QwiicPIREventStatusReg, 1)
}

// ClearEventBits clears the raw object detected, object removed and event available bits
func (p *QwiicPIR) ClearEventBits() (err error) {
	return p.i2c.WriteRegBits(QwiicPIREventStatusReg, 1, 3, 0)
}

// GetDebounceEvents gets the availability of the event, if it is detected or removed and clears the bits after reading
func (p *QwiicPIR) GetDebounceEvents() (available bool, detected bool, removed bool, err error) {
	available = false
	detected = false
	removed = false

	var eventStatus byte
	if eventStatus, err = p.i2c.ReadRegU8(QwiicPIREventStatusReg); err != nil {
		return
	}

	if ((eventStatus & 0x02) >> 1) == 1 {
		available = true
	}

	if ((eventStatus & 0x04) >> 2) == 1 {
		removed = true
	}

	if ((eventStatus & 0x08) >> 3) == 1 {
		detected = true
	}

	// Clear event status bits
	err = p.i2c.WriteRegU8(QwiicPIREventStatusReg, eventStatus & ^byte(0x0E))
	return
}

// GetDebounceTime returns the debounce time in milliseconds
func (p *QwiicPIR) GetDebounceTime() (debounceTime uint16, err error) {
	return p.i2c.ReadReg16U16LE(QwiicPIREventDebounceTimeReg)
}

// SetDebounceTime the time in milliseconds to set the debounce time in
func (p *QwiicPIR) SetDebounceTime(debounceTime uint16) (err error) {
	return p.i2c.WriteReg16U16LE(QwiicPIREventDebounceTimeReg, debounceTime)
}

// IsDetectedQueueFull checks in the detect queue is full
func (p *QwiicPIR) IsDetectedQueueFull() (full bool, err error) {
	return p.i2c.ReadRegBit(QwiicPIRDetectedQueueStatusReg, 2)
}

// IsDetectedQueueEmpty checks in the detect queue is empty
func (p *QwiicPIR) IsDetectedQueueEmpty() (empty bool, err error) {
	return p.i2c.ReadRegBit(QwiicPIRDetectedQueueStatusReg, 1)
}

// TimeSinceLastDetect number of milliseconds since the last detect event
func (p *QwiicPIR) TimeSinceLastDetect() (lastDetect uint32, err error) {
	return p.i2c.ReadRegU32LE(QwiicPIRDetectedQueueFrontReg)
}

// TimeSinceFirstDetect number of milliseconds since the first detect event
func (p *QwiicPIR) TimeSinceFirstDetect() (firstDetect uint32, err error) {
	return p.i2c.ReadRegU32LE(QwiicPIRDetectedQueueBackReg)
}

// PopDetectedQueue returns the oldest value in the detected queue in milliseconds and then removes it
func (p *QwiicPIR) PopDetectedQueue() (firstDetect uint32, err error) {
	if firstDetect, err = p.TimeSinceFirstDetect(); err != nil {
		return
	}
	err = p.i2c.WriteRegBit(QwiicPIRDetectedQueueStatusReg, 0, true)
	return
}

// IsRemovedQueueFull checks in the detect queue is full
func (p *QwiicPIR) IsRemovedQueueFull() (full bool, err error) {
	return p.i2c.ReadRegBit(QwiicPIRRemovedQueueStatusReg, 2)
}

// IsRemovedQueueEmpty checks in the detect queue is empty
func (p *QwiicPIR) IsRemovedQueueEmpty() (empty bool, err error) {
	return p.i2c.ReadRegBit(QwiicPIRRemovedQueueStatusReg, 1)
}

// TimeSinceLasRemove number of milliseconds since the last remove event
func (p *QwiicPIR) TimeSinceLasRemove() (lastDetect uint32, err error) {
	return p.i2c.ReadRegU32LE(QwiicPIRRemovedQueueFrontReg)
}

// TimeSinceFirstRemove number of milliseconds since the first remove event
func (p *QwiicPIR) TimeSinceFirstRemove() (firstDetect uint32, err error) {
	return p.i2c.ReadRegU32LE(QwiicPIRRemovedQueueBackReg)
}

// PopRemoveQueue returns the oldest value in the remove queue in milliseconds and then removes it
func (p *QwiicPIR) PopRemoveQueue() (firstDetect uint32, err error) {
	if firstDetect, err = p.TimeSinceFirstRemove(); err != nil {
		return
	}
	err = p.i2c.WriteRegBit(QwiicPIRRemovedQueueStatusReg, 0, true)
	return
}

// Close cleans up the connection for the QwiicPIR instances
func (p *QwiicPIR) Close() {
	p.i2c.Close()
}
