package devicedata_model

import "time"

type SendDataReportRequest struct {
	DeviceUID     string  `json:"deviceUID"`
	DeviceID      int64   `json:"deviceId"`
	Temperature   float64 `json:"temperature"`
	Humidity      float64 `json:"humidity"`
	Light         float64 `json:"light"`
	Noise         float64 `json:"noise"`
	CO            float64 `json:"co"`
	Fire          bool    `json:"fire"`
	LightStatus   bool    `json:"light_Status"`
	FanStatus     bool    `json:"fan_Status"`
	DataTimeStamp int64   `json:"dataTimeStamp"`
}

type SensorDataReportResponse struct {
	Message string `json:"message"`
}

type DataHistoryRequest struct {
	DeviceUID string `json:"device_uid" binding:"required"`
	Limit     int    `json:"limit"`
}
type DataHistoryResponse struct {
	Temperature float64   `json:"temperature"`
	Humidity    float64   `json:"humidity"`
	Light       float64   `json:"light"`
	Noise       float64   `json:"noise"`
	File        bool      `json:"file"`
	Co          float64   `json:"co"`
	Time        time.Time `json:"time"`
}
