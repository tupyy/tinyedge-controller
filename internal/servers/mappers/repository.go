package mappers

import (
	"github.com/tupyy/tinyedge-controller/internal/entity"
	"github.com/tupyy/tinyedge-controller/pkg/grpc/admin"
)

func RepositoryToModel(r entity.Repository) *admin.Repository {
	return &admin.Repository{
		Id:             r.Id,
		Url:            r.Url,
		Branch:         r.Branch,
		LocalPath:      r.LocalPath,
		CurrentHeadSha: r.CurrentHeadSha,
		TargetHeadSha:  r.TargetHeadSha,
		PullPeriod:     int32(r.PullPeriod),
	}
}
