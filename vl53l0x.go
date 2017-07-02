package vl53l0x

import (
	"encoding/binary"
	"errors"

	"github.com/kidoman/embd"
)

// Constants used for VL54L0X driver
const (
	VL53L0XAddress = 0x29
)

// Errors and things
var (
	ErrMeasureTimeout = errors.New("vl53l0x: measure timeout")
	ErrOutOfBounds    = errors.New("vl53l0x: measurement out of bounds")
)

// VL34L0XDriver represents the I2C driver for the VL34L0X proximity chip.
type VL34L0XDriver struct {
	bus     embd.I2CBus
	address byte
}

// NewDriver returns a new VL34L0X driver and starts it on the provided I2C bus.
func NewDriver(bus embd.I2CBus) *VL34L0XDriver {
	d := &VL34L0XDriver{
		bus:     bus,
		address: VL53L0XAddress,
	}

	return d
}

// Measure measures the distance detected by the driver.
func (d *VL34L0XDriver) Measure() (int, error) {
	byteA, err := d.bus.ReadByteFromReg(d.address, 0x1E)
	if err != nil {
		return 0, err
	}
	byteB, err := d.bus.ReadByteFromReg(d.address, 0x1F)
	if err != nil {
		return 0, err
	}

	d.bus.WriteByteToReg(d.address, 0x00, 0x01)

	result := int(binary.BigEndian.Uint16([]byte{byteA, byteB}))

	if result <= 20 {
		return 0, ErrOutOfBounds
	} else if result > 2000 {
		return 0, ErrOutOfBounds
	}

	return result, nil
}
