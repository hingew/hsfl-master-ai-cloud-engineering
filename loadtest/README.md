# Run loadtest:
```go run main.go```

# Config-File:

```
{
    "users": 5,
    "rampup": 5,
    "duration": 5,
    "cooldown": 5,
    "target": "localhost:8000",
    "path": "/admin/",
    "targets": [
        "localhost:8000"
    ]
}
```

# Diagram:
![diagram](loadtest-example-config.png)



