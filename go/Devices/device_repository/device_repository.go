package device_repository

import (
	"biye/model/device_model"
	"biye/share/error_code"
	"context"
	"database/sql"
	"errors"
)

type DeviceRepository struct {
	db *sql.DB
}

func NewDeviceRepository(db *sql.DB) *DeviceRepository {
	return &DeviceRepository{
		db: db,
	}
}
func (R *DeviceRepository) CreateDevice(ctx context.Context, req *device_model.Device) error {
	query := `INSERT devices(device_uid,device_status,last_online,created_at,updated_at)
	VALUES($1,$2,$3,$4,$5)
	RETURNING device_id,created_at,updated_at`
	err := R.db.QueryRowContext(ctx, query,
		req.DevicesUID,
		req.DeviceStatus,
		req.LastOnline,
		req.CreatedAt,
		req.UpdatedAt).Scan(&req.DeviceID, &req.DeviceStatus, &req.LastOnline, &req.CreatedAt, &req.UpdatedAt)
	if err != nil {
		return error_code.DatabaseError
	}
	return nil
}
func (R *DeviceRepository) UpdateDeviceUserID(ctx context.Context, req *device_model.UpdateDeviceUserIDRequest) error {

	query := `UPDATE devices SET user_id = $1,device_name = $2 WHERE device_uid = $3 `
	_, err := R.db.ExecContext(ctx, query, req.UserID, req.DeviceName, req.DeviceUID)
	if err != nil {
		return error_code.DatabaseError
	}
	return nil
}
func (R *DeviceRepository) GetUserIDByDeviceUID(ctx context.Context, deviceUID string) (*int64, error) {
	var userID *int64
	query := `SELECT user_id FROM devices WHERE device_uid = $1`
	err := R.db.QueryRowContext(ctx, query, deviceUID).Scan(
		&userID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, error_code.DeviceNotFound
		}
		return nil, error_code.DatabaseError
	}
	return userID, nil

}
func (R *DeviceRepository) GetDeviceByID(ctx context.Context, userId int64) ([]device_model.DeviceInfoResponse, error) {
	query := `SELECT device_id, device_uid,device_name,device_status, last_online, created_at, updated_at
              FROM devices WHERE user_id = $1`

	rows, err := R.db.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, error_code.DatabaseError
	}
	var devices []device_model.DeviceInfoResponse
	for rows.Next() {
		var device device_model.DeviceInfoResponse
		err := rows.Scan(
			&device.DeviceID,
			&device.DevicesUID,
			&device.DeviceName,
			&device.DeviceStatus,
			&device.LastOnline,
			&device.CreatedAt,
			&device.UpdatedAt,
		)
		if err != nil {
			return nil, error_code.DatabaseError
		}
		devices = append(devices, device)
	}
	if err = rows.Err(); err != nil {
		return nil, error_code.DatabaseError
	}
	return devices, nil
}
