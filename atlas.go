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
	// Provides opportunity to extend to include serial
	// Improvement Note: Use enum {i2c,serial}?
	protocol string
	// Represents chip of atlas chip.
	// Improvement NOte: Use enum {ph,ec,ec,orp}?
	chip      string
	reading   float64
	i2cAccess *i2c.I2C
}

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
		fmt.Println("Unable to open DO I2C Device")
	}
	return &Atlas{i2cAccess: i2cAccess, chip: chip, reading: 0.0, protocol: protocol}, nil
}

// UpdateReading gets a reading from the appropriate atlas chip
// and stores it in Atlas.reading
func (atlas *Atlas) UpdateReading(temp []byte) {
	switch atlas.chip {
	case "do":
		atlas.reading = getDoReading(atlas.i2cAccess, temp)
	case "orp":
		atlas.reading = getOrpReading(atlas.i2cAccess, temp)
	case "ph":
		atlas.reading = getPhReading(atlas.i2cAccess, temp)
	case "ec":
		atlas.reading = getEcReading(atlas.i2cAccess, temp)
	}

}

// GetReading returns the current value of Atlas.reading
func (atlas *Atlas) GetReading() float64 {
	return atlas.reading
}

// Stub for future development
func getDoReading(i2cAccess *i2c.I2C, temp []byte) float64 {
	return 1.0
}

// Stub for future development
func getOrpReading(i2cAccess *i2c.I2C, temp []byte) float64 {
	return 1.0
}

func getPhReading(i2cAccess *i2c.I2C, temp []byte) float64 {
	_, err := i2cAccess.Write(temp)
	if err != nil {
		fmt.Println("Error writing byte")
	}
	time.Sleep(time.Millisecond * 300)

	_, err = i2cAccess.WriteByte(0x52)
	if err != nil {
		fmt.Println("Error writing byte")
	}

	time.Sleep(time.Second * 1)

	buf := make([]byte, 7)
	_, err = i2cAccess.Read(buf)
	if err != nil {
		fmt.Println("Error reading bytes")
	}
	reading := strings.Trim(string(buf), "\x01")
	reading = strings.Trim(reading, "\x00")
	r, err := strconv.ParseFloat(reading, 64)
	if err != nil {
		fmt.Println("getPhReading() - Error converting reading to float: ", err)
	}
	//atlas.reading = r
	return r
}

func getEcReading(i2cAccess *i2c.I2C, temp []byte) float64 {
	_, err := i2cAccess.Write(temp)
	if err != nil {
		fmt.Println("Error writing byte")
	}
	time.Sleep(time.Millisecond * 300)

	_, err = i2cAccess.WriteByte(0x52)
	if err != nil {
		fmt.Println("Error writing byte")
	}
	time.Sleep(time.Second * 1)

	buf := make([]byte, 32)
	_, err = i2cAccess.Read(buf)
	if err != nil {
		fmt.Println("Error reading bytes")
		//return 1000.0
	}
	reading := strings.Split(string(buf), ",")[1]
	reading = strings.Trim(reading, "\x01")
	r, err := strconv.ParseFloat(reading, 64)
	if err != nil {
		fmt.Println("getEcReading() - Error converting reading to float: ", err)
	}
	//atlas.reading = r
	return r
}

// Sleep puts the chip into a low power state
func (atlas *Atlas) Sleep() {
	sleep := "Sleep"
	byteArray := []byte(sleep)
	_, err := atlas.i2cAccess.Write(byteArray)
	if err != nil {
		fmt.Println("chipSleep() - Error writing byte")
	}
	time.Sleep(time.Millisecond * 300)
}

// SetEcProbe tells the EC chip what type of EC probe is being used
// Improvement Note: Needs to provide support for mid and high range ec probes
func (atlas *Atlas) SetEcProbe() {
	lowRange := []byte{0x4B, 0x2C, 0x30, 0x2E, 0x31}
	n, err := atlas.i2cAccess.Write(lowRange)

	if err != nil {
		fmt.Println("setProbeType() - Error writing byte")
	}

	fmt.Println("setProbeType() - Wrote bytes", n)
	time.Sleep(time.Millisecond * 1000)
	buf := make([]byte, 20)
	_, err = atlas.i2cAccess.Read(buf)
	if err != nil {
		fmt.Println("Error reading bytes")
	}
}

// Wake sends a command to the chip to wake it up from sleep mode
func (atlas *Atlas) Wake() {
	_, err := atlas.i2cAccess.WriteByte(0x52)
	time.Sleep(time.Millisecond * 300)
	_, err = atlas.i2cAccess.WriteByte(0x52)
	if err != nil {
		fmt.Println("chipWake() - Error writing byte")
	}
	time.Sleep(time.Second * 1)
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

func calibratePh(i2cAccess *i2c.I2C) {
	var dummy string
	fmt.Println("Press any key to perform 7.0 solution")
	_, _ = fmt.Scanln(&dummy)

	calibrateMid := "Cal,mid,7.00"
	byteArray := []byte(calibrateMid)

	_, err := i2cAccess.Write(byteArray)
	if err != nil {
		fmt.Println("Error writing byte")
	}

	time.Sleep(time.Millisecond * 2000)
	fmt.Println("7.0 Calibration Complete")
	buf := make([]byte, 20)
	_, err = i2cAccess.Read(buf)
	if err != nil {
		fmt.Println("Error reading bytes")
	}

	fmt.Println("Press any key to perform 4.0 calibration")
	_, _ = fmt.Scanln(&dummy)

	calibrateLow := "Cal,low,4.00"
	byteArray = []byte(calibrateLow)

	_, err = i2cAccess.Write(byteArray)
	if err != nil {
		fmt.Println("Error writing byte")
	}
	time.Sleep(time.Millisecond * 2000)
	fmt.Println("4.0 Calibration Complete")
	_, err = i2cAccess.Read(buf)
	if err != nil {
		fmt.Println("Error reading bytes")
	}

	fmt.Println("Press any key to perform 10.0 calibration")
	_, _ = fmt.Scanln(&dummy)

	calibrateHigh := "Cal,high,10.00"
	byteArray = []byte(calibrateHigh)
	_, err = i2cAccess.Write(byteArray)
	if err != nil {
		fmt.Println("Error writing byte")
	}
	time.Sleep(time.Millisecond * 2000)
	fmt.Println("10.0 Calibration Complete")
	_, err = i2cAccess.Read(buf)
	if err != nil {
		fmt.Println("Error reading bytes")
	}
}

// Stub for future development
func calibrateEc(i2cAccess *i2c.I2C) {
}

// Status returns the status of the chip
func (atlas *Atlas) Status() string {
	status := []byte{0x53, 0x54, 0x41, 0x54, 0x55, 0x53}
	_, err := atlas.i2cAccess.Write(status)

	if err != nil {
		fmt.Println("Error writing byte")
	}
	time.Sleep(time.Second * 1)

	buf := make([]byte, 20)
	_, err = atlas.i2cAccess.Read(buf)
	if err != nil {
		fmt.Println("Error reading bytes")
	}
	return (string(buf))
}
