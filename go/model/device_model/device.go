package device_model

import (
	"time"
)

type DeviceAction int

const (
	Bind DeviceAction = iota
	Query
)

type Device struct {
	DeviceID     int64     `json:"device_id"`
	DevicesUID   string    `json:"devices_uid"`
	UserID       *int64    `json:"user_id"`
	DeviceStatus string    `json:"device_status"`
	LastOnline   time.Time `json:"last_online_time"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UpdateDeviceUserIDRequest struct {
	DeviceUID  string  `json:"device_uid" binding:"required"`
	DeviceName *string `json:"device_name" binding:"required"`
	UserID     *int64  `json:"user_id" binding:"required"`
}
type UpdateDeviceUserResponse struct {
	DeviceUID string `json:"device_uid"`
	Message   string `json:"message"`
}
type DeviceInfoResponse struct {
	DeviceID     int64     `json:"device_id"`
	DevicesUID   string    `json:"devices_uid"`
	UserID       *int64    `json:"user_id"`
	DeviceName   *string   `json:"device_name"`
	DeviceStatus string    `json:"device_status"`
	LastOnline   time.Time `json:"last_online_time"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
type GetDeviceInfoResponse struct {
	TotalCount int                  `json:"total_count"`
	Devices    []DeviceInfoResponse `json:"devices"`
	Message    string               `json:"message"`
}
