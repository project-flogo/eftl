# SQL Detector

The `sqld` service type implements SQL injection attack detection. Regular expressions and a [GRU](https://en.wikipedia.org/wiki/Gated_recurrent_unit) recurrent neural network are used to detect SQL injection attacks.

The available trigger `settings` are as follows:

| Name     | Type   | Description                           |
|:---------|:-------|:--------------------------------------|
| url      | string | The URL of the EFTL server            |
| id       | string | The client ID of the EFTL trigger     |
| user     | string | The login user for the EFTL server    |
| password | string | The login passwod for the EFTL server |
| ca       | string | The certificate for the EFTL server   |

The available trigger `handler settings` are as follows:

| Name | Type   | Description                  |
|:-----|:-------|:-----------------------------|
| dest | string | The destination to listen on |

The available `output` for the request are as follows:

| Name       |  Type       | Description                     |
|:-----------|:------------|:--------------------------------|
| content    | JSON object | The content of the EFTL message |

## Example

A sample `trigger` definition is:

```json
{
  "name": "flogo-eftl",
  "id": "MyProxy",
  "ref": "github.com/project-flogo/eftl/trigger",
  "settings": {
    "url": "ws://localhost:9191/channel"
  },
  "handlers": [
    {
      "settings": {
        "dest": "sample"
      },
      "actions": [
        {
          "id": "microgateway:Pets"
        }
      ]
    }
  ]
}
```
