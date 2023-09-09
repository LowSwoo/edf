package main

import (
	"edf_reader/models"
	"encoding/json"
	"os"
)

func main() {

	var sampleName = "samples/OXY_R_AC.edf"

	edf, err := models.NewEdfParser(sampleName)
	if err != nil {
		panic(err)
	}

	file, err := os.Create("data.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	json.NewEncoder(file).Encode(edf)

}
