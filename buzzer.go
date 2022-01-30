// Core Electronics PiicoDev Buzzer
// The Python and firmware implementations: https://github.com/CoreElectronics/CE-PiicoDev-Buzzer-MicroPython-Module
package piicodev

const BuzzerAddress = 0x5C

const (
	BuzzerStatusReg               = 0x01
	BuzzerFirmwareVersionMajorReg = 0x02
	BuzzerFirmwareVersionMinorReg = 0x03
	BuzzerI2CAddressReg           = 0x04
	BuzzerToneReg                 = 0x05
	BuzzerVolumeReg               = 0x06
	BuzzerPowerLEDReg             = 0x07
	BuzzerDeviceIDReg             = 0x11
)

type Buzzer struct {
	i2c *I2C
}

// NewBuzzer creates a new Buzzer instances
func NewBuzzer(addr uint8, bus int) (b *Buzzer, err error) {
	b = &Buzzer{}
	if b.i2c, err = OpenI2C(addr, bus); err != nil {
		return
	}

	return
}

// GetDeviceID gets the device ID of the Buzzer (should be 0x51)
func (b *Buzzer) GetDeviceID() (id byte, err error) {
	id, err = b.i2c.ReadRegU8(BuzzerDeviceIDReg)
	return
}

// GetFirmwareVersion gets the firmware version of the Buzzer (currently 1.1)
func (b *Buzzer) GetFirmwareVersion() (ver [2]byte, err error) {
	if ver[0], err = b.i2c.ReadRegU8(BuzzerFirmwareVersionMajorReg); err != nil {
		return
	}
	ver[1], err = b.i2c.ReadRegU8(BuzzerFirmwareVersionMinorReg)
	return
}

// GetStatus reads the status of the Buzzer where bit 1 is last command succeeded and bit 2 is last command known
func (b *Buzzer) GetStatus() (status byte, err error) {
	status, err = b.i2c.ReadRegU8(BuzzerStatusReg)
	return
}

// EnablePowerLED sets the state of the green power LED to on or off depending on the state passed in
func (b *Buzzer) EnablePowerLED(state bool) (err error) {
	v := byte(0)
	if state {
		v = 1
	}

	err = b.i2c.WriteRegU8(BuzzerPowerLEDReg, v)
	return
}

// SetVolume sets the volume of the buzzer between 0 which is quietest to 2 which is loudest
func (b *Buzzer) SetVolume(volume int) (err error) {
	err = b.i2c.WriteRegU8(BuzzerVolumeReg, byte(volume))
	return
}

// SetTone sets the tone of the buzzer ...
func (b *Buzzer) SetTone(freq uint16, duration uint16) (err error) {
	err = b.i2c.WriteReg(BuzzerToneReg, []byte{byte((freq >> 8) & 0xFF), byte(freq & 0xFF),
		byte((duration >> 8) & 0xFF), byte(duration & 0xFF)})
	return
}

// NoTone turns off the buzzer
func (b *Buzzer) NoTone() (err error) {
	b.SetTone(0, 0)
	return
}

// Close cleans up the connection for the Buzzer instances
func (b *Buzzer) Close() {
	b.i2c.Close()
}
