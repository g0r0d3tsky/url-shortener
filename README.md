# url-shortener
[![Lint Status](https://img.shields.io/github/actions/workflow/status/MisterZurg/TBank-backend-academy-URL-Shortener/golangci-lint.yml?branch=main&style=for-the-badge)](https://github.com/MisterZurg/TBank-backend-academy-URL-Shortener/actions?workflow=golangci-lint)
[![Coverage Status](https://img.shields.io/codecov/c/gh/github.com/MisterZurg/TBank_URL_shortener.svg?logo=codecov&style=for-the-badge)](https://codecov.io/gh/MisterZurg/TBank_URL_shortener)
[![](http://img.shields.io/badge/godoc-reference-5272B4.svg?style=for-the-badge)](https://pkg.go.dev/github.com/MisterZurg/TBank-backend-academy-URL-Shortener)

В результате проектных работ ожидается:
1. Подробное архитектурное описание с тщательным анализом каждого принятого решения. Здесь будут освещены такие аспекты, как причины выделения функциональности в отдельный микросервис, выбор способа коммуникации – Kafka/GRPC, логика за выбором определённого типа базы данных, и другие ключевые моменты.
2. Полноценная реализация сервиса, отвечающая всем поставленным требованиям и стандартам качества. Не забудьте написать тесты для вашего проекта.
3. Docker compose файл, содержащий все необходимые настройки для быстрого и безболезненного запуска сервиса в любой среде.
4. Документация интерфейса сервиса, включающая в себя спецификацию REST запроса для генерации короткой ссылки и прочие важные детали взаимодействия с сервисом.

<p align="center">
    <img src="src/gopher.png" alt="gopher" />
</p>

## Архитектурное описание
### Функциональные требования
- Для звданного URL, сервис генерирует уникальный сокращенный URL.
- При отправки GET запроса на сокращенный URL, происходит Redirect на оригинальный.
### Нефункциональные требования
- Сервис должен быть расширяемым и эффективным
- Надёжность, доступность
### Приблизительные оценки
- Операции записи: 1 миллион URL-адресов в день
- Операции записи в секунду: $1 милллион / 24 / 3600 = 12$
- Операции чтения: 10 к 1, т. е. 120 операций чтения в секунду
- Пусть длина среднего URL составляет 100 символов, и среднее время жизни одной ссылки - 1 год, тогда $148 * 10^6 * 365 / 1024 / 1024 / 1024 = 50 ГБ$

### Архитектура сервиса
![Architecture](src/arch.jpg)

### Database Schema
```mermaid
erDiagram
    url_data {
        uuid id
        string short_url
        string long_url
        datetime visited_at
    }

    keys{
        uuid id
        serial key
        encode string
        url_id uuid
    }
```
## Реализация
1. API Documentation first: OpenAPI 3.0
![API](src/api.png)
2. Prometheus metrics
![Prometheus](src/prometheus.png)
3. Grafana Dashboard
![Grafana](src/grafana.png)
