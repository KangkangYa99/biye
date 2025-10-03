package devicedata_repository

import (
	"biye/model/devicedata_model"
	"biye/share/error_code"
	"context"
	"database/sql"
	time2 "time"
)

type DeviceDataRepository struct {
	db *sql.DB
}

func NewDeviceDataRepository(db *sql.DB) *DeviceDataRepository {
	return &DeviceDataRepository{
		db: db,
	}
}
func (R *DeviceDataRepository) InsertData(ctx context.Context, data *devicedata_model.SendDataReportRequest) error {
	time := time2.Unix(data.DataTimeStamp, 0)
	query := `INSERT INTO device_data (device_id,temperature,humidity,light,noise,fire,co,light_on,fan_on,data_timestamp)
			VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`
	_, err := R.db.ExecContext(ctx, query,
		data.DeviceID,
		data.Temperature,
		data.Humidity,
		data.Light,
		data.Noise,
		data.Fire,
		data.CO,
		data.LightStatus,
		data.FanStatus,
		time,
	)
	if err != nil {

		return error_code.DatabaseError.WithDetail(err.Error())
	}
	return nil
}
func (R *DeviceDataRepository) GetDeviceIDByUID(ctx context.Context, uid string) (*int, error) {
	query := `SELECT device_id FROM devices WHERE device_uid = $1`
	var deviceID *int
	err := R.db.QueryRowContext(ctx, query, uid).Scan(&deviceID)
	if err != nil {
		return nil, error_code.DeviceNotFound.WithDetail(err.Error())
	}
	return deviceID, nil
}
