package devicedata_services

import (
	"biye/DeviceData/devicedata_repository"
	"biye/Devices/device_repository"
)

type DeviceDataServices struct {
	req    devicedata_repository.DataInterFace
	device device_repository.DeviceInterface
}

func NewDeviceDataServices(req devicedata_repository.DataInterFace, device device_repository.DeviceInterface) *DeviceDataServices {
	return &DeviceDataServices{
		req:    req,
		device: device,
	}
}
