# EPAgent

A simple rest API service that executes scripts on running systems. 

## Caution 

This is by far not a recommended and secure way to run and execute scripts on a remote system. This is just an experiment.


## Using it

Make sure the binary is ran with sufficient permission to execute the scripts on the local system it is running on.

Then make you generated self-signed certificate. Included a very simple script in the repository that generates one with openssl.

Set a secret in config.json and paths to the certificate/key.

Theres two endpoints. http://localhost:8080/pwsh for powershell and /sh for linux. 

## API Reference

#### Header

| header | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `x-api-key` | `string` | **Required**. Your API key |

#### Linux 

```http
  POST /sh
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `path`      | `string` | **Required**. path to script or binary |

#### Windows 

```http
  POST /pwsh
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `path`      | `string` | **Required**. path to script |


## Example

### Linux shell
```sh
curl --header "Content-Type: application/json" --header "X-API-KEY: secret123" \
--request POST \
--data '{"path":"scripts\\test.ps1"}' \
https://172.20.10.10:8080/pwsh --insecure
```
### Powershell
```powershell
Invoke-RestMethod -SkipCertificateCheck -Method post -Uri "https://172.20.35.137:8080/pwsh" -Headers @{"X-API-KEY" = "secret123"} -Body (@{"path" = "scripts/test.ps1"} |ConvertTo-Json)
```