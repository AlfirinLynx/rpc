# json rpc сервер с использованием пакета net/rpc/json
Сервер слушает на двух портах: один - на чистом tcp, другой - по http.
Для работы с БД была использована gorm из проекта github.com/jinzhu/gorm.
Вся конфигурационная информация находится в файле config.yaml. Тесты находятся в файле client_test.go
