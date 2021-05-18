# task-tracker
API для отслеживания выполнения задач  

## Оглавление
- [Конфигурация](#configuration)
- [Развертывание](#deployment)

<a name="configuration"></a>
## Конфигурация
Конфигурация происходит следующим образом:  
1. Читается конфиг по пути, указанному в переменной окружения `CONFIG_PATH`.
Файл должен иметь формат YAML и иметь определенную структуру.
2. Читаются оставшиеся настройки из переменных окружения.
При дублировании настроек переменные окружения затирают параметры конфига.

Список переменных окружения:  
```
CONFIG_PATH=configs/config.yaml
LOGGER_LEVEL=debug
LOGGER_FORMAT=default
PG_ADDRESS=0.0.0.0:5432
PG_USER=task-tracker
PG_PASSWORD=123
PG_DATABASE=task-tracker
REDIS_ADDRESS=0.0.0.0:6379
JWT_SIGNING_KEY=some_key
COOKIE_HASH_KEY=some_key
COOKIE_BLOCK_KEY=some_key
COOKIE_DOMAIN=task-tracker.com
EMAIL_SERVER_ADDRESS=smtp.gmail.com:587
EMAIL_USERNAME=user@test.com
EMAIL_PASSWORD=some_password
```

<a name="deployment"></a>
## Развертывание
1. Для того, чтобы развернуть сервис в docker:  
```docker-compose up -d --build```  

    Опустить контейнеры:  
```docker-compose down```  
2. Чтобы выполнить начальную миграцию для базы данных, нужно установить <a href="https://github.com/golang-migrate/migrate">эту утилиту</a> и выполнить команду:  
```migrate -path ./schema -database 'postgres://task-tracker:123@localhost:54320/task-tracker?sslmode=disable' up```  

    Откатить миграцию:  
```migrate -path ./schema -database 'postgres://task-tracker:123@localhost:54320/task-tracker?sslmode=disable' down```  
