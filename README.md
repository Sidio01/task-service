## Service for tasks of approval and sending letters

* implements a REST API for creating a task for approval and sending unique solution links to participants.
* The list of participants is indicated explicitly (list of email) for each task.
* Authentication of calls to the REST API is validated on the authentication service through GRPC calls.
* The REST API must implement CRUDL for negotiation tasks. Operations U and D are allowed only to the author of the problem.
* When creating/updating a task, each participant is sent a letter with two unique links - "agreed" and "not agreed". First, a letter is sent to the first approver, his reaction is expected, then the next one, his reaction is expected, etc. until the last match. The API must have methods for handling "clicks" on sent links and registering the appropriate response from the approver.
* If at least one "not approved" link was clicked, the task is considered not approved as a whole, and all participants are sent a notification about the end of approval with a negative result, and letters with links to the next approvers for this task are no longer sent.
* The service generates and sends events for creating tasks, sending letters, clicking links to kafka.

---

## Run using Docker

### Project setup

Create a `.env` file at the root of the repository:

```bash
cp .env.example .env
```

Make adjustments to the environment variables as needed.

### Building images and running containers

At the root of the repository, run the command:

```bash
docker-compose up --build
```

### Stop containers

To stop containers, run the command:

```bash
docker-compose stop
```

---

Documentation - http://localhost:3000/swagger/index.html

Requests should be sent to http://localhost:3000/

---
---

## Сервис задач согласования и отправки писем 

* реализует REST API для создания задачи на согласование и рассылки уникальных ссылок-решений участникам. 
* Список участников указывается явно (список email) для каждой задачи. 
* Аутентификация обращений на REST API валидируется на сервисе аутентификации посредством GRPC-вызовов. 
* REST API должно реализовывать CRUDL для задач согласования. Операции U и D позволены только автору задачи. 
* Каждому участнику при создании/обновлении задачи высылается письмо с двумя уникальными ссылками - "согласовано" и "не согласовано". Сначала отправляется письмо первому согласующему, ожидается его реакция, затем следующему, ожидается его реакция, и т.д. до последнего согласующего. API должно иметь методы для обработки "нажатий" на высланные ссылки и регистрации соответствующей реакции согласующего. 
* Если была нажата хотя бы одна ссылка "не согласовано", задача считается в целом не согласованной, и всем участникам рассылается уведомление об окончании согласования с негативным результатом, и письма со ссылками следующим согласующим по этой задаче уже не отправляются. 
* Сервис формирует и отправляет в kafka события создания задач, отправки писем, нажатия на ссылки.

---

## Запуск с использованием Docker

### Настройка проекта

Создайте `.env` файл в корне репозитория:

```bash
cp .env.example .env
```

Внесите при необходимости корректировки в переменные окружения.

### Сборка образов и запуск контейнеров

В корне репозитория выполните команду:

```bash
docker-compose up --build
```

### Остановка контейнеров

Для остановки контейнеров выполните команду:

```bash
docker-compose stop
```

---

Документация - http://localhost:3000/swagger/index.html

Запросы следует отправлять на http://localhost:3000/
