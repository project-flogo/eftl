# EFTL Activity
This activity provides your microgateway application the ability to send EFTL messages.

The available service `settings` are as follows:

| Name       |  Type   | Description                                   |
|:-----------|:--------|:----------------------------------------------|
| url        | string  | The EFTL server URL                           |
| id         | string  | The id for this EFTL client                   |
| user       | string  | The user name for the EFTL server             |
| password   | string  | The password for the EFTL server              |
| ca         | string  | The certificate authority for the EFTL client |

The available `input` for the request are as follows:

| Name        |  Type       | Description                           |
|:------------|:------------|:--------------------------------------|
| content     | JSON object | The message to send                   |
| dest        | string      | The EFTL dest to send the messages to |

The available response `outputs` are as follows:

| Name   |   Type   | Description   |
|:-------|:---------|:--------------|


A sample `service` definition is:

```json
{
  "name": "EFTLGateway",
  "description": "EFTL gateway",
  "ref": "github.com/project-flogo/eftl/activity",
  "settings": {
    "url": "ws://localhost:9191/channel"
  }
}
```

An example `step` that invokes the above `SQLSecurity` service using `payload` is:

```json
{
  "service": "EFTLGateway",
  "input": {
    "content": "=$.payload.content",
    "dest": "sample"
  }
}
```
