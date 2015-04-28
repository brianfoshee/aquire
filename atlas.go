package atlas

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/davecheney/i2c"
)

type Atlas struct {
	protocol string
	chip      string
	reading   float64
	i2cAccess *i2c.I2C
}

// New create an i2c connection to one of the following atlas chips:
// {do, orp, ph, ec}
func New(chip string) (*Atlas, error) {
	// Default i2c address for atlas chips
	do := uint8(0x61)
	orp := uint8(0x62)
	ph := uint8(0x63)
	ec := uint8(0x64)

	var addr uint8

	// Hard coding protocol for now
	protocol := "i2c"

	chip = strings.ToLower(chip)
	switch chip {
	case "do":
		addr = do
	case "opr":
		addr = orp
	case "ph":
		addr = ph
	case "ec":
		addr = ec
	default:
		unsupportedChip := errors.New("Unsupported chip")
		return nil, unsupportedChip
	}

	i2cAccess, err := i2c.New(addr, 1)
	if err != nil {
		return nil, err
	}

	return &Atlas{i2cAccess: i2cAccess, chip: chip, reading: 0.0, protocol: protocol}, nil
}

// UpdateReading gets a reading from the appropriate atlas chip
// and stores it in Atlas.reading
func (atlas *Atlas) UpdateReading(temp []byte) error {
	var newReading float64
	var err error

	// Call chip specific read command
	switch atlas.chip {
	case "do":
		newReading, err = getDoReading(atlas.i2cAccess, temp)
	case "orp":
		newReading, err = getOrpReading(atlas.i2cAccess, temp)
	case "ph":
		newReading, err = getPhReading(atlas.i2cAccess, temp)
	case "ec":
		newReading, err = getEcReading(atlas.i2cAccess, temp)
	default:
		newReading = 0
		err = errors.New("usupported chip selection")
	}
	if err != nil {
		return err
	}

	// Update Atlas.reading
	atlas.reading = newReading

	// Return success
	return nil
}

// GetReading returns the current value of Atlas.reading
func (atlas *Atlas) GetReading() float64 {
	return atlas.reading
}

// Stub for future development
func getDoReading(i2cAccess *i2c.I2C, temp []byte) (float64, error) {
	return 1.0, nil
}

// Stub for future development
func getOrpReading(i2cAccess *i2c.I2C, temp []byte) (float64, error) {
	return 1.0, nil
}

func getPhReading(i2cAccess *i2c.I2C, temp []byte) (float64, error) {
	// Temperature calibration
	_, err := i2cAccess.Write(temp)
	if err != nil {
		return 0, err
	}

	// Wait
	time.Sleep(time.Millisecond * 300)

	// Request data
	_, err = i2cAccess.WriteByte(0x52)
	if err != nil {
		return 0, err
	}

	// Wait
	time.Sleep(time.Second * 1)

	// Receive data
	buf := make([]byte, 7)
	_, err = i2cAccess.Read(buf)
	if err != nil {
		return 0, err
	}

	// Clean up data
	reading := strings.Trim(string(buf), "\x01")
	reading = strings.Trim(reading, "\x00")
	r, err := strconv.ParseFloat(reading, 64)
	if err != nil {
		return 0, err
	}

	// return success
	return r, nil
}

func getEcReading(i2cAccess *i2c.I2C, temp []byte) (float64, error) {
	// Temperature calibration
	_, err := i2cAccess.Write(temp)
	if err != nil {
		return 0, err
	}

	// Wait
	time.Sleep(time.Millisecond * 300)

	// Request data
	_, err = i2cAccess.WriteByte(0x52)
	if err != nil {
		return 0, err
	}

	// Wait
	time.Sleep(time.Second * 1)

	// Receive data
	buf := make([]byte, 32)
	_, err = i2cAccess.Read(buf)
	if err != nil {
		return 0, err
	}

	// Clean up data
	reading := strings.Split(string(buf), ",")[1]
	reading = strings.Trim(reading, "\x01")
	r, err := strconv.ParseFloat(reading, 64)
	if err != nil {
		return 0, err
	}

	// Return success
	return r, nil
}

// Sleep puts the chip into a low power state
func (atlas *Atlas) Sleep() error {
	sleep := "Sleep"
	byteArray := []byte(sleep)
	_, err := atlas.i2cAccess.Write(byteArray)
	if err != nil {
		return err
	}
	
	// Wait
	time.Sleep(time.Millisecond * 300)
	
	// Return success
	return nil
}

// SetEcProbe tells the EC chip what type of EC probe is being used
// Improvement Note: Needs to provide support for mid and high range ec probes
func (atlas *Atlas) SetEcProbe() error {
	
	// Send command
	lowRange := []byte{0x4B, 0x2C, 0x30, 0x2E, 0x31}
	_, err := atlas.i2cAccess.Write(lowRange)
	if err != nil {
		return err
	}

	// Wait
	time.Sleep(time.Millisecond * 1000)

	// Basic response validation
	buf := make([]byte, 20)
	_, err = atlas.i2cAccess.Read(buf)
	if err != nil {
		return err
	}

	// Return success
	return nil
}

// Wake sends a command to the chip to wake it up from sleep mode
func (atlas *Atlas) Wake() error {
	// Send arbitrary command
	_, err := atlas.i2cAccess.WriteByte(0x52)
	if err != nil {
		return err
	}

	time.Sleep(time.Millisecond * 300)
	return nil
	
}

// Calibrate uses the Atlas.chip type to calibrate the sensor
// based on what type of chip it is.
func (atlas *Atlas) Calibrate() {
	switch atlas.chip {
	case "do":
		calibrateDo(atlas.i2cAccess)
	case "orp":
		calibrateOrp(atlas.i2cAccess)
	case "ph":
		calibratePh(atlas.i2cAccess)
	case "ec":
		calibrateEc(atlas.i2cAccess)
	}
}

// Stub for future development
func calibrateDo(i2cAccess *i2c.I2C) {
}

// Stub for future development
func calibrateOrp(i2cAccess *i2c.I2C) {
}

func calibratePh(i2cAccess *i2c.I2C) error {
	var dummy string
	fmt.Println("Press any key to perform 7.0 solution")
	_, _ = fmt.Scanln(&dummy)

	calibrateMid := "Cal,mid,7.00"
	byteArray := []byte(calibrateMid)

	_, err := i2cAccess.Write(byteArray)
	if err != nil {
		return err
	}

	time.Sleep(time.Millisecond * 2000)
	fmt.Println("7.0 Calibration Complete")
	buf := make([]byte, 20)
	_, err = i2cAccess.Read(buf)
	if err != nil {
		return err
	}

	fmt.Println("Press any key to perform 4.0 calibration")
	_, _ = fmt.Scanln(&dummy)

	calibrateLow := "Cal,low,4.00"
	byteArray = []byte(calibrateLow)

	_, err = i2cAccess.Write(byteArray)
	if err != nil {
		return err
	}
	time.Sleep(time.Millisecond * 2000)
	fmt.Println("4.0 Calibration Complete")
	_, err = i2cAccess.Read(buf)
	if err != nil {
		return err
	}

	fmt.Println("Press any key to perform 10.0 calibration")
	_, _ = fmt.Scanln(&dummy)

	calibrateHigh := "Cal,high,10.00"
	byteArray = []byte(calibrateHigh)
	_, err = i2cAccess.Write(byteArray)
	if err != nil {
		return err
	}
	time.Sleep(time.Millisecond * 2000)
	fmt.Println("10.0 Calibration Complete")
	_, err = i2cAccess.Read(buf)
	if err != nil {
		return err
	}

	// Return success
	return nil
}

// Stub for future development
func calibrateEc(i2cAccess *i2c.I2C) {
}

// Status returns the status of the chip
func (atlas *Atlas) Status() (string, error) {
	status := []byte{0x53, 0x54, 0x41, 0x54, 0x55, 0x53}
	_, err := atlas.i2cAccess.Write(status)

	if err != nil {
		return "", err
	}
	time.Sleep(time.Second * 1)

	buf := make([]byte, 20)
	_, err = atlas.i2cAccess.Read(buf)
	if err != nil {
		return "", err
	}

	// Return success
	return (string(buf)), nil
}
