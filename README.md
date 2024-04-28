# Entrypoint for docker containers

Entrypoint for running apps in containers with:
1. Optional generation env variables (only for child process) from Vault secrets. Windows version also set env variables in Registry system-wide
2. SIGTERM and SIGINT propagation to child process
3. Wait for child process for finish and exit with child's exit code

## Entrypoint binaries delivery

### With built-in base image

You could use next Dockerfiles as example to build your base image:
- [Dockerfile for Linux for Node.JS apps](Dockerfile)
- [Dockerfile for Windows for Dot.Net apps](Dockerfile.windows)

Applications CI will use those base images in `FROM`

### With S3 storage and Kubernetes host_mount

1. Create an S3 bucket (like `infra-binaries`)
2. Upload binaries (for linux and windows) to the S3 bucket
    1. New binary should be uploaded to the temp name like `entrypoint.tmp`
    2. Old binary should be renamed to the `entrypoint.old`
    3. New binary should be renamed from temp name `entrypoint.tmp` to `entrypoint`
3. Every k8s node contains a bootstrap code to download relevant entrypoint binary
    1. For linux nodes:
    ```
                pre_bootstrap_user_data = <<-EOT
                #!/bin/bash
                mkdir -p /entrypoint
                aws s3 cp s3://infra-binaries/entrypoint/entrypoint /entrypoint/entrypoint || aws s3 cp s3://infra-binaries/entrypoint/entrypoint.old /entrypoint/entrypoint
                chmod +x /entrypoint/entrypoint
                EOT
    ```
    2. For windows nodes:
    ```
                pre_bootstrap_user_data = <<-EOT
                Read-S3Object -BucketName "infra-binaries" -Key "entrypoint/entrypoint.exe" -Region "eu-west-2" -File "/entrypoint/entrypoint.exe"; if (-not $?) { Read-S3Object -BucketName "infra-binaries" -Key "entrypoint/entrypoint.exe.old" -Region "eu-west-2" -File "/entrypoint/entrypoint.exe" }
                EOT
    ```
4. Configure POD with host volume mount `/entrypoint/`
5. Configure POD's `command` (entrypoint) changed to `/entrypoint/entrypoint` for linix and `/entrypoint/entrypoint.exe` for windows
6. To update `entrypoint` on nodes, could use project [go-entrypoint-updater](https://github.com/alt-dima/go-entrypoint-updater)

## Entrypoint logic workflow
1. Check if `VAULT_ADDR` env var configured and Vault is reacheble and ready by endpoint `/v1/sys/health`
3. If list with required Vault secrets is not empty: 
    1. Read secrets list from `SECRETS_SOURCE_CONFIG` env var, by default: `./secrets_config.json#secrets_list` (`./secrets_config.json` - json file path, `secrets_list` - json path inside file)
    2.  Init Vault Client with credentials (env vars `VAULT_APPROLE_RID` and `VAULT_APPROLE_SID`)
    3. Read required secrets from Vault and set env varibales with these values to the child
4.  Run child app process with defined arguments
6.  Wait until process will be terminated (with signals propagation) or exited by itself

### Vault secrets:
Regular `/secret/{secret_path}` will be used.

Required secrets configuration (`secrets_config.json` example):
```json
{
    "secrets_list": [
        "mongodb",
        "rabbitmq",
        {
            "secretname": "mysql#local",
            "envvarname": "env1"
        }, 
    ]
}
```
## Usage:
```bash
export SECRETS_SOURCE_CONFIG=./secrets_config.json#secrets_list
export VAULT_APPROLE_SID=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
export VAULT_APPROLE_RID=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
export VAULT_ADDR=https://vault-api-address

entrypoint node app.js appparam1 appparam2 appparam3
```
Listed secrets from `secrets_config.json` file will be provided as a child's process env vars (and container-wide for windows) in the following format:
Non `[^a-zA-Z0-9_]` characters in the secret path will be replaced with `_` (like envconsul did)

```bash
echo $secret_mongodb_url1
secret_mongodb_url1="xxx"
``` 
if one of listed secret's path doesn't exist in Vault - entrypoint will fail.
