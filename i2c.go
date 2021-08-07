// I2C device access through the Linux /dev/i2c-{bus} for Piicodev devices
package piicodev

import (
	"os"
	"reflect"
	"syscall"
	"fmt"
	"unsafe"
)

const (
	I2C_SLAVE uintptr = 0x0703
	I2C_RDWR uintptr = 0x0707
)

type I2C struct {
	dev *os.File
	address uint8
}


type i2c_msg struct {
	addr  uint16
	flags uint16
	len   uint16
	buf   uintptr
}

type i2c_rdwr_ioctl_data struct {
	msgs uintptr
	nmsg uint32
}


// OpenI2C opens an I2C device at a particular address on a bus
func OpenI2C(address uint8, bus int) (i2c *I2C, err error) {
	i2c = &I2C{ address: address }

	if i2c.dev, err = os.OpenFile(fmt.Sprintf("/dev/i2c-%d", bus), os.O_RDWR, 0600); err != nil {
		return
	}

	var errno syscall.Errno
	if _, _, errno = syscall.Syscall(syscall.SYS_IOCTL, i2c.dev.Fd(), I2C_SLAVE, uintptr(address)); errno != 0 {
		err = fmt.Errorf("failed to set the I2C address on bus %d: %s\n", bus, errno.Error())
		return
	}

	return
}

func uintptrToByteSliceData(s []byte) uintptr {
	return uintptr(unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&s)).Data))
}

// readReg uses the RDWR ioctl call to read from an I2C register
func (i2c *I2C) ReadReg(reg byte, length int) (val []byte, err error) {
	val = make([]byte, length)

	messages := [2]i2c_msg{
		{
			addr: uint16(i2c.address),
			flags: 0,
			len: 1,
			buf: uintptr(unsafe.Pointer(&reg)),
		},
		{
			addr: uint16(i2c.address),
			flags: 1,
			len: uint16(length),
			buf: uintptrToByteSliceData(val),
		},
	}

	request := i2c_rdwr_ioctl_data{
		msgs: uintptr(unsafe.Pointer(&messages)),
		nmsg: 2,
	}

	
	var errno syscall.Errno
	if _, _, errno = syscall.Syscall(syscall.SYS_IOCTL, i2c.dev.Fd(), I2C_RDWR, uintptr(unsafe.Pointer(&request))); errno != 0 {
		err = fmt.Errorf("failed to read from I2C register 0x%X at address 0x%X: %s\n", reg, i2c.address, errno.Error())
	}

	return
}

// ReadRegU8 reads an unsigned 8-bit value a register
func (i2c *I2C) ReadRegU8(reg byte) (val byte, err error) {
	var buf []byte
	if buf, err = i2c.ReadReg(reg, 1); err != nil {
		return
	}

	val = buf[0]
	return
}

// ReadRegU16BE reads an unsigned 16-bit value in big endian format from a register
func (i2c *I2C) ReadRegU16BE(reg byte) (val uint16, err error) {
	var buf []byte
	if buf, err = i2c.ReadReg(reg, 2); err != nil {
		return
	}

	val = uint16(buf[0]) << 8 + uint16(buf[1])
	return
}

// ReadRegU16LE reads an unsigned 16-bit value in little endian format from a register
func (i2c *I2C) ReadRegU16LE(reg byte) (val uint16, err error) {
	var buf []byte
	if buf, err = i2c.ReadReg(reg, 2); err != nil {
		return
	}

	val = uint16(buf[1]) << 8 + uint16(buf[0])
	return
}

// ReadRegU24BE reads an unsigned 24-bit value in big endian format from a register
func (i2c *I2C) ReadRegU24BE(reg byte) (val uint32, err error) {
	var buf []byte
	if buf, err = i2c.ReadReg(reg, 3); err != nil {
		return
	}

	val = uint32(buf[0]) << 16 + uint32(buf[1]) << 8 + uint32(buf[2])
	return
}

// ReadReg16 uses the RDWR ioctl call to read from an I2C register with a 16-bit address
func (i2c *I2C) ReadReg16(reg uint16, length int) (val []byte, err error) {
	
	val = make([]byte, length)

	messages := [2]i2c_msg{
		{
			addr: uint16(i2c.address),
			flags: 0,
			len: 2,
			buf: uintptr(unsafe.Pointer(&([2]byte{byte((reg >> 8) & 0xFF), byte(reg & 0xFF)}))),
		},
		{
			addr: uint16(i2c.address),
			flags: 1,
			len: uint16(length),
			buf: uintptrToByteSliceData(val),
		},
	}

	request := i2c_rdwr_ioctl_data{
		msgs: uintptr(unsafe.Pointer(&messages)),
		nmsg: 2,
	}

	
	var errno syscall.Errno
	if _, _, errno = syscall.Syscall(syscall.SYS_IOCTL, i2c.dev.Fd(), I2C_RDWR, uintptr(unsafe.Pointer(&request))); errno != 0 {
		err = fmt.Errorf("failed to read from I2C register 0x%X at address 0x%X: %s\n", reg, i2c.address, errno.Error())
	}

	return
}

// ReadReg16U8 reads an unsigned 8-bit value a register
func (i2c *I2C) ReadReg16U8(reg uint16) (val byte, err error) {
	var buf []byte
	if buf, err = i2c.ReadReg16(reg, 1); err != nil {
		return
	}

	val = buf[0]
	return
}

// ReadReg16U16BE reads an unsigned 16-bit value in big endian format from a register with a 16-bit address
func (i2c *I2C) ReadReg16U16BE(reg uint16) (val uint16, err error) {
	var buf []byte
	if buf, err = i2c.ReadReg16(reg, 2); err != nil {
		return
	}

	val = uint16(buf[0]) << 8 + uint16(buf[1])
	return
}

// ReadReg16U16LE reads an unsigned 16-bit value in little endian format from a register with a 16-bit address
func (i2c *I2C) ReadReg16U16LE(reg uint16) (val uint16, err error) {
	var buf []byte
	if buf, err = i2c.ReadReg16(reg, 2); err != nil {
		return
	}

	val = uint16(buf[1]) << 8 + uint16(buf[0])
	return
}

// i2c_ioctl_rdwr_write makes a call to the ioctl RDWR with a write package and the passed data
func (i2c *I2C) i2c_ioctl_rdwr_write(val []byte) (errno syscall.Errno) {
	message := i2c_msg{
		addr: uint16(i2c.address),
		flags: 0,
		len: uint16(len(val)),
		buf: uintptrToByteSliceData(val),
	}

	request := i2c_rdwr_ioctl_data{
		msgs: uintptr(unsafe.Pointer(&message)),
		nmsg: 1,
	}

	_, _, errno = syscall.Syscall(syscall.SYS_IOCTL, i2c.dev.Fd(), I2C_RDWR, uintptr(unsafe.Pointer(&request)))
	return
}

// Write uses the RDWR ioctl call to write
func (i2c *I2C) Write(val []byte) (err error) {
	if errno := i2c.i2c_ioctl_rdwr_write(val); errno != 0 {
		err = fmt.Errorf("failed to write to I2C at address 0x%X: %s\n", i2c.address, errno.Error())
	}

	return
}

// WriteU8 uses the RDWR ioctl call to write a byte
func (i2c *I2C) WriteU8(val byte) (err error) {
	err = i2c.Write([]byte{val})
	return
}

// WriteReg uses the RDWR ioctl call to write to an I2C register
func (i2c *I2C) WriteReg(reg byte, val []byte) (err error) {
	msgbuf := append([]byte{reg}, val...)

	if errno := i2c.i2c_ioctl_rdwr_write(msgbuf); errno != 0 {
		err = fmt.Errorf("failed to write to I2C register 0x%X at address 0x%X: %s\n", reg, i2c.address, errno.Error())
	}

	return
}

// WriteRegU8 writes an unsigned 8-bit value to an I2C register
func (i2c *I2C) WriteRegU8(reg byte, val byte) (err error) {
	err = i2c.WriteReg(reg, []byte{val})
	return
}

// WriteRegU16BE writes an unsigned 16-bit big endian value to an I2C register with a 16-bit address
func (i2c *I2C) WriteRegU16BE(reg byte, val uint16) (err error) {
	err = i2c.WriteReg(reg, []byte{byte((val >> 8) & 0xFF), byte(val & 0xFF)})
	return
}

// WriteRegU16LE writes an unsigned 16-bit little endian value to an I2C register with a 16-bit address
func (i2c *I2C) WriteRegU16LE(reg byte, val uint16) (err error) {
	err = i2c.WriteReg(reg, []byte{byte(val & 0xFF), byte((val >> 8) & 0xFF)})
	return
}

// WriteReg16 uses the RDWR ioctl call to write to an I2C register with a 16-bit address
func (i2c *I2C) WriteReg16(reg uint16, val []byte) (err error) {
	msgbuf := append([]byte{byte((reg >> 8) & 0xFF), byte(reg & 0xFF)}, val...)

	if errno := i2c.i2c_ioctl_rdwr_write(msgbuf); errno != 0 {
		err = fmt.Errorf("failed to write to I2C register 0x%X at address 0x%X: %s\n", reg, i2c.address, errno.Error())
	}

	return
}

// WriteReg16U8 writes an unsigned 8-bit value to an I2C register with a 16-bit address
func (i2c *I2C) WriteReg16U8(reg uint16, val byte) (err error) {
	err = i2c.WriteReg16(reg, []byte{val})
	return
}

// WriteReg16U16BE writes an unsigned 16-bit big endian value to an I2C register with a 16-bit address
func (i2c *I2C) WriteReg16U16BE(reg uint16, val uint16) (err error) {
	err = i2c.WriteReg16(reg, []byte{byte((val >> 8) & 0xFF), byte(val & 0xFF)})
	return
}

// WriteReg16U16LE writes an unsigned 16-bit little endian value to an I2C register with a 16-bit address
func (i2c *I2C) WriteReg16U16LE(reg uint16, val uint16) (err error) {
	err = i2c.WriteReg16(reg, []byte{byte(val & 0xFF), byte((val >> 8) & 0xFF)})
	return
}

// Close an I2C device
func (i2c *I2C) Close() {
	i2c.dev.Close()
}
