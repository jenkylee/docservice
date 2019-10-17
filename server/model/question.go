// Package model 定义数据库相关模型
// 试题模型定义

package model

import "time"

type Question struct {
	ID            string    `gorm:"primary_key;size:128;not null;"`
	Type          string    `gorm:"size:32;not null;"`
	Title         string    `gorm:"type:text;not null;"`
	Option        string    `gorm:"type:text;not null;"`
	Answer        string    `gorm:"type:varchar(255);"`
	Mark          int
	Subject       string    `gorm:"type:varchar(128);"`
	Difficulty    string    `gorm:"type:varchar(255);"`
	Analysis      string    `gorm:"type:text;"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}