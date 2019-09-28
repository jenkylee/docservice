// Package repository 定义数据库操作相关
// 用户创建、更新、删除等操作
package repository

import (
	"github.com/jinzhu/gorm"

	"yokitalk.com/docservice/server/model"
)

type UserRepository interface {
	Find(id int32) (*model.User, error)
	Create(*model.User) error
	Update(*model.User, int32) (*model.User, error)
	FindByField(string, string, string) (*model.User, error)
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userDBRepository{db: db}
}

type userDBRepository struct {
	db *gorm.DB
}

func(this *userDBRepository) Find(id int32)(*model.User, error){
	user := &model.User{}
	user.ID = uint(id)

	if err := this.db.First(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func(this *userDBRepository) Create(user *model.User) error{
	if err := this.db.Create(user).Error; err != nil {
		return err
	}

	return nil
}

func (this *userDBRepository) Update(user *model.User, id int32) (*model.User, error) {
	user.ID = uint(id)
	if err := this.db.Model(user).Updates(&user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (this *userDBRepository) FindByField(key string, value string, fields string) (*model.User, error) {
	if len(fields) == 0 {
		fields = "*"
	}
	user :=  &model.User{}
	if err := this.db.Select(fields).Where(key+" = ?", value).First(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}


