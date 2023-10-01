## TODO-APP
### Описание
REST API приложение для ведения **списков задач** (todo-списков).

### Функционал
Данное серверное (backend) приложение предоставляет API для **регистрации** и **аутентификации** пользователя, а также **работы** с todo-списками и отдельными задачами внутри них.

## Структура REST API
**Публичные endpoint'ы**, доступные всем:
```
POST /users - регистрация пользователя
POST /tokens - аутентификация пользователя и выдача JWT
```

**Приватные endpoint'ы**, доступные только аутентифицированным пользователям:
```
GET /profile - просмотр профиля пользователя

GET /lists - просмотр всех списков

POST /lists - создание списка
GET /lists/{id} - просмотр списка
PUT /lists/{id} - редактирование списка
DELETE /lists/{id} - удаление списка

POST /lists/{id}/tasks - добавление задачи в список
GET /lists/{id}/tasks - просмотр всех задач в списке

GET /tasks/{id} - просмотр задачи в списке
PUT /tasks/{id} - редактирование задачи в списке  
DELETE /tasks/{id} - удаление задачи из списка
```

## Схема базы данных

<p align="center">
    <img src="/assets/images/er_schema.png" width="800">
</p>

## Архитектура

<p align="center">
    <img src="/assets/images/architecture.png" width="800">
</p>

## Структура проекта
```
├── cmd
|   ├── app
|
├── configs
├── internal
|   ├── app
|   ├── controller
|   ├── entity
|   ├── store
|   ├── usecase
|
├── migrations
├── .env
├── .gitignore
├── Makefile
├── README.md
├── go.mod
├── go.sum
```

## Запуск и отладка
Все команды, используемые в процессе разработки и тестирования, фиксировались в `Makefile`.

Чтобы поднять проект, необходимо выполнить **две задачи** из `Makefile`:

```bash
make compose-build
make compose-up
```

## Техническая статья

### [Мой опыт создания REST API сервера для ведения todo-списков](/todo_paper.md)

<p align="center">
  <img src="/assets/images/todo_paper_cover.png" width="800">
</p>

`todo_paper.md` - файл со статьей

## Примеры запросов
* [Регистрация](#регистрация)
* [Аутентификация](#аутентификация)
* [Просмотр профиля](#просмотр-профиля)
* [Создание списка](#создание-списка)

### Регистрация
Регистрация пользователя:

```bash
curl --location --request POST http://localhost:8080/users \
--data-raw '{
    "email":"user@example.org",
    "password":"password"
}'
```

Пример ответа:

```bash
{
    "user_id":1,
    "email":"user@example.org"
}
```

### Аутентификация
Аутентификация пользователя и выдача JWT:

```bash
curl --location --request POST http://localhost:8080/tokens \
--data-raw '{
    "email":"user@example.org",
    "password":"password"
}' -v
```

Пример ответа:

```bash
HTTP/1.1 200 OK
Set-Cookie: token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJleHAiOjE2OTM3NjY5OTB9.FhxzjhKtylOZQYrpG88r_lH7-kssye9IWh7UsZ8_t6k
X-Request-Id: f7e4e8ab-ac7f-4266-b683-d35564927fc7
Date: Sun, 03 Sep 2023 18:44:50 GMT
Content-Length: 0
```

### Просмотр профиля
Просмотр профиля аутентифицированного пользователя:

```bash
curl --location --request GET http://localhost:8080/profile \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJleHAiOjE2OTM3NjY5OTB9.FhxzjhKtylOZQYrpG88r_lH7-kssye9IWh7UsZ8_t6k'
```

Пример ответа:

```bash
{
    "user_id":1,
    "email":"user@example.org"
}
```

### Создание списка
Создание списка задач:

```bash
curl --location --request POST http://localhost:8080/lists/ \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJleHAiOjE2OTM3NjY5OTB9.FhxzjhKtylOZQYrpG88r_lH7-kssye9IWh7UsZ8_t6k' \
--data-raw '{
  "list_title":"test title 1"
}'
```

Пример ответа:

```bash
{
    "list_id": 1,
    "list_title": "TEST TITLE 1",
    "user_id": 1
}
```