Тестовое задание Cloud.ru Camp 2025

# Load Balancer на Go

Простой балансировщик нагрузки на языке Go с поддержкой round-robin алгоритма, 
проверкой здоровья бэкендов и graceful shutdown.

## Возможности
* Балансировка нагрузки по алгоритму round-robin

* Проверка здоровья бэкендов (health checks)

* Поддержка graceful shutdown

* Конфигурация через JSON/YAML файлы

* Подробное логирование

## Установка

1. Клонируйте репозиторий:
````
git clone git@github.com:lolita-tr/test-go-cloud.git
cd test-go-cloud/loadbalancer
````

2. Установите зависимости:
````
go mod tidy
````
## Запуск
````
go run cmd/loadbalancer/main.go --config=configs/config.yaml
````

## Параметры командной строки
| Парамерт   | Описание                  | По умолчанию |
|------------|---------------------------|--------------|
|--config| Путь к конфигурационному файлу |config.yaml|

## Структура проекта
````

loadbalancer/
├── cmd/
│   └── loadbalancer/
│       └── main.go        # Точка входа
├── internal/
│   ├── balancer/          # Логика балансировки
│   │   ├── backend.go     # Определение бэкенд-сервера
│   │   ├── roundrobin.go  # Round-robin алгоритм
│   │   ├── balancer.go    # Основная логика балансировщика
│   │   └── healthcheck.go # Проверка здоровья
│   └── config/            # Конфигурация
│       └── config.go      # Загрузка конфигурации
├── configs          
│   └──config.yaml         # Пример конфига
├── Dockerfile
├── backend.Dockerfile
├── docker-compose.yml
├── backends.go            # Запуск тестовых серверов  
├── go.mod
├── go.sum
└── README.md
````
## Логирование

````
[YYYY/MM/DD HH:MM:SS.MICROS] [LEVEL] Message
````

Уровни логирования 

* INFO: Основная информация о работе

* WARN: Потенциальные проблемы

* ERROR: Критические ошибки

## Пример использования

1. Запустить тестовые серверы
````
PORT=8081 go run backends.go
PORT=8082 go run backends.go
PORT=8083 go run backends.go
````
2. Запустить балансировщик
````
go run cmd/loadbalancer/main.go --config=configs/config.json
````

3. Отправлять запросы(в другом терминале):
````
curl http://localhost:8080
````

## Docker

1. Собрать сервисы(балансировщик и тестовые сервисы)
````
docker-compose build
````
2. Запустить собранные сервисы 
````
docker-compose up
````
3. Отправлять запросы(в другом терминале):
````
curl http://localhost:8080
````