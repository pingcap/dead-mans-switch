# dead-mans-switch
dead mans switch is a simple prometheus alert manager webhook service. it provides a basic mechanisms to ensure alerting pipeline is healthy.


prometheus provider a mechanisms to always firing. this is call `WatchDog` in prometheus-operator. If `WatchDog` does not firing, we can think alerting pipeline is unhealthy.
dead mans switch can used for receive `WatchDog` webhook payload. evaluate alerting payload as expected optional. (sometimes your have a aggregated alert manager services, will have multiply prometheus
instance send `WatchDog` alert to alert manager). we can use evaluator to ensure all prometheus is healthy and external label as expected.

```yaml
- alert: Watchdog
  annotations:
    message: |
      This is an alert meant to ensure that the entire alerting pipeline is functional.
      This alert is always firing, therefore it should always be firing in Alertmanager
      and always fire against a receiver. There are integrations with various notification
      mechanisms that send a notification when this alert is not firing. For example the
      "DeadMansSnitch" integration in PagerDuty.
  expr: vector(1)
  labels:
    severity: none
```
      
## Develop
build binary in local environment
```sh
make build
```

running in local environment
```sh
dms -config ./manifest/config.example.yaml
```

send alert manager webhook payload
```sh
curl -H "Content-Type: application/json" --data @payload.json http://localhost:8080/webhook
```

## Deploy

The `manifest/deploy` directory have k8s deploy yaml files, you can copy it and update <pagerduty> in configmap.
The `manifest/monitoring` directory have `ServiceMonitor` and `PrometheusRule` crd file, if you use prometheus-operator monitor your k8s clusters, you can trying for it.

### AlertManager config
let WatchDog alert send to dead mans switch receivers
```yaml
route:
  routes:
    - receiver: dead-mans-switch
      group_wait: 10s
      group_interval: 30s
      repeat_interval: 15s
      match:
        alertname: 'Watchdog'
```

add dead mans switch service as a new webhook receiver
```yaml
receivers:
- name: dead-mans-switch
  webhook_configs:
  - url: http://dead-mans-switch:8080/webhook
```