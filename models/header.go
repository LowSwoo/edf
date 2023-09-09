package models

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Header struct {
	Version                      string
	LPI                          string
	LRI                          string
	StartDate                    string
	StartTime                    string
	NumbersOfBytesInHeader       string
	Reserved1                    string
	NumberOfDataRecords          string
	DurationOfDataRecordInSecond string
	NumberOfSignals              string
	Label                        string
	TransducerType               string
	PhysicalDimension            string
	PhysicalMinimum              string
	PhysicalMaximum              string
	DigitalMinimum               string
	DigitalMaximum               string
	Prefiltering                 string
	NrNs                         string
	Reserved2                    string
}

type EdfParser struct {
	Header
	Body [][]int16
}

func NewEdfParser(filePath string) (*EdfParser, error) {
	rawBytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var counter = 0
	var nr int
	var ns int

	read := func(count int) []byte {
		counter += count
		return rawBytes[counter-count : counter]
	}

	parseHeader := func() Header {
		var header Header

		header.Version = string(read(8))

		header.LPI = string(read(80))
		header.LRI = string(read(80))

		header.StartDate = string(read(8))
		header.StartTime = string(read(8))

		header.NumbersOfBytesInHeader = string(read(8))
		header.Reserved1 = string(read(44))

		header.NumberOfDataRecords = string(read(8))
		header.DurationOfDataRecordInSecond = string(read(8))
		header.NumberOfSignals = string(read(4))

		ns, err = strconv.Atoi(strings.Split(header.NumberOfSignals, " ")[0])
		if err != nil {
			panic(err)
		}

		header.Label = string(read(ns * 16))
		header.TransducerType = string(read(ns * 80))
		header.PhysicalDimension = string(read(ns * 8))

		header.PhysicalMinimum = string(read(ns * 8))
		header.PhysicalMaximum = string(read(ns * 8))

		header.DigitalMinimum = string(read(ns * 8))
		header.DigitalMaximum = string(read(ns * 8))

		header.Prefiltering = string(read(ns * 80))
		header.NrNs = string(read(ns * 8))

		nr, err = strconv.Atoi(strings.Split(header.NrNs, " ")[0])
		if err != nil {
			panic(err)
		}
		nr /= ns

		header.Reserved2 = string(read(ns * 32))

		return header
	}

	parseBody := func() [][]int16 {
		read(1)
		getOneRecord := func() []int16 {

			var record = make([]int16, nr)

			for i := 0; i < nr; i++ {
				// record[i] = int(binary.BigEndian.Uint16(read(2)))
				// record[i], _ = strconv.Atoi(fmt.Sprintf("%d", read(2)))

				// tmp := hex.EncodeToString(read(2))
				tmp, _ := strconv.ParseInt(fmt.Sprintf("%x", read(2)), 16, 32)
				record[i] = int16(tmp)
				fmt.Printf("%d\t", record[i])
			}

			return record
		}

		var records = make([][]int16, ns)

		for i := 0; i < ns; i++ {
			records[i] = getOneRecord()
		}

		return records
	}

	header := parseHeader()
	body := parseBody()

	return &EdfParser{
		Header: header,
		Body:   body,
	}, nil
}
