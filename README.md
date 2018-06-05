# stressed

Stress testing tool for APIs via [vegeta](https://github.com/tsenart/vegeta)

## Config

### JSON

```json
{
  "baseUrl": "http://localhost:8080",
  "tests": [
    {
      "name": "basic",
      "rps": 100,
      "duration": 1,
      "target": {
        "method": "POST",
        "path": "/test",
        "headers": [
          {
            "name": "X-Request-Id",
            "value": "12345"
          }
        ],
        "body": "ewoic2F5IiA6ICJoZWxsbyIKfQ=="
      },
      "sla": {
        "latency": 150,
        "successRate": 99.9
      }
    }, {
      "name": "other",
      "rps": 100,
      "duration": 1,
      "target": {
        "method": "POST",
        "path": "/test",
        "headers": [
          {
            "name": "X-Request-Id",
            "value": "54321"
          }
        ],
        "body": "ewoic2F5IiA6ICJoZWxsbyIKfQ=="
      },
      "sla": {
        "latency": 150,
        "successRate": 99.9
      }
    }
  ]
}
```

## YAML

```yaml
---
baseUrl: http://localhost:8080
tests:
- name: 'basic' 
  rps: 100
  duration: 1
  target:
    method: POST
    path: '/test'
    headers:
    - name: X-Request-Id
      value: '12345'
    body: ewoic2F5IiA6ICJoZWxsbyIKfQ==
  sla:
    latency: 150
    successRate: 99.9
- name: 'other' 
  rps: 100
  duration: 1
  target:
    method: POST
    path: '/test'
    headers:
    - name: X-Request-Id
      value: '54321'
    body: ewoic2F5IiA6ICJoZWxsbyIKfQ==
  sla:
    latency: 150
    successRate: 99.9
```
