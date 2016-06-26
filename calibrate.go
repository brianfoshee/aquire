package main

import (
	"fmt"
	"github.com/brianfoshee/aquire/atlas"
)

func getProbe() int {
	var probe int

	fmt.Println("Please select the type of probe you want to calibrate")
	fmt.Println("1. PH")
	fmt.Println("2. EC")
	fmt.Println("3. DO")
	fmt.Println("4. ORP")
	fmt.Print("Probe number: ")

	_, _ = fmt.Scanln(&probe)

	// While user input is not valid, continue to ask
	for probe < 1 || probe > 4 {
		fmt.Println("")
		fmt.Println("Invalid input, please input the number that corresponds with the probe you are using")
		fmt.Println("For example, if you are calibrating a PH probe enter '1'")
		fmt.Println("1. PH")
		fmt.Println("2. EC")
		fmt.Println("3. DO")
		fmt.Println("4. ORP")
		fmt.Print("Probe number: ")
		_, _ = fmt.Scanln(&probe)
	}

	return probe
}

func main() {
	var dummy string
	var calSolution float64

	probe := getProbe()

	switch probe {
	// PH
	case 1:
		fmt.Print("Enter PH Calibration Solution: ")
		fmt.Scanf("%f", &calSolution)
		for calSolution < 0 || calSolution > 14 {
			fmt.Println("Invalid input, valid PH range is 0-14 inclusive")
			fmt.Print("Enter Calibration Solution: ")
			fmt.Scanf("%f", &calSolution)
		}

		phChip, err := atlas.New("ph")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Print("Submerge PH probe in solution for 30 seconds then press any key to calibrate")
		_, _ = fmt.Scanln(&dummy)
		fmt.Println("Calibrating...")
		err = phChip.Calibrate(calSolution)
		if err != nil {
			fmt.Println("Calibration unsuccessfull:", err)
		} else {
			fmt.Println("Calibration Successfull")
		}
	// EC
	case 2:
		fmt.Println("Calibration not yet available for EC probe")
	// DO
	case 3:
		fmt.Println("Calibration not yet available for DO probe")
	// ORP
	case 4:
		fmt.Println("Calibration not yet available for ORP probe")
	}

}
