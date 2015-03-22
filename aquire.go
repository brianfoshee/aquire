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
)

// Info is 0x49
// Reading is 0x52

func getPhReading(addr uint8) float64 {
	i2cAccess, err := i2c.New(addr, 1)
	if err != nil {
		fmt.Println("Unable to open I2C Device")
	}
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

func getEcReading(addr uint8) float64 {
	i2cAccess, err := i2c.New(addr, 1)
	if err != nil {
		fmt.Println("Unable to open I2C Device")
	}
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
	//status := []byte{0x53, 0x54, 0x41, 0x54, 0x55, 0x53}
	sleep := []byte{0x53, 0x4C, 0x45, 0x45, 0x50}
	n, err := i2cAccess.Write(sleep)

	if err != nil {
		fmt.Println("chipSleep() - Error writing byte")
	}

	fmt.Println("chipSleep() - Wrote bytes", n)
	time.Sleep(time.Millisecond * 300)
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

	phReading := getPhReading(ph)
	tdsReading := getEcReading(tds)
	wTemp := 72.5
	fmt.Println("PH: ", phReading)
	fmt.Println("TDS: ", tdsReading)

	isoDateTime := time.Now().UTC().Format(time.RFC3339)

	rdgns := map[string]interface{}{
		"created_at": isoDateTime,
		"sensor_data": map[string]float64{
			"water_temperature": wTemp,
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
