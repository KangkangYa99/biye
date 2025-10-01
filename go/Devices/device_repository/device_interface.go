package device_repository

import (
	"biye/model/device_model"
	"context"
)

type DeviceInterface interface {
	CreateDevice(ctx context.Context, req *device_model.Device) error
	UpdateDeviceUserID(ctx context.Context, req *device_model.UpdateDeviceUserIDRequest) error
	GetUserIDByDeviceUID(ctx context.Context, deviceUID string) (*int64, error)
	GetDeviceByID(ctx context.Context, userId int64) ([]device_model.DeviceInfoResponse, error)
}
