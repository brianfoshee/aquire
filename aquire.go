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
                return 1000.0
        }
        reading := strings.Split(string(buf), ",")[1]
        reading = strings.Trim(reading, "\x01")
        r, err := strconv.ParseFloat(reading, 64)
        if err != nil {
                fmt.Println("getEcReading() - Error converting reading to float: ", err)
        }
        return r
}

func chipSleep(i2cAccess *i2c.I2C) {
        //sleep := []byte{0x53, 0x4C, 0x45, 0x45, 0x50}
        sleep := "Sleep"
        byteArray := []byte(sleep)
        _, err := i2cAccess.Write(byteArray)
        if err != nil {
                fmt.Println("chipSleep() - Error writing byte")
        }
        time.Sleep(time.Millisecond * 300)
}
func setProbeType(i2cAccess *i2c.I2C) {
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
func chipWake(i2cAccess *i2c.I2C) {
        _, err := i2cAccess.WriteByte(0x52)
        time.Sleep(time.Millisecond * 300)
        _, err = i2cAccess.WriteByte(0x52)
        if err != nil {
                fmt.Println("chipWake() - Error writing byte")
        }
        fmt.Println("Woke up chip")
}
func calibrateLow(i2cAccess *i2c.I2C) {
        var dummy string
        fmt.Println("Press any key to perform 4.0 calibration")
        _,_ = fmt.Scanln(&dummy)

        calibrateLow := "Cal,low,4.00"
        byteArray := []byte(calibrateLow)

        _, err := i2cAccess.Write(byteArray)
        if err != nil {
                fmt.Println("Error writing byte")
        }
        time.Sleep(time.Millisecond * 2000)
        fmt.Println("Calibration Complete")
        buf := make([]byte, 20)
        _, err = i2cAccess.Read(buf)
        if err != nil {
                fmt.Println("Error reading bytes")
        }
}

func calibrateMid(i2cAccess *i2c.I2C) {
        var dummy string
        fmt.Println("Press any key to calibrate using mid range solution")
        _,_ = fmt.Scanln(&dummy)

        calibrateMid := "Cal,mid,7.00"
        byteArray := []byte(calibrateMid)

        _, err := i2cAccess.Write(byteArray)
        if err != nil {
                fmt.Println("Error writing byte")
        }

        time.Sleep(time.Millisecond * 2000)
        fmt.Println("Calibration Complete")
        buf := make([]byte, 20)
        _, err = i2cAccess.Read(buf)
        if err != nil {
                fmt.Println("Error reading bytes")
        }
}

func calibrateHigh(i2cAccess *i2c.I2C) {
        var dummy string
        fmt.Println("Press any key to calibrate using high range solution")
        _,_ = fmt.Scanln(&dummy)

        calibrateHigh := "Cal,high,10.00"
        byteArray := []byte(calibrateHigh)
        _, err := i2cAccess.Write(byteArray)
        if err != nil {
                fmt.Println("Error writing byte")
        }
        time.Sleep(time.Millisecond * 2000)
        fmt.Println("Calibration Complete")
        buf := make([]byte, 20)
        _, err = i2cAccess.Read(buf)
        if err != nil {
                fmt.Println("Error reading bytes")
        }
}
func getStatus(i2cAccess *i2c.I2C) {
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
func main() {
        oneWire, err := onewire.NewDS18S20("28-031466321eff")
        if err != nil {
                fmt.Print(err)
        }

        ph := uint8(0x63)
        i2cPH, err := i2c.New(ph, 1)
        if err != nil {
                fmt.Println("Unable to open I2C Device")
        }

        tds := uint8(0x64)
        i2cTDS, err := i2c.New(tds, 1)
        if err != nil {
                fmt.Println("Unable to open I2C Device")
        }

        for {
                tempRaw, err := oneWire.Read()
                if err != nil {
                        fmt.Print(err)
                }
                tempC := float64(tempRaw/1000)
                tempF := float64(tempC * 9.0/5.0 + 32.0)

                stringTemp := strconv.FormatFloat(tempC, 'f', 2, 64)
                stringTemp = "T,"+stringTemp
                byteTemp := []byte(stringTemp)

                phReading := getPhReading(i2cPH, byteTemp)
                tdsReading := getEcReading(i2cTDS, byteTemp)

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
