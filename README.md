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

## Техническая статья

### [Мой опыт создания REST API сервера для ведения todo-списков](/todo_paper.md)

<p align="center">
  <img src="/assets/images/todo_paper_cover.png" width="800">
</p>

`todo_paper.md` - файл со статьей
