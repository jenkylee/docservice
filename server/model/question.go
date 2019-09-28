// Package model 定义数据库相关模型
// 试题模型定义

package model

import "time"

type Question struct {
	ID        string  `gorm:"primary_key;size:128;not null;"`
	Type      string  `gorm:"size:32;not null;"`
	Content   string  `gorm:"type:text;not null;"`
	Answer    string  `gorm:"type:varchar(255);"`
	Mark      int
	Class     string  `gorm:"type:varchar(128);"`
	Tag       string  `gorm:"type:varchar(255);"`
	Analysis      string   `gorm:"type:text;"`
	BlankDisorder int `gorm:"type:tinyint(1)"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}