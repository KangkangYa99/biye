package devicedata_services

import (
	"biye/DeviceData/devicedata_repository"
	"biye/Devices/device_repository"
	"biye/model/devicedata_model"
	"biye/share/error_code"
	"context"
)

type DeviceDataServices struct {
	data   devicedata_repository.DataInterFace
	device device_repository.DeviceInterface
}

func NewDeviceDataServices(req devicedata_repository.DataInterFace, device device_repository.DeviceInterface) *DeviceDataServices {
	return &DeviceDataServices{
		data:   req,
		device: device,
	}
}
func (d *DeviceDataServices) Insert(ctx context.Context, data *devicedata_model.SendDataReportRequest) (*devicedata_model.SensorDataReportResponse, error) {
	deviceID, err := d.data.GetDeviceIDByUID(ctx, data.DeviceUID)
	if err != nil {
		return nil, error_code.DeviceNotFound.WithDetail("设备不存在")
	}
	data.DeviceID = int64(*deviceID)

	err = d.data.InsertData(ctx, data)
	if err != nil {
		return nil, err
	}
	return &devicedata_model.SensorDataReportResponse{
		Message: "数据接收成功",
	}, nil
}
func (d *DeviceDataServices) GetDeviceHistoryData(ctx context.Context, req *devicedata_model.DataHistoryRequest) ([]*devicedata_model.DataHistoryResponse, error) {
	rec, err := d.data.GetDataHistory(ctx, req)
	if err != nil {
		return nil, err // Repository 应该已经包装了错误
	}
	return rec, nil
}
