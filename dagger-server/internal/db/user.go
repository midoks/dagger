package db

import (
	"errors"
	"time"
)

type User struct {
	Id       int64  `gorm:"primaryKey"`
	Name     string `gorm:"unique;size:50;comment:用户名"`
	Password string `gorm:"size:32;comment:用户密码"`
	Status   int    `gorm:"comment:状态"`

	Created     time.Time `gorm:"autoCreateTime;comment:创建时间"`
	CreatedUnix int64     `gorm:"autoCreateTime;comment:创建时间"`
	Updated     time.Time `gorm:"autoCreateTime;comment:更新时间"`
	UpdatedUnix int64     `gorm:"autoCreateTime;comment:更新时间"`
}

func (User) TableName() string {
	return "dd_users"
}

func UserGetByName(name string) (User, error) {
	var user User
	err := db.First(&user, "name = ?", name).Error
	if err != nil {
		return user, err
	}
	return user, nil
}

func UsersList() ([]*User, error) {
	page := 1
	pageSize := 1000
	u := make([]*User, 0, pageSize)
	err := db.Limit(pageSize).Offset((page - 1) * pageSize).Order("id desc").Find(&u)
	return u, err.Error
}

func UserAdd(name, password string) error {
	var u User

	err := db.First(&u, "name = ?", name).Error
	if err == nil {
		return errors.New("already user exists!")
	}

	u.Password = password
	u.Name = name

	result := db.Create(&u)
	return result.Error
}

func UserDel(name string) error {
	err := db.Where("name = ?", name).Delete(&User{})
	return err.Error
}

func UserMod(name, password string) error {
	var u User
	err := db.First(&u, "name = ?", name).Error
	if err != nil {
		return errors.New("user not exists!")
	}

	err = db.Model(&User{}).Where("name = ?", name).Update("password", password).Error
	return err
}

func UserAclCheck(name, password string) bool {
	u, err := UserGetByName(name)
	if err != nil {
		return false
	}

	if u.Password != password {
		return false
	}

	return true
}
