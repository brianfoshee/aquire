package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/brianfoshee/aquire/atlas"
	"github.com/brianfoshee/raspberrypi/onewire"
)

func main() {

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

		// jot down the time
		isoDateTime := time.Now().UTC().Format(time.RFC3339)

		// create data structure dictated by server
		rdgns := map[string]interface{}{
			"created_at": isoDateTime,
			"sensor_data": map[string]float64{
				"water_temperature": tempF,
				"tds":               tdsReading,
				"ph":                phReading,
			},
		}

		// json enocde structure
		b, err := json.Marshal(rdgns)
		if err != nil {
			fmt.Println(err)
		}
		buf := bytes.NewBuffer(b)

		// post readings to server
		resp, err := http.Post("http://gowebz.herokuapp.com/devices/MockClient1/readings", "application/json", buf)
		if err != nil {
			fmt.Println(err)
		}

		// Log Post and response
		fmt.Println("Data: ", rdgns)
		fmt.Println("Response: ", resp)
	}
}
