# verification_code_service
Service for generation and verification code which is sent.

## Description
Now a service has only a mock implementation of code's sender (which output it).

Need to write implementation of 
> sender/sms-code-sender.go

Just use your favorite sender provider.

## Examples
### HealthCheck
#### Request
```
curl http://localhost:8080/health-check
```
#### Response
```json
{"statusCode":200,"msg":"Ok"}
```

### Send code
```
required params: "phone" and "action"
```
#### Request
```
curl -X POST -H "Content-Type: application/json" -d '{"phone": "89067001910","action":"auth"}' http://localhost:8080/send-code
```
#### Response
```json
{"isSent":true,"statusCode":201,"error":""}
```

### Verify code
```
required params: "phone", "action", "code"
```
#### Request
```
curl -X POST -H "Content-Type: application/json" -d '{"phone": "89067001910","code":"646482","action":"auth"}' http:/localhost:8080/validate-code
```
#### Response
```json
{"isValid":true,"statusCode":200,"error":""}
```

