# dead-mans-switch

## Develop
build binary in local environment
```sh
make build
```

running in local environment
```sh
dms -config config.example.yaml
```

send alert manager webhook payload
```sh
curl -H "Content-Type: application/json" --data @payload.json http://localhost:8080/webhook
```
