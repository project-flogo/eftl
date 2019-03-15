# EFTL Producer

## Installation
* Docker [docker](https://www.docker.com)
* Install [Go](https://golang.org/)
* Install the flogo [cli](https://github.com/project-flogo/cli)

## Setup
```bash
git clone https://github.com/project-flogo/eftl
cd eftl/examples/api/producer
```

## Testing
Start the EFTL server (this will take ~1 minute):
```bash
go run main.go -app
```

Create the gateway:
```bash
flogo create -f flogo.json
cd MyProxy
flogo build
```

Start the microgateway application:
```bash
bin/MyProxy
```

In another terminal start the EFTL client:
```bash
go run main.go -client
```

In another terminal make a request to the microgateway application:
```bash
curl -d "{\"message\": \"hello world\"}" -H "Content-Type: application/json" http://localhost:9096
```

The below should be visible in the EFTL client terminal:
```
{"message":"hello world"}
```

The EFTL messages has been forwarded from the Rest trigger to the EFTL server, and the message was then received by the EFTL client.
