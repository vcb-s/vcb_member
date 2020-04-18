package models

import "time"

// SoftDeletedModel 设置软删键值
type SoftDeletedModel struct {
	// 这里统一omit掉，不对外输出。需要输出的时候再覆盖定义
	CreatedAt time.Time `json:"-" form:"-" gorm:"column:create_at"`
	UpdatedAt time.Time `json:"-" form:"-" gorm:"column:update_at"`
	DeletedAt time.Time `json:"-" form:"-" gorm:"column:deleted_at"`
}
