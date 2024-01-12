# Config-File:

To configure the loadtest adapt the config.json file to your needs.

The following example for the config file produces the following graph:
![diagram](example-config-curve.png)

```
{
    "rampSpecifications": [
        {"duration": 20, "targetRPS": 20},
        {"duration": 10, "targetRPS": 50},
        {"duration": 20, "targetRPS": 50},
        {"duration": 30, "targetRPS": 100},
        {"duration": 30, "targetRPS": 60},
        {"duration": 20, "targetRPS": 60},
        {"duration": 10, "targetRPS": 80},
        {"duration": 20, "targetRPS": 40},
        {"duration": 40, "targetRPS": 0}
    ],
    "target": "192.168.2.134:32674",
    "paths": [
        "/"
    ]
}
```

Replace the target with the target address of your api-gateway service.

# Run loadtest:

To run the load test just run `go run main.go`
