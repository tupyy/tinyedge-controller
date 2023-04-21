package postgres

import (
	"github.com/tupyy/tinyedge-controller/internal/entity"
	"gorm.io/gorm"
)

type Options interface {
	Sort(query *gorm.DB) *gorm.DB
	Filter(query *gorm.DB) *gorm.DB
}

type FilterByKind struct {
	Kind entity.ManifestKind
}

func (f *FilterByKind) Filter(query *gorm.DB) *gorm.DB {
	refType := "configuration"
	if f.Kind == entity.WorkloadManifestKind {
		refType = "workload"
	}
	return query.Where("ref_type = ?", refType)
}

func (f *FilterByKind) Sort(query *gorm.DB) *gorm.DB {
	return query // no-op
}
