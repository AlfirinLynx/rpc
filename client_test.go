package main_test

import (
	"net"
	"net/rpc/jsonrpc"
	"testing"
	"fmt"
	"time"
	"github.com/antipin1987@gmail.com/rpcj/server"
	"github.com/antipin1987@gmail.com/rpcj/db"
	"github.com/antipin1987@gmail.com/rpcj/models/orm"
	"github.com/stretchr/testify/assert"
	"net/rpc"
	"log"
	"github.com/antipin1987@gmail.com/rpcj/models"
	"net/http"
	"bytes"
	"encoding/json"
)
var (
	c *rpc.Client
    rpcServer *rpc.Server
)

func init() {
	//Старт RPC сервера
	rpcServer = rpc.NewServer()
	go server.StartServerTCP(rpcServer, server.AddrTCP)
	time.Sleep(3 * time.Second)

	client, err := net.Dial("tcp", server.AddrTCP)
	if err != nil {
		log.Fatal(err)
	}
	//Клиент для обработки реквестов в формате json
	c = jsonrpc.NewClient(client)
}

func countRecords(f *models.Filter) int {
	var count int
	q, err := f.GetQuery()
	if err != nil {
		panic(err)
	}
	q.Count(&count)
	return count
}

func dateByTime(t time.Time) string {
	y, m, d := t.Date()
	return fmt.Sprintf("%04d-%02d-%02d", y, m, d)
}


//Тест добавления пользователя с данным логином
func TestAddUser(t *testing.T) {
	var (
		rep bool //ответ - успех/неуспех (true/false) выполнения метода
		argLogin = "kvAnt"
	)

	f := &models.Filter{Login: argLogin} //Фильтр (логин - "kvAnt")
	//Перед вызовом метода считаем записи в БД
	countBefore := countRecords(f)
	fmt.Println("Records in db before call to User.Add: ", countBefore)

	//Запрос на json-RPC сервер по TCP
	err := c.Call("User.Add", argLogin, &rep)
	if err != nil {
		t.Error("User.Add error:", err)
		return
	}
	fmt.Println("Success: ", rep)
	assert.True(t, rep)

	//Считаем записи в БД после вызова метода
	countAfter := countRecords(f)
	fmt.Println("Records in db after call to User.Add: ", countAfter)

	//Проверяем правильность выполнения запроса
	assert.True(t, countAfter - countBefore == 1)
}

//Тест поиска пользователей по фильтру
func TestUserFind(t *testing.T) {
	var (
		login = "Cat5"
		date = dateByTime(time.Now()) // Извлекаем сегодняшнюю дату
		f = &models.Filter{Login: login, Date: date}  // Фильтр для поиска (логин "Cat5", создан сегодня)
		usrs = make([]orm.User, 0)  // для результата поиска
	)
	//Создадим запись в БД
	u, err := f.ToORM()
	if err != nil {
		t.Error(err)
		return
	}
	if err := db.DB().Create(u).Error; err != nil {
		t.Error(err)
		return
	}

	//Считаем записи
	count := countRecords(f)
	fmt.Println("Records in db before call to User.Find: ", count)

	//Теперь ищем запись (логин "Cat", создан сегодня), отправляя запрос по чистому TCP на json-RPC сервер
	err = c.Call("User.Find", f, &usrs)
	if err != nil {
		t.Error("User.Find error: ", err)
		return
	}

	//Выводим найденные записи
	fmt.Println("Success: ", usrs)
	fmt.Printf("Returned %d records in response\n", len(usrs))
	//Проверяем правильность выполнения запроса
	assert.True(t, count == len(usrs))
}

//Тест удаления пользователей по фильтру
func  TestUserDelete(t *testing.T) {
	var (
		rep string // Для строки с ожидаемым результатом - "... records deleted"
		login = "trash"
		f = &models.Filter{Login: login} //фильтр
	)

	//Создадим 3 записи с одинаковым логином
	for i := 0; i < 3; i++ {
		if err := db.DB().Create(&orm.User{Login: login}).Error; err != nil {
			t.Error(err)
			return
		}
	}

	//Вызываем метод удаления записей с данным логином
	err := c.Call("User.Delete", f, &rep)
	if err != nil {
		t.Error("User.Delete error: ", err)
		return
	}
	//Сообщение-ответ сервера
	fmt.Println(rep)

	//Удостоверимся, что все записи удалены
	count := countRecords(f)
	assert.True(t, count == 0)
}



//Тест http-сервера
//Отправляем запросы, в теле которых json вида: {"method":"User.Find","params":[{"Date":"2017-07-08", "login":"Cat7"}],"id":0},
// {"method":"User.Add","params":["Dog172967"],"id":0}, {"method":"User.Delete","params":[{"uuid":"","login":"Dog172967"}],"id":0} и т. д.

func TestHttp(t *testing.T)  {
	//Слушаем http-запросы
	http.HandleFunc("/", (&server.Serv{Server: rpcServer}).HttpHandler)
	go func() {
		log.Fatal(http.ListenAndServe(server.AddrHTTP, nil))
	}()
	time.Sleep(2 * time.Second)

	//Делаем http-запрос, например, на метод "User.Add"
	url := fmt.Sprintf("http://%s", server.AddrHTTP)
	jsonBody := bytes.NewBufferString(`{"method":"User.Add","params":["Dog"],"id":15}`)
	resp, err := http.Post(url, "application/json", jsonBody)
	if err != nil {
		t.Error(err)
		return
	}
	defer resp.Body.Close()

	//Посмотрим ответ
	v := struct{
		Id int `json:"id"`
		Result interface{} `json:"result"`
		Error error `json:"error"`
	}{}
	dec := json.NewDecoder(resp.Body)
	 if err :=dec.Decode(&v); err != nil {
		 t.Error(err)
		 return
	 }
	fmt.Println("Response: ", v)
	assert.True(t, v.Id == 15)
	r, ok := v.Result.(bool)
	if !ok {
		t.Error("Wrong response result type")
	}
	assert.True(t, r)
}

