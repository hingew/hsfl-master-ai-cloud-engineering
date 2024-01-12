# Run loadtest:

`go run main.go`

# Config-File:

```
{
    "rampSpecifications": [
        {"duration": 5, "RPSincrement": 1},
        {"duration": 5, "RPSincrement": 2},
        {"duration": 5, "RPSincrement": 3}
    ],
    "target": "192.168.178.98:31153",
    "paths": [
        "/admin",
        "/api"
    ]
}
```

# Diagram:

![diagram](loadtest-example-config.png)