# cmd/gophermart

В данной директории будет содержаться код накопительной системы лояльности, который скомпилируется в бинарное приложение.

// адрес и порт запуска сервиса: переменная окружения ОС `RUN_ADDRESS` или флаг `-a`
// адрес подключения к базе данных: переменная окружения ОС `DATABASE_URI` или флаг `-d`
// адрес системы расчёта начислений: переменная окружения ОС `ACCRUAL_SYSTEM_ADDRESS` или флаг `-r`

## Accrual

### Accrual start

```bash
./projects/gofermart/cmd/accrual/accrual_darwin_amd64 \
-a ":8081" \
-d "postgres://accrual:accrual@127.0.0.1:5432/accrual?sslmode=disable"
```

### Accrual registers new mechanics in accrual service

```bash
curl -v -X POST http://127.0.0.1:8081/api/goods \
-H 'Content-Type: application/json' \
-d '{"match": "XXXXXXXXXXX","reward": 5,"reward_type": "%"}'
```

### Accrual register new order

```bash
curl -v -X POST http://127.0.0.1:8081/api/orders \
-H 'Content-Type: application/json' \
-d '{"order": "717413","goods": [{"description": "Стиральная машинка LG","price": 47399.99},{"description": "Телевизор XXXXXXXXXXX","price": 14599.50}]}'
```

### Accrual get

```bash
curl -v -X GET http://127.0.0.1:8081/api/orders/717413
```

## Gofermart

### User

#### **Регистрация пользователя**

```bash
curl -v -X POST http://127.0.0.1:8080/api/user/register \
-H 'Content-Type: application/json' \
-d '{"login":"user002","password":"P@ssword002"}'
```

#### **Аутентификация пользователя**

```bash
curl -s -D - -o /dev/null \
-X POST http://127.0.0.1:8080/api/user/login \
-H 'Content-Type: application/json' \
-d '{"login":"user002","password":"P@ssword002"}' \
| awk '/thorization/ {print $3}' | tr -d '\r' > token.txt
```

### Order

#### **Загрузка номера заказа**

```bash
curl -v -X POST http://127.0.0.1:8080/api/user/orders \
-H "Authorization: Bearer $(cat token.txt)" \
-H "Content-Type: text/plain" \
-d '717413'
```

// Luhn valid numbers
4111111111111111
4532015112830366
5454545454545454
371449635398431
6011111111111117
30569309025904
3566002020360505
6200000000000005
5063511470022499
2223000048400011

#### **Получение списка загруженных номеров заказов**

```bash
curl -v -X GET http://127.0.0.1:8080/api/user/orders \
-H "Authorization: Bearer $(cat token.txt)"
```

### Balance

#### **Получение текущего баланса пользователя**

```bash
curl -v -X GET  http://127.0.0.1:8080/api/user/balance \
-H "Authorization: Bearer $(cat token.txt)"
```

#### **Запрос на списание средств**

```bash
curl -v -X POST http://127.0.0.1:8080/api/user/balance/withdraw \
-H "Authorization: Bearer $(cat token.txt)" \
-H 'Content-Type: application/json' \
-d '{"order":"717413","sum":33.355}'
```

// Luhn valid numbers
2720991234567890
4916730123456782
4000056655665556
6271136264806190

#### **Получение информации о выводе средств**

```bash
curl -v -X GET http://127.0.0.1:8080/api/user/withdrawals \
-H "Authorization: Bearer $(cat token.txt)"
```
