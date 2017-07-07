package models

import (
	"github.com/antipin1987@gmail.com/rpcj/db"
	"github.com/antipin1987@gmail.com/rpcj/models/orm"
	"fmt"
	"time"
)


type User struct {

}

//Структура-фильтр для поиска пользователей (методы Add и Delete)
type Filter struct {
	UUId string `json:"uuid"`
	Login string `json:"login"`
	Date time.Time `json:"date"`
}

func (f *Filter) ToORM() *orm.User {
	return &orm.User{
		Id: f.UUId,
		Login: f.Login,
		CreatedAt: f.Date,
	}
}

//Метод добавления пользователя с данным логином (uuid и дата регистрации проставляются автоматически)
func (u *User) Add(login string, reply *bool) error {
	if err := db.DB().Create(&orm.User{Login: login}).Error; err != nil {
		*reply = false
		return err
	}
	*reply = true
	return nil
}

//Метод поиска пользователей по фильтру (по полям uuid, логин, дата регистрации, могут быть заданы произвольно -  вместе или по отдельности)
func (u *User) Find(f *Filter, usrs *[]orm.User) error {
	var users []orm.User
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
	for _, u := range users {
		*usrs = append(*usrs, u)
	}
	return nil
}

//Метод удаления пользователей по фильтру
func (u *User) Delete(f *Filter, reply *string) error {
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
	resp := q.Delete(&orm.User{})
	if resp.Error != nil {
		return resp.Error
	}
	*reply = fmt.Sprintf("%d rows deleted", resp.RowsAffected)
	return nil
}

func init() {
	register("User", &User{})
}