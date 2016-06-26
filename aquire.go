package main

import (
	"github.com/alexcesaro/statsd"
	"github.com/brianfoshee/aquire/atlas"
	"github.com/brianfoshee/raspberrypi/onewire"
	"github.com/pborman/uuid"
	"log"
	"os"
	"os/user"
	"strconv"
)

var (
	Debug *log.Logger
	Info  *log.Logger
	Error *log.Logger
)

func initLogging(logOutput string) {
	file, err := os.OpenFile(logOutput, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file:", err)
	}

	Debug = log.New(file, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
	Info = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func main() {
	// if temp sensor is not available, default temp to 25.428C, 77.7704F
	// using distinct number so it's obvious when temp sensor is not available/working
	var tempRaw int64 = 25428
	var logOutput string = "hydroPi.log"
	var deviceIdFile = ".hydroPiId"

	//get user information
	usr, err := user.Current()
	if err != nil {
		Info.Println(err, " - unable to get current user information, using current directory for file io")
	} else {
		logOutput = usr.HomeDir + "/hydroPi.log"
		//store uuid in users home directory
		deviceIdFile = usr.HomeDir + "/.hydroPiId"
	}

	//setup file for logging
	initLogging(logOutput)

	//generate uuid
	deviceId := uuid.NewRandom().String()

	// if the device id already exists on the device
	if _, err := os.Stat(deviceIdFile); err == nil {
		// open the file containing the id
		f, err := os.Open(deviceIdFile)
		if err != nil {
			Info.Println(err)
		}

		b := make([]byte, 36)
		_, err = f.Read(b)
		if err != nil {
			Info.Println(err)
		}

		deviceId = string(b)

		Info.Println("Starting up using existing device id:", deviceId)

		// if the device id does not exists
	} else {
		f, err := os.Create(deviceIdFile)
		if err != nil {
			Error.Println(err, " - unable to create file to store device id")
		} else {
			f.WriteString(deviceId)
		}
		Info.Println("Starting up with new device id:", deviceId)
		Info.Println("Saving id to", deviceIdFile)
	}

	// open i2c communication to ph
	phChip, err := atlas.New("ph")
	if err != nil {
		Error.Println(err)
	}

	// open i2c communication to ec
	ecChip, err := atlas.New("ec")
	if err != nil {
		Error.Println(err)
	}

	stats, err := statsd.New(
		statsd.Address("159.203.144.95:8125"),
		statsd.Prefix("aquaponics"),
	)
	if err != nil {
		Error.Println(err)
		return
	}
	defer stats.Close()

	// Forever
	for {
		// open 1-wire communication to temp sensor
		oneWire, err := onewire.NewDS18S20("28-031466321eff")
		if err != nil {
			Info.Println("Temperature sensor not available, using default or temp from last reading if available")
		} else {
			tempBuf, err := oneWire.Read()
			if err != nil {
				Info.Println("Temperature sensor not available, using default or temp from last reading if available")
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
			Error.Println(err)
		}

		// update ecChip.reading
		ecChip.UpdateReading(byteTemp)
		if err != nil {
			Error.Println(err)
		}

		// access new readings
		phReading := phChip.GetReading()
		tdsReading := ecChip.GetReading()

		Debug.Printf("temp '%v', ph'%v', tds '%v'\n", tempF, phReading, tdsReading)

		// send to statsd
		stats.Gauge(deviceId+".watertempf", tempF)
		stats.Gauge(deviceId+".watertempc", tempC)
		stats.Gauge(deviceId+".tds", tdsReading)
		stats.Gauge(deviceId+".ph", phReading)
	}
}
