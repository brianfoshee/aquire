package atlas

import (
	"fmt"
	"github.com/davecheney/i2c"
	"errors"
	"strings"
)

type Atlas struct {
	// Provides opportunity to extend to include serial
	// Improvement Note: Use enum {i2c,serial}?
	protocol string
	// Represents type of atlas chip.
	// Improvement NOte: Use enum {ph,tds,ec,orp}?
	type string
	reading float64
	i2cAccess *i2c.I2C
}

func New(chip string) (*Atlas, error) {

	// Default i2c address for atlas chips
	do := uint8(0x61)
	orp := uint8(0x62)
	ph := uint8(0x63)
        tds := uint8(0x64)

	// Hard coding protocol for now
	protocol := "i2c"

	chip = strings.ToLower(chip)
	switch chip {
	case "tds":
		i2cAccess, err := i2c.New(tds, 1)
		return &AtlasI2C{i2cAccess: i2cAccess, type: "tds", reading: 0.0, protocol: protocol}
	case "ph":
		i2cAccess, err := i2c.New(ph, 1)
		return &AtlasI2C{i2cAccess: i2cAccess, type: "ph", reading: 0.0, protocol: protocol}
	case "do":
		i2cAccess, err := i2c.New(do, 1)
		return &AtlasI2C{i2cAccess: i2cAccess, type: "do", reading: 0.0, protocol: protocol}
	case "orp":
		i2cAccess, err := i2c.New(orp, 1)
		return &AtlasI2C{i2cAccess: i2cAccess, type: "orp", reading: 0.0, protocol: protocol}
	default:
		unsupportedChip := errors.New("Unsupported chip type")
		return nil,unsupportedChip
}

func (atlas *Atlas) getPhReading(temp []byte) {
	_, err := atlas.i2cAccess.Write(temp)
	if err != nil {
		fmt.Println("Error writing byte")
	}
	time.Sleep(time.Millisecond * 300)


	_, err = atlas.i2cAccess.WriteByte(0x52)
	if err != nil {
		fmt.Println("Error writing byte")
	}

	time.Sleep(time.Second * 1)

	buf := make([]byte, 7)
	_, err = atlas.i2cAccess.Read(buf)
	if err != nil {
		fmt.Println("Error reading bytes")
	}
	reading := strings.Trim(string(buf), "\x01")
	reading = strings.Trim(reading, "\x00")
	r, err := strconv.ParseFloat(reading, 64)
	if err != nil {
		fmt.Println("getPhReading() - Error converting reading to float: ", err)
	}
	atlas.reading = r
}

func (atlas *Atlas) getEcReading(temp []byte) {
	_, err := atlas.i2cAccess.Write(temp)
	if err != nil {
		fmt.Println("Error writing byte")
	}
	time.Sleep(time.Millisecond * 300)

	_, err = atlas.i2cAccess.WriteByte(0x52)
	if err != nil {
		fmt.Println("Error writing byte")
	}
	time.Sleep(time.Second * 1)

	buf := make([]byte, 32)
	_, err = atlas.i2cAccess.Read(buf)
	if err != nil {
		fmt.Println("Error reading bytes")
		return 1000.0
	}
	reading := strings.Split(string(buf), ",")[1]
	reading = strings.Trim(reading, "\x01")
	r, err := strconv.ParseFloat(reading, 64)
        if err != nil {
                fmt.Println("getEcReading() - Error converting reading to float: ", err)
        }
        atlas.reading = r
}
func (atlas *Atlas) chipSleep() {
        sleep := "Sleep"
        byteArray := []byte(sleep)
        _, err := atlas.i2cAccess.Write(byteArray)
        if err != nil {
                fmt.Println("chipSleep() - Error writing byte")
        }
        time.Sleep(time.Millisecond * 300)
}
// Improvement Note: Provide support for mid and high range ec probes
func (atlas *Atlas) setProbeType() {
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
func (atlas *Atlas) chipWake() {
        _, err := i2cAccess.WriteByte(0x52)
        time.Sleep(time.Millisecond * 300)
        _, err = i2cAccess.WriteByte(0x52)
        if err != nil {
                fmt.Println("chipWake() - Error writing byte")
        }
	time.Sleep(time.Second * 1)
}
func (atlas *Atlas) calibrateLow() {
        var dummy string
        fmt.Println("Press any key to perform 4.0 calibration")
        _,_ = fmt.Scanln(&dummy)

        calibrateLow := "Cal,low,4.00"
        byteArray := []byte(calibrateLow)

        _, err := atlas.i2cAccess.Write(byteArray)
        if err != nil {
                fmt.Println("Error writing byte")
        }
        time.Sleep(time.Millisecond * 2000)
        fmt.Println("Calibration Complete")
        buf := make([]byte, 20)
        _, err = atlas.i2cAccess.Read(buf)
        if err != nil {
                fmt.Println("Error reading bytes")
        }
}
func (atlas *Atlas) calibrateMid() {
        var dummy string
        fmt.Println("Press any key to calibrate using mid range solution")
        _,_ = fmt.Scanln(&dummy)

        calibrateMid := "Cal,mid,7.00"
        byteArray := []byte(calibrateMid)

        _, err := atlas.i2cAccess.Write(byteArray)
        if err != nil {
                fmt.Println("Error writing byte")
        }

        time.Sleep(time.Millisecond * 2000)
        fmt.Println("Calibration Complete")
        buf := make([]byte, 20)
        _, err = atlas.i2cAccess.Read(buf)
        if err != nil {
                fmt.Println("Error reading bytes")
        }
}
func (atlas *Atlas) calibrateHigh() {
        var dummy string
        fmt.Println("Press any key to calibrate using high range solution")
        _,_ = fmt.Scanln(&dummy)

        calibrateHigh := "Cal,high,10.00"
        byteArray := []byte(calibrateHigh)
        _, err := atlas.i2cAccess.Write(byteArray)
        if err != nil {
                fmt.Println("Error writing byte")
        }
        time.Sleep(time.Millisecond * 2000)
        fmt.Println("Calibration Complete")
        buf := make([]byte, 20)
        _, err = atlas.i2cAccess.Read(buf)
        if err != nil {
                fmt.Println("Error reading bytes")
        }
}
func (atlas *Atlas) getStatus() {
        status := []byte{0x53, 0x54, 0x41, 0x54, 0x55, 0x53}
        n, err := atlas.i2cAccess.Write(status)

        if err != nil {
                fmt.Println("Error writing byte")
        }
        time.Sleep(time.Second * 1)

        buf := make([]byte, 20)
        r, err := atlas.i2cAccess.Read(buf)
        if err != nil {
                fmt.Println("Error reading bytes")
        }
        fmt.Println(string(buf))
}
