package main

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"strconv"
	"time"
	"github.com/pborman/uuid"
	"github.com/brianfoshee/aquire/atlas"
	"github.com/brianfoshee/raspberrypi/onewire"
	"github.com/quipo/statsd"
)

func main() {

	deviceId := uuid.NewRandom().String()
	usr,err := user.Current()
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

		b := make([]byte,36)
		_, err = f.Read(b)
		if err != nil {
			fmt.Println(err)
		}

		deviceId = string(b)

		fmt.Println("Using existing device id:", deviceId)

	// if the device id does not exists
	} else {
		f,err := os.Create(deviceIdPath)
		if err != nil {
			fmt.Println(err)
		}
		f.WriteString(deviceId)
		fmt.Println("Generated new device id:", deviceId)
		fmt.Println("Saving id to", deviceIdPath)
	}

	// open 1-wire communication to temp sensor
	oneWire, err := onewire.NewDS18S20("28-031466321eff")
	if err != nil {
		fmt.Print(err)
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

	prefix := "aquaponics."
	statsdclient := statsd.NewStatsdClient("159.203.144.95:8125", prefix)
	if err := statsdclient.CreateSocket(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	interval := time.Second * 10
	stats := statsd.NewStatsdBuffer(interval, statsdclient)
	defer stats.Close()

	// Forever
	for {
		// grab latest reading from temp sensor
		tempRaw, err := oneWire.Read()
		if err != nil {
			fmt.Print(err)
			tempRaw = 2222
		}

		// clean up reading
		tempC := float64(tempRaw / 1000)
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

		// send to statsd
		stats.FGauge(ns+".watertempf", tempF)
		stats.FGauge(ns+".watertempc", tempC)
		stats.FGauge(ns+".tds", tdsReading)
		stats.FGauge(ns+".ph", phReading)
	}
}
