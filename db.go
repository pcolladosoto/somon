package main

import (
	"fmt"
	"log/slog"
	"time"

	"database/sql"

	pq "github.com/lib/pq"
)

const sqlStatement string = `
INSERT INTO %s (ts, imei, sensor, measure)
VALUES ($1, $2, $3, $4)
RETURNING imei
`

func openDB() (*sql.DB, error) {
	dbUri, err := getEnvVar("DB_URI")
	if err != nil {
		return nil, fmt.Errorf("couldn't retrieve the db uri: %w", err)
	}
	slog.Debug("db info", "uri", dbUri)

	db, err := sql.Open("postgres", dbUri)
	if err != nil {
		return nil, err
	}

	slog.Debug("correctly opened the db")

	// Be sure to check https://pkg.go.dev/database/sql and
	// https://gorm.io/docs/generic_interface.html for info
	// on the following settings.
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}

func insertDataPoint(db *sql.DB, dataPoint extractedValue) {
	ts := time.Now().UTC()
	imeiCheck := ""
	for sensorId, sensorData := range dataPoint.Values {
		for parameterType, parameterValue := range sensorData {
			slog.Debug("inserting value", "imei", dataPoint.IMEI, "type", parameterType.String(), "val", parameterValue)
			if err := db.QueryRow(
				fmt.Sprintf(sqlStatement, pq.QuoteIdentifier(parameterType.String())),
				ts, dataPoint.IMEI, sensorId, parameterValue,
			).Scan(&imeiCheck); err != nil {
				slog.Error("error inserting a value", "err", err)
				continue
			}
			slog.Debug("inserted value", "imeiCheck", imeiCheck)
		}
	}
}
