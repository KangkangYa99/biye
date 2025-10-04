package devicedata_repository

import (
	"biye/model/devicedata_model"
	"biye/share/error_code"
	"context"
	"database/sql"
	"fmt"
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
func (R *DeviceDataRepository) GetDataHistory(ctx context.Context, req *devicedata_model.DataHistoryRequest) ([]*devicedata_model.DataHistoryResponse, error) {
	deviceID, err := R.GetDeviceIDByUID(ctx, req.DeviceUID)
	if err != nil {
		return nil, err
	}
	if req.Limit <= 0 {
		req.Limit = 20
	}
	if req.Limit > 100 {
		req.Limit = 100
	}
	query := `
		SELECT temperature, humidity, light, noise, fire, co,data_timestamp FROM device_data WHERE device_id = $1 ORDER BY data_timestamp DESC 
		LIMIT $2`
	rows, err := R.db.QueryContext(ctx, query, deviceID, req.Limit)
	if err != nil {
		return nil, error_code.DatabaseError.WithDetail(fmt.Sprintf("查询历史数据失败: %v", err))
	}
	defer rows.Close()
	var records []*devicedata_model.DataHistoryResponse
	for rows.Next() {
		var record devicedata_model.DataHistoryResponse
		err := rows.Scan(
			&record.Temperature,
			&record.Humidity,
			&record.Light,
			&record.Noise,
			&record.File,
			&record.Co,
			&record.Time,
		)
		if err != nil {
			return nil, error_code.DatabaseError.WithDetail(fmt.Sprintf("扫描历史数据失败: %v", err))
		}
		records = append(records, &record)
	}
	if err = rows.Err(); err != nil {
		return nil, error_code.DatabaseError.WithDetail(fmt.Sprintf("读取历史数据时发生错误: %v", err))
	}
	return records, nil
}
