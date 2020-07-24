# Authorization
## User registration and authorization service  

Requirements:  
`Go 1.13.5`  
`Postgres 10.12`   

Сервис для регистрации и авторизации пользователя.  
API:  
POST /api/registration  
Request:  
`{  "email": "test@example.com",  "login": "username123",  "phone_number": "89000000000",  "password": "qwerty123"}`  
Response:  
`{ "Success": true}`
POST /api/login  
Request:  
` {"login": "username123", "password": "qwerty123"}`  
Response:
` { "access_token": ... , "refresh_token": ...}`


### Запуск:  
Создайте базу и пользователя.  
`sudo -u postgres psql`   

`CREATE DATABASE auth_db;`  
`CREATE USER auth_user WITH password 'auth_password';`  
`GRANT ALL ON DATABASE auth_db TO auth_user;`

Заполните конфиг нужными данными. Пример конфига представлен в файле config.json.  
"secrets" - Криптографические ключи для создания токенов  
"http" - Порт и таймауты хттп-сервера    
"db" - Соединение с базой  

Укажите в переменной среды `AUTH_CONFIG` путь к конфигу и запустите сервис:  
`export AUTH_CONFIG="$PWD/config.json" && go run internal/cmd/main.go`


### Тесты:  
Для удобства конфиг для тестов указывается в другой переменной `TEST_AUTH_CONFIG`  
`export TEST_AUTH_CONFIG="$PWD/test_config.json" && go test -v ./internal/cmd/`