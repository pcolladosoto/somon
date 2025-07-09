package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log/slog"
)

type parameter int

const (
	temperature parameter = iota
	humidity
	conductivity
)

var parameterTypeMap = map[parameter]string{
	temperature:  "temperature",
	humidity:     "humidity",
	conductivity: "conductivity",
}

func (p parameter) String() string {
	s, ok := parameterTypeMap[p]
	if !ok {
		return "unknown"
	}
	return s
}

type OnceDecPayload struct {
	Received string `json:"received"`
	Id       string `json:"id"`
	Source   string `json:"source"`
	Type     string `json:"type"`
	Version  string `json:"version"`
	Device   struct {
		ICCID string `json:"iccid"`
		Ipv4  string `json:"ip"`
		IMSI  string `json:"imsi"`
	} `json:"device"`
	Payload struct {
		Encoding string `json:"encoding"`
		Type     string `json:"type"`
		Value    string `json:"value"`
	} `json:"payload"`
}

type SE0XNBDecPayload struct {
	IMEI  string `json:"IMEI"`
	IMSI  string `json:"IMSI"`
	Model string `json:"Model"`

	Time string `json:"time"`

	Interrupt      int `json:"interrupt"`
	InterruptLevel int `json:"interrupt_level"`

	Battery float32 `json:"battery"`
	Signal  float32 `json:"signal"`

	SensorFlag string `json:"sensor_flag"`

	WaterA        float32 `json:"water_soil1"`
	TemperatureA  float32 `json:"temp_soil1"`
	ConductivityA float32 `json:"conduct_soil1"`

	WaterB        float32 `json:"water_soil2"`
	TemperatureB  float32 `json:"temp_soil2"`
	ConductivityB float32 `json:"conduct_soil2"`

	WaterC        float32 `json:"water_soil3"`
	TemperatureC  float32 `json:"temp_soil3"`
	ConductivityC float32 `json:"conduct_soil3"`

	WaterD        float32 `json:"water_soil4"`
	TemperatureD  float32 `json:"temp_soil4"`
	ConductivityD float32 `json:"conduct_soil4"`
}

type extractedValue struct {
	IMEI   string
	Values map[int]map[parameter]float32
}

func parsePayload(body []byte) (interface{}, error) {
	tmpPayload := map[string]interface{}{}
	if err := json.Unmarshal(body, &tmpPayload); err != nil {
		return nil, fmt.Errorf("couldn't unmarshal into temporary struct")
	}

	// Try to parse a 1nce payload
	payloadType, ok := tmpPayload["type"]
	if ok {
		payloadTypeStr, ok := payloadType.(string)
		if !ok || payloadTypeStr != "TELEMETRY_DATA" {
			return nil, fmt.Errorf("malformed 1nce payload")
		}
		payload := OnceDecPayload{}
		if err := json.Unmarshal(body, &payload); err != nil {
			return nil, fmt.Errorf("couldn't unmarshal into 1nce struct: %w", err)
		}
		return payload, nil
	}

	return nil, fmt.Errorf("couldn't detect struct type")
}

func parse1ncePayload(rawPayload string) (extractedValue, error) {
	decPayload, err := base64.StdEncoding.DecodeString(rawPayload)
	if err != nil {
		return extractedValue{}, fmt.Errorf("couldn't decode the bas64-encoded payload: %v", err)
	}

	tmpPayload := map[string]interface{}{}
	if err := json.Unmarshal(decPayload, &tmpPayload); err != nil {
		return extractedValue{}, fmt.Errorf("couldn't unmarshal into 1nce struct: %v", err)
	}
	sensorModel, ok := tmpPayload["Model"]
	if !ok {
		return extractedValue{}, fmt.Errorf("no 'Model' field in the payload")
	}

	sensorModelStr, ok := sensorModel.(string)
	if !ok {
		return extractedValue{}, fmt.Errorf("the 'Model' field was not a string")
	}

	dataPoints := map[int]map[parameter]float32{}
	switch sensorModelStr {
	case "SE0X-NB":
		slog.Debug("decoded a SE0X-NB payload")
		se0XPayload := SE0XNBDecPayload{}
		if err := json.Unmarshal(decPayload, &se0XPayload); err != nil {
			return extractedValue{}, fmt.Errorf("couldn't unmarshal into SE0X-NB struct: %v", err)
		}

		for i, present := range se0XPayload.SensorFlag {
			if present == '1' {
				slog.Debug("got data for a sensor", "imei", se0XPayload.IMEI, "sensor", i)
				dataPoints[i] = map[parameter]float32{}

				switch i {
				case 0:
					dataPoints[0][temperature] = se0XPayload.TemperatureA
					dataPoints[0][humidity] = se0XPayload.WaterA
					dataPoints[0][conductivity] = se0XPayload.ConductivityA

				case 1:
					dataPoints[1][temperature] = se0XPayload.TemperatureB
					dataPoints[1][humidity] = se0XPayload.WaterB
					dataPoints[1][conductivity] = se0XPayload.ConductivityB

				case 2:
					dataPoints[2][temperature] = se0XPayload.TemperatureC
					dataPoints[2][humidity] = se0XPayload.WaterC
					dataPoints[2][conductivity] = se0XPayload.ConductivityC

				case 3:
					dataPoints[3][temperature] = se0XPayload.TemperatureD
					dataPoints[3][humidity] = se0XPayload.WaterD
					dataPoints[3][conductivity] = se0XPayload.ConductivityD
				}
			}
		}

		return extractedValue{IMEI: se0XPayload.IMEI, Values: dataPoints}, nil
	default:
		return extractedValue{}, fmt.Errorf("couldn't detect any payload type for %q", sensorModelStr)
	}
}

func extractData(body []byte) (extractedValue, error) {
	parsedPayload, err := parsePayload(body)
	if err != nil {
		return extractedValue{}, fmt.Errorf("couldn't parse the incoming payload: %v", err)
	}

	switch decPayload := parsedPayload.(type) {
	case OnceDecPayload:
		slog.Debug("decoded a 1nce payload", "payload", decPayload)
		return parse1ncePayload(decPayload.Payload.Value)
	default:
		return extractedValue{}, fmt.Errorf("got unidentified payload type: %#v", parsedPayload)
	}
}
