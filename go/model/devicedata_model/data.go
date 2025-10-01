package devicedata_model

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
