package postgres

import (
	"gorm.io/gorm"
)

type Options interface {
	Sort(query *gorm.DB) *gorm.DB
	Filter(query *gorm.DB) *gorm.DB
}
