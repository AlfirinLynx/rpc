package models

import (
	"github.com/antipin1987@gmail.com/rpcj/db"
	"github.com/antipin1987@gmail.com/rpcj/models/orm"
	"fmt"
	"time"
	"github.com/jinzhu/gorm"
)

const (
	DATE_LAYOUT = "2006-01-02"  //формат представления даты (для time.Parse)
)

type User struct {

}

//Структура-фильтр для поиска пользователей (методы Add и Delete)
type Filter struct {
	UUId string `json:"uuid"`
	Login string `json:"login"`
	Date string `json:"date"`  //дата регистрации пользователя задается строкой вида "2017-07-07" (ГГГГ-ММ-ДД)
}

//По фильтру подготовить запрос в БД (через ОРМ)
func (f *Filter) GetQuery() (*gorm.DB, error) {
	m, err := f.ToORM()
	if err != nil {
		return nil, err
	}
	q :=  db.DB().Model(&orm.User{}).Where(m)
	return q, nil
}

//Из фильтра - в модель ОРМ для БД
func (f *Filter) ToORM() (*orm.User, error) {
	var (
		date time.Time
		err error
	)
	if f.Date != "" {
		if date, err = time.Parse(DATE_LAYOUT, f.Date); err != nil {
			return nil, err
		}
	}
	return &orm.User{
		Id: f.UUId,
		Login: f.Login,
		CreatedAt: date,
	}, nil
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
	q, err := f.GetQuery()
	if err != nil {
		return err
	}
	if err := q.Find(&users).Error; err != nil {
		return err
	}
	for _, u := range users {
		*usrs = append(*usrs, u)
	}
	return nil
}

//Метод удаления пользователей по фильтру
func (u *User) Delete(f *Filter, reply *string) error {
	q, err := f.GetQuery()
	if err != nil {
		return err
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