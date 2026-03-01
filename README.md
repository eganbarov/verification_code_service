# verification_code_service
Service for generation and verification code which is sent.

## Примеры использования
### Отправка кода
#### Запрос
```
curl -X POST -H "Content-Type: application/json" -d '{"phone": "89067001910","action":"auth"}' http://localhost:8080/send-code
```
#### Ответ
```json
{"isSent":true,"statusCode":201,"error":""}
```

### Верификация кода
#### Запрос
```
curl -X POST -H "Content-Type: application/json" -d '{"phone": "89067001910","code":"646482","action":"auth"}' http:/localhost:8080/validate-code
```
#### Ответ
```json
{"isValid":true,"statusCode":200,"error":""}
```

