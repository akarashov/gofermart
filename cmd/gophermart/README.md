# cmd/gophermart

В данной директории будет содержаться код накопительной системы лояльности, который скомпилируется в бинарное приложение.

// адрес и порт запуска сервиса: переменная окружения ОС `RUN_ADDRESS` или флаг `-a`
// адрес подключения к базе данных: переменная окружения ОС `DATABASE_URI` или флаг `-d`
// адрес системы расчёта начислений: переменная окружения ОС `ACCRUAL_SYSTEM_ADDRESS` или флаг `-r`

#### Accrual start
./projects/gofermart/cmd/accrual/accrual_darwin_amd64 \
-a ":8081" \
-d "postgres://admin:admin@192.168.0.20:5432/demo?search_path=gofermart&sslmode=disable"

#### Accrual get
curl -v -X GET http://127.0.0.1:8081/api/orders/4111111111111111



#### **Регистрация пользователя**
curl -v -X POST http://127.0.0.1:8080/api/user/register \
-H 'Content-Type: application/json' \
-d '{"login":"user002","password":"P@ssword002"}'


#### **Аутентификация пользователя**
curl -s -D - -o /dev/null \
-X POST http://127.0.0.1:8080/api/user/login \
-H 'Content-Type: application/json' \
-d '{"login":"user002","password":"P@ssword002"}' \
| awk '/thorization/ {print $3}' | tr -d '\r' > token.txt


#### **Загрузка номера заказа**
curl -v -X POST http://127.0.0.1:8080/api/user/orders \
-H "Authorization: Bearer $(cat token.txt)" \
-H "Content-Type: text/plain" \
-d '717413'

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
curl -v -X GET http://127.0.0.1:8080/api/user/orders \
-H "Authorization: Bearer $(cat token.txt)"


#### **Получение текущего баланса пользователя**
curl -v -X GET  http://127.0.0.1:8080/api/user/balance \
-H "Authorization: Bearer $(cat token.txt)"


#### **Запрос на списание средств**
curl -v -X POST http://127.0.0.1:8080/api/user/balance/withdraw \
-H "Authorization: Bearer $(cat token.txt)" \
-H 'Content-Type: application/json' \
-d '{"order":"5555555555554444","sum":333}'

# TODO
Поправить в запросе число, а json ждет строку
json: cannot unmarshal number into Go struct field WithdrawalRequest.sum of type string

2720991234567890
4916730123456782
4000056655665556
6271136264806190

#### **Получение информации о выводе средств**
curl -v -X GET http://127.0.0.1:8080/api/user/withdrawals \
-H "Authorization: Bearer $(cat token.txt)"


