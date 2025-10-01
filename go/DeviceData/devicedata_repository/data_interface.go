package devicedata_repository

import (
	"biye/model/devicedata_model"
	"context"
)

type DataInterFace interface {
	InsertData(ctx context.Context, data *devicedata_model.SendDataReportRequest) error
}
