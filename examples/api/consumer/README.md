# EFTL Consumer
This recipe demonstrates the use of the EFTL trigger to generate HTTP requests from EFTL messages.

## Installation
* Docker [docker](https://www.docker.com)
* Install [Go](https://golang.org/)

## Setup
```bash
git clone https://github.com/project-flogo/eftl
cd eftl/examples/api/consumer
```

## Testing
Start the EFTL server and microgateway application (this will take ~1 minute):
```bash
go run main.go -app
```

In another terminal start the target HTTP server:
```bash
go run main.go -target
```

In another terminal execute the EFTL client:
```bash
go run main.go -client
```

The target terminal should print out a message as below:
```
2019/03/15 14:55:14 /a
2019/03/15 14:55:14 application/json; charset=UTF-8
2019/03/15 14:55:14 {"message":"hello world"}
```

This demonstrates the EFTL payload was forwarded to the HTTP target service.
