## Мой опыт создания TODO REST API сервера

### Подготовка репозитория
На Windows
```bash
E:
cd E:\HDD\Magistracy\EducationCourses\ITMO\
mkdir todo-app
cd .\todo-app
```

Создал `README.md` и `LICENSE`.

В README записал структуру проекта.

Также добавил ER-диаграмму базы данных, которую создал в [dbdesigner](https://app.dbdesigner.net) и сохранил в директорию `/assets/images`.

Создал `.gitignore` для игнорирования исполняемого файла приложения `app.exe` и файла с переменными окружения `.env` (`/` обозначает, что будут исключены только файлы находящиеся в корневой директории проекта).
```
/app.exe
/.env
```

Инициализировал репозиторий, сделал `initial commit` и опубликовал его через IDE.

```bash
git init
git add -A
git commit -m "initial commit"
```


