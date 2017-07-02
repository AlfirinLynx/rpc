package models

import (
	"github.com/antipin1987@gmail.com/rpcj/db"
	"github.com/antipin1987@gmail.com/rpcj/models/orm"
)


type User struct {
}


func (u *User) Add(login string, reply *bool) error {
	if err := db.DB().Create(&orm.User{Login: login}).Error; err != nil {
		*reply = false
		return err
	}
	*reply = true
	return nil
}


func init() {
	register("User", &User{})
}