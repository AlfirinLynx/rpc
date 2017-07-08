package orm

import (
	"time"
	"github.com/jinzhu/gorm"
	"github.com/pborman/uuid"
)

type User struct {
	Id string `gorm:"column:uuid"`
	Login string
	CreatedAt time.Time `gorm:"column:registration_date"`
}

func (User) TableName() string {
	return "users"
}


func (user *User) BeforeCreate(scope *gorm.Scope) error {
	if user.Id == "" {
		scope.SetColumn("Id", uuid.New())
	}
	return nil
}

