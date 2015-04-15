package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/davecheney/i2c"
	"net/http"
	"strconv"
	"strings"
	"time"
	"github.com/raspberrypi/onewire"
)

// Info is 0x49
// Reading is 0x52
// status := []byte{0x53, 0x54, 0x41, 0x54, 0x55, 0x53}
// setProbeType(tds)
func getPhReading(addr uint8, temp []byte) float64 {
	i2cAccess, err := i2c.New(addr, 1)
	if err != nil {
		fmt.Println("Unable to open I2C Device")
	}
	_, err = i2cAccess.Write(temp)
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
	return r
}

func getEcReading(addr uint8, temp []byte) float64 {
	i2cAccess, err := i2c.New(addr, 1)
	if err != nil {
		fmt.Println("Unable to open I2C Device")
	}
	_, err = i2cAccess.Write(temp)
	if err != nil {
		fmt.Println("Error writing byte")
	}
	time.Sleep(time.Millisecond * 300)

	_, err = i2cAccess.WriteByte(0x52)
	if err != nil {
		fmt.Println("Error writing byte")
	}

	time.Sleep(time.Second * 1)

	buf := make([]byte, 20)
	_, err = i2cAccess.Read(buf)
	if err != nil {
		fmt.Println("Error reading bytes")
	}
	fmt.Println(buf)
	reading := strings.Split(string(buf), ",")[1]
	reading = strings.Trim(reading, "\x01")
	r, err := strconv.ParseFloat(reading, 64)
	if err != nil {
		fmt.Println("getEcReading() - Error converting reading to float: ", err)
	}
	return r
}

func chipSleep(addr uint8) {
	i2cAccess, err := i2c.New(addr, 1)
	if err != nil {
		fmt.Println("Unable to open I2C Device")
	}
	defer i2cAccess.Close()
	sleep := []byte{0x53, 0x4C, 0x45, 0x45, 0x50}
	n, err := i2cAccess.Write(sleep)

	if err != nil {
		fmt.Println("chipSleep() - Error writing byte")
	}

	fmt.Println("chipSleep() - Wrote bytes", n)
	time.Sleep(time.Millisecond * 300)
}
func setProbeType(addr uint8) {
	i2cAccess, err := i2c.New(addr, 1)
	if err != nil {
		fmt.Println("Unable to open I2C Device")
	}
	defer i2cAccess.Close()
	lowRange := []byte{0x4B, 0x2C, 0x30, 0x2E, 0x31}
	n, err := i2cAccess.Write(lowRange)

	if err != nil {
		fmt.Println("setProbeType() - Error writing byte")
	}

	fmt.Println("setProbeType() - Wrote bytes", n)
	time.Sleep(time.Millisecond * 1000)
	buf := make([]byte, 20)
	_, err = i2cAccess.Read(buf)
	if err != nil {
		fmt.Println("Error reading bytes")
	}
}
func chipWake(addr uint8) {
	i2cAccess, err := i2c.New(addr, 1)
	if err != nil {
		fmt.Println("Unable to open I2C Device")
	}
	defer i2cAccess.Close()
	_, err = i2cAccess.WriteByte(0x52)
	time.Sleep(time.Millisecond * 300)
	_, err = i2cAccess.WriteByte(0x52)
	if err != nil {
		fmt.Println("chipWake() - Error writing byte")
	}
	fmt.Println("Woke up chip")
}

func calibrateLow(addr uint8, temp) {
        i2cAccess, err := i2c.New(addr, 1)
        if err != nil {
                fmt.Println("chipStatus() - Unable to open I2C Device")
        }
	defer i2cAccess.Close()
	calibrateLow := "Cal,low,"+temp
	byteArray := []byte(calibrateLow)
	
	_, err = i2cAccess.Write(byteArray)
	if err != nil {
		fmt.Println("Error writing byte")
	}
	time.Sleep(time.Millisecond * 300)
}

func calibrateMid(addr uint8, temp) {
        i2cAccess, err := i2c.New(addr, 1)
        if err != nil {
                fmt.Println("chipStatus() - Unable to open I2C Device")
        }
	defer i2cAccess.Close()
	calibrateMid := "Cal,mid,"+temp
	byteArray := []byte(calibrateMid)

	_, err = i2cAccess.Write(byteArray)
	if err != nil {
		fmt.Println("Error writing byte")
	}
	time.Sleep(time.Millisecond * 300)
}

func calibrateHigh(addr uint8, temp) {
        i2cAccess, err := i2c.New(addr, 1)
        if err != nil {
                fmt.Println("chipStatus() - Unable to open I2C Device")
        }
	defer i2cAccess.Close()
	calibrateHigh := "Cal,high,"+temp
	byteArray := []byte(calibrateHigh)
	_, err = i2cAccess.Write(byteArray)
	if err != nil {
		fmt.Println("Error writing byte")
	}
	time.Sleep(time.Millisecond * 300)
}

func getStatus(addr uint8) {
	i2cAccess, err := i2c.New(addr, 1)
	if err != nil {
		fmt.Println("chipStatus() - Unable to open I2C Device")
	}
	defer i2cAccess.Close()
	status := []byte{0x53, 0x54, 0x41, 0x54, 0x55, 0x53}
	n, err := i2cAccess.Write(status)

	if err != nil {
		fmt.Println("Error writing byte")
	}

	fmt.Println("Wrote bytes", n)
	time.Sleep(time.Second * 1)

	buf := make([]byte, 20)
	r, err := i2cAccess.Read(buf)
	if err != nil {
		fmt.Println("Error reading bytes")
	}
	fmt.Printf("Read %v bytes: %v\n", r, buf)

	fmt.Println(string(buf))
}

/*

	data2 = json.dumps({	'created_at': isoDateTime,
				'sensor_data': {
					'water_temperature': fahrenheitWaterTempFloat,
				    	'ph': phReadingFloat,
				    	'tds': tdsReadingFloat
				}

*/

func main() {
	ph := uint8(0x63)
	tds := uint8(0x64)
	
	oneWire, err := onewire.NewDS18S20("28-031466321eff")
	if err != nil {
		fmt.Print(err)
	}

	for {
		tempRaw, err := oneWire.Read()
		if err != nil {
			fmt.Print(err)
		}
		tempC := float64(tempRaw/1000.0)
		tempF := float64(tempC * 9.0/5.0 + 32.0)
		
		stringTemp := strconv.FormatFloat(tempC, 'f', 2, 64)
		stringTemp = "T,"+stringTemp
		fmt.Println("String Temp: ", stringTemp)
	
		byteArray := []byte(stringTemp)
		fmt.Println(byteArray)	
	
		phReading := getPhReading(ph, byteArray)
		tdsReading := getEcReading(tds, byteArray)
	
		fmt.Println("PH: ", phReading)
		fmt.Println("TDS: ", tdsReading)
		fmt.Println("Temp: ", tempF)
	
		isoDateTime := time.Now().UTC().Format(time.RFC3339)
	
		rdgns := map[string]interface{}{
			"created_at": isoDateTime,
			"sensor_data": map[string]float64{
				"water_temperature": tempF,
				"tds":               tdsReading,
				"ph":                phReading,
			},
		}
	
		b, err := json.Marshal(rdgns)
		if err != nil {
			fmt.Println("Unable to marshal readings")
		}
		fmt.Printf("JSON: %q\n", b)
	
		buf := bytes.NewBuffer(b)
		resp, err := http.Post("http://gowebz.herokuapp.com/devices/MockClient1/readings", "application/json", buf)
		if err != nil {
			fmt.Println("Error posting data: ", err)
		}
		fmt.Println("Response: ", resp)
	}
}
