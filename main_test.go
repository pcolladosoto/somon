package main

import (
	"fmt"
	"log/slog"
	"os"
	"testing"
)

func init() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug, AddSource: false})))
}

func TestExtractData(t *testing.T) {
	rawPayload, err := os.ReadFile("testdata/raw.json")
	if err != nil {
		t.Errorf("couldn't open the test file: %v", err)
	}

	dataPoint, err := extractData(rawPayload)
	if err != nil {
		t.Errorf("couldn't extract the data: %v", err)
	}
	fmt.Printf("extracted data: %+v\n", dataPoint)
}

func TestInsertDataPoint(t *testing.T) {
	db, err := openDB()
	if err != nil {
		t.Errorf("error opening the db: %v", err)
		t.FailNow()
	}
	defer db.Close()

	insertDataPoint(db, extractedValue{
		IMEI: "101010101010",
		Values: map[int]map[parameter]float32{
			0: {
				temperature:  10,
				humidity:     10,
				conductivity: 10,
			},
			1: {
				temperature:  20,
				humidity:     20,
				conductivity: 20,
			},
			2: {
				temperature:  30,
				humidity:     30,
				conductivity: 30,
			},
			3: {
				temperature:  40,
				humidity:     40,
				conductivity: 40,
			},
		},
	})
}
