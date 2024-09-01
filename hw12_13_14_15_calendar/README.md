#### Результатом выполнения следующих домашних заданий является сервис «Календарь»:
- [Домашнее задание №12 «Заготовка сервиса Календарь»](./docs/12_README.md)
- [Домашнее задание №13 «Внешние API от Календаря»](./docs/13_README.md)
- [Домашнее задание №14 «Кроликизация Календаря»](./docs/14_README.md)
- [Домашнее задание №15 «Докеризация и интеграционное тестирование Календаря»](./docs/15_README.md)

#### Ветки при выполнении
- `hw12_calendar` (от `master`) -> Merge Request в `master`
- `hw13_calendar` (от `hw12_calendar`) -> Merge Request в `hw12_calendar` (если уже вмержена, то в `master`)
- `hw14_calendar` (от `hw13_calendar`) -> Merge Request в `hw13_calendar` (если уже вмержена, то в `master`)
- `hw15_calendar` (от `hw14_calendar`) -> Merge Request в `hw14_calendar` (если уже вмержена, то в `master`)

**Домашнее задание не принимается, если не принято ДЗ, предшедствующее ему.**


## Migrations:

1) Создать базу данных:
**psql --host localhost --username postgres --password -c "create database calendar;"**
2) Запустить  make migrate-up psqlInfo="CONNECTION_STRING"
Example:
make migrate-up psqlInfo="postgresql://postgres:postgres@127.0.0.1:5432/calendar?sslmode=disable" 
3) Отскатить изменения можно командой  make migrate-down psqlInfo="CONNECTION_STRING"


### HTTP Server localhost:7777:

**Postman в папке test_help**

### gRPC Server localhost:8888:

**Create event request:**

{
    "userID": 1,
    "event": {
        "userID": 1,
        "title": "TitleGR44PC",
        "descr": "DescRRRr",
        "startDate": {
            "seconds": "1724612989"
        },
        "endDate": {
            "seconds": "1724612999"
        }
    }
}


**Update event request:**

{
"userID": 1,
"event": {
    "ID": 1,
    "userID": 1,
    "title": "New updated title",
    "descr": "DescRRRr",
    "startDate": {
        "seconds": "1724612985"
    },
    "endDate": {
        "seconds": "1724612985"
    }
   
}
}

**Delete event request:**

{
    "event": {
        "ID": 4
    },
    "userID": 1
}

**Get daily/weekly/monthly events request:**

{
    "date": {
        "seconds": "1724612900"
    },
    "userID": 1
}


