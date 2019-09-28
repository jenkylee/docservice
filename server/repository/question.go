// Package repository 定义数据库操作相关
// 试题创建、更新、删除等操作
package repository

import (
	"github.com/jinzhu/gorm"
	"yokitalk.com/docservice/server/model"
)

type QuestionRepository interface {
	Find(id string) (*model.Question, error)
	Create(*model.Question) error
	Update(*model.Question, string) (*model.Question, error)
	FindByField(string, string, string) (*model.Question, error)
}

func NewQuestionRepository(db *gorm.DB) QuestionRepository {
	return &questionDBRepository{db: db}
}

type questionDBRepository struct {
	db *gorm.DB
}

func (this *questionDBRepository) Find(id string) (*model.Question, error) {
	question := &model.Question{}
	question.ID = id

	if err := this.db.First(question).Error; err != nil {
		return nil, err
	}

	return question, nil
}

func (this *questionDBRepository) Create(question *model.Question) error {
	if err := this.db.Create(question).Error; err != nil {
		return err
	}

	return nil
}

func (this *questionDBRepository) Update(question *model.Question, id string) (*model.Question, error) {
	question.ID = id
	if err := this.db.Model(question).Updates(&question).Error; err != nil {
		return nil, err
	}

	return question, nil
}

func (this *questionDBRepository) FindByField(key string, value string, fields string) (*model.Question, error) {
	if len(fields) == 0 {
		fields = "*"
	}

	question := &model.Question{}
	if err := this.db.Select(fields).Where(key + " = ?", value).First(question).Error; err != nil {
		return nil, err
	}

	return question, nil
}