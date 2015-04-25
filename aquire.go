package main

import (
	"github.com/crakalakin/atlas"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"github.com/raspberrypi/onewire"
)

func main() {
	oneWire, err := onewire.NewDS18S20("28-031466321eff")
	if err != nil {
		fmt.Print(err)
	}

	phChip, err := atlas.New("ph")
	if err != nil {
		fmt.Println("Unable to open I2C Device")
	}
	ecChip, err := atlas.New("ec")
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

		phChip.UpdateReading(byteTemp)
		ecChip.UpdateReading(byteTemp)
		phReading := phChip.GetReading()
		tdsReading := ecChip.GetReading()

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

		buf := bytes.NewBuffer(b)
		resp, err := http.Post("http://gowebz.herokuapp.com/devices/MockClient1/readings", "application/json", buf)
		if err != nil {
			fmt.Println("Error posting data: ", err)
		}
		fmt.Println("Response: ", resp)
	}
}
