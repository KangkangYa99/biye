package devicedata_repository

import (
	"biye/model/devicedata_model"
	"context"
)

type DataInterFace interface {
	InsertData(ctx context.Context, data *devicedata_model.SendDataReportRequest) error
	GetDeviceIDByUID(ctx context.Context, uid string) (*int, error)
	GetDataHistory(ctx context.Context, req *devicedata_model.DataHistoryRequest) ([]*devicedata_model.DataHistoryResponse, error)
}
