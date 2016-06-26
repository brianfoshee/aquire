package main

import (
	"fmt"
	"github.com/alexcesaro/statsd"
	"github.com/brianfoshee/aquire/atlas"
	"github.com/brianfoshee/raspberrypi/onewire"
	"github.com/pborman/uuid"
	"os"
	"os/user"
	"strconv"
	"time"
)

func main() {

	// if temp sensor is not available, default temp to 25.428C, 77.7704F
	// using distinct number so it's obvious when temp sensor is not available/working
	var tempRaw int64 = 25428

	deviceId := uuid.NewRandom().String()
	usr, err := user.Current()
	if err != nil {
		fmt.Println(err)
	}

	deviceIdPath := usr.HomeDir + "/.hydroPiId"

	// if the device id exists on the device
	if _, err := os.Stat(deviceIdPath); err == nil {
		// open the file containing the id
		f, err := os.Open(deviceIdPath)
		if err != nil {
			fmt.Println(err)
		}

		b := make([]byte, 36)
		_, err = f.Read(b)
		if err != nil {
			fmt.Println(err)
		}

		deviceId = string(b)

		fmt.Println("Starting up using existing device id:", deviceId)

		// if the device id does not exists
	} else {
		f, err := os.Create(deviceIdPath)
		if err != nil {
			fmt.Println(err)
		}
		f.WriteString(deviceId)
		fmt.Println("Starting up with new device id:", deviceId)
		fmt.Println("Saving id to", deviceIdPath)
	}

	// open i2c communication to ph
	phChip, err := atlas.New("ph")
	if err != nil {
		fmt.Println(err)
	}

	// open i2c communication to ec
	ecChip, err := atlas.New("ec")
	if err != nil {
		fmt.Println(err)
	}

	stats, err := statsd.New(
		statsd.Address("159.203.144.95:8125"),
		statsd.Prefix("aquaponics"),
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer stats.Close()

	// Forever
	for {
		// open 1-wire communication to temp sensor
		oneWire, err := onewire.NewDS18S20("28-031466321eff")
		if err != nil {
			fmt.Println(err)
		} else {
			tempBuf, err := oneWire.Read()
			if err != nil {
				fmt.Println(err)
			} else {
				tempRaw = tempBuf
			}
		}

		// clean up reading
		tempC := float64(tempRaw) / 1000
		tempF := float64(tempC*9.0/5.0 + 32.0)

		// convert reading to bytes
		stringTemp := strconv.FormatFloat(tempC, 'f', 2, 64)
		stringTemp = "T," + stringTemp
		byteTemp := []byte(stringTemp)

		// update phChip.reading
		err = phChip.UpdateReading(byteTemp)
		if err != nil {
			fmt.Println(err)
		}

		// update ecChip.reading
		ecChip.UpdateReading(byteTemp)
		if err != nil {
			fmt.Println(err)
		}

		// access new readings
		phReading := phChip.GetReading()
		tdsReading := ecChip.GetReading()

		ns := deviceId

		t := time.Now()
		fmt.Printf("%s: temp '%v', ph'%v', tds '%v'\n", t, tempF, phReading, tdsReading)

		// send to statsd
		stats.Gauge(ns+".watertempf", tempF)
		stats.Gauge(ns+".watertempc", tempC)
		stats.Gauge(ns+".tds", tdsReading)
		stats.Gauge(ns+".ph", phReading)
	}
}
