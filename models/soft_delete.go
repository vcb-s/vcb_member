package models

import "time"

// SoftDeletedModel 设置软删键值
type SoftDeletedModel struct {
	CreatedAt time.Time `json:"create_at" form:"create_at" gorm:"column:create_at"`
	UpdatedAt time.Time `json:"update_at" form:"update_at" gorm:"column:update_at"`
	DeletedAt time.Time `json:"delete_at" form:"delete_at" gorm:"column:deleted_at"`
}
