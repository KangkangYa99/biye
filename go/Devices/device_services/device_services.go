package device_services

import (
	"biye/Devices/device_repository"
	"biye/model/device_model"
	"biye/share/error_code"
	"context"
	"fmt"
	"log"
)

type DeviceServices struct {
	deviceRepo device_repository.DeviceInterface
}

func NewDeviceService(deviceRepo device_repository.DeviceInterface) *DeviceServices {
	return &DeviceServices{
		deviceRepo: deviceRepo,
	}
}
func (S *DeviceServices) BindDevices(ctx context.Context, req *device_model.UpdateDeviceUserIDRequest) (*device_model.UpdateDeviceUserResponse, error) {
	err := S.ValidateDevice(ctx, req.DeviceUID, *req.UserID, device_model.Bind)
	if err != nil {
		return nil, err
	}
	err = S.deviceRepo.UpdateDeviceUserID(ctx, req)
	if err != nil {
		return nil, err
	}
	return &device_model.UpdateDeviceUserResponse{
		DeviceUID: req.DeviceUID,
		Message:   "绑定成功。",
	}, nil
}
func (S *DeviceServices) UnBindDevices(ctx context.Context, req *device_model.UpdateDeviceUserIDRequest) (*device_model.UpdateDeviceUserResponse, error) {
	err := S.ValidateDevice(ctx, req.DeviceUID, *req.UserID, device_model.Query)
	if err != nil {

		log.Printf(err.Error())
		return nil, err
	}
	req.UserID = nil
	err = S.deviceRepo.UpdateDeviceUserID(ctx, req)
	if err != nil {
		return nil, err

	}
	return &device_model.UpdateDeviceUserResponse{
		DeviceUID: req.DeviceUID,
		Message:   "解绑成功。",
	}, nil
}
func (S *DeviceServices) ValidateDevice(ctx context.Context, deviceUID string, userID int64, action device_model.DeviceAction) error {
	ID, err := S.deviceRepo.GetUserIDByDeviceUID(ctx, deviceUID)
	if err != nil {
		return err
	}
	switch action {
	case device_model.Bind:
		if ID != nil {
			return error_code.DeviceIsBind.WithDetail(fmt.Sprintf("设备: %v 已被绑定。", deviceUID))
		}
	case device_model.Query:
		if ID == nil {
			return error_code.DeviceNotBind.WithDetail(fmt.Sprintf("设备: %v 未被绑定。", deviceUID))
		}
		if *ID != userID {
			return error_code.NotDeviceOwner.WithDetail(fmt.Sprintf("您不是设备 %v 的拥有者", deviceUID))
		}
	}
	return nil
}
func (S *DeviceServices) GetDevicesByUserID(ctx context.Context, userID int64) (*device_model.GetDeviceInfoResponse, error) {

	devices, err := S.deviceRepo.GetDeviceByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	response := &device_model.GetDeviceInfoResponse{
		TotalCount: len(devices),
		Devices:    devices, // 返回设备列表
		Message:    fmt.Sprintf("获取到 %d 个设备", len(devices)),
	}
	return response, nil
}
