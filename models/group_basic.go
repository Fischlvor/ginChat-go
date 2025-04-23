package models

import "gorm.io/gorm"

type GroupBasic struct {
	gorm.Model
	Name    string `json:"name" gorm:"type:varchar(100);not null"` // 群组名称
	OwnerId uint   `json:"owner_id" gorm:"not null"`               // 群组拥有者ID
	Icon    string `json:"icon" gorm:"type:varchar(255)"`          // 群组图标
	Type    int    `json:"type" gorm:"type:int"`                   // 群组类型
	Desc    string `json:"desc" gorm:"type:varchar(255)"`          // 群组描述
}

func (table *GroupBasic) TableName() string {
	return "group_basic"
}
