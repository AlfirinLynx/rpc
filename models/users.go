package models

import (
	"github.com/antipin1987@gmail.com/rpcj/db"
	"github.com/antipin1987@gmail.com/rpcj/models/orm"
	"fmt"
	"time"
)


type User struct {

}

type Filter struct {
	UUId string `json:"uuid"`
	Login string `json:"login"`
	Date time.Time `json:"date"`
}

func (f *Filter) toORM() *orm.User {
	return &orm.User{
		Id: f.UUId,
		Login: f.Login,
		CreatedAt: f.Date,
	}
}


func (u *User) Add(login string, reply *bool) error {
	if err := db.DB().Create(&orm.User{Login: login}).Error; err != nil {
		*reply = false
		return err
	}
	*reply = true
	return nil
}

func (u *User) Find(f *Filter, usrs *[]orm.User) error {
	var users []orm.User
	fmt.Println(f)
	q :=  db.DB()
	if !f.Date.IsZero(){
		q = q.Where("registration_date = DATE(?)", f.Date.String())
	}
	if f.UUId != "" {
		q = q.Where("uuid = ?", f.UUId)
	}
	if f.Login != "" {
		q = q.Where("login = ?", f.Login)
	}
	if err := q.Find(&users).Error; err != nil {
		fmt.Println(err)
		return err
	}
	//fmt.Println(users, len(users))
	for _, u := range users {
		*usrs = append(*usrs, u)
	}
	return nil
}


func init() {
	register("User", &User{})
}