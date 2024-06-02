## To do list

### Запуск
```
docker-compose up --build
```
### Документация
```
localhost:8080/swagger
```
### Описание запросов
Получить все задачи (без пагинации)
```
Status OK

curl -X GET 'http://localhost:8080/tasks'
``` 
Получить все задачи с паганицией
```
Status OK 

curl -X GET 'http://localhost:8080/tasks?pageSize=2&page=1'

curl -X GET 'http://localhost:8080/tasks?pageSize=1000&page=1000'

curl -X GET 'http://localhost:8080/tasks?pageSize=10'
```
```
Bad Request 

curl -X GET 'http://localhost:8080/tasks?pageSizeeeeee=2&page=1'
```
Получить все задачи с фильтром по дате
```
Status OK

curl -X GET 'http://localhost:8080/tasks?&date=2024-06-02'

curl -X GET 'http://localhost:8080/tasks?&date=10024-06-02'

curl -X GET 'http://localhost:8080/tasks?&data=10024-06-02' (игнорируются неправильно написанные поля)
```
```
Bad Request

curl -X GET 'http://localhost:8080/tasks?&date=2024-06-02'
```
Получить все задачи с фильтром по статусу выполнено/не выполнено
```
Status OK

curl -X GET 'http://localhost:8080/tasks?done=true'
```

Создать задачу
```
Status Created

curl -X POST 'http://localhost:8080/tasks'   -H "Content-Type: application/json"   -d '{
    "header":"docker wait",
    "description": "make docker wait for postgres",
    "task_date": "2024-06-03"
}' 
```
Получить задачу по айди
```
Status OK

curl -X GET 'http://localhost:8080/tasks/1
```
Обновить задачу (1+ поле)
```
Status OK

curl -X PATCH "http://localhost:8080/tasks/1"   -H "Content-Type: application/json"   -d '{
    "done": true
}' 

```
Удалить задачу
```
Status OK
curl -X DELETE "http://localhost:8080/tasks/2" 
```