package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func GenerateChildEnvs() []string {
	vaultAddr := os.Getenv("VAULT_ADDR")
	removeEnvVars := []string{"VAULT_APPROLE_RID", "VAULT_APPROLE_SID"}
	var childEnvs []string
	for _, removeEnv := range removeEnvVars {
		err := os.Unsetenv(removeEnv)
		if err != nil {
			Logger.Error("Entrypoint failed to unset env var: " + err.Error())
		}
	}
	childEnvs = os.Environ()

	if vaultAddr != "" {
		checkCriticalSvcReady(vaultAddr + "/v1/sys/health?standbyok=true")

		var vaultStruct vaultStruct
		vaultStruct.initVaultClient(vaultAddr)

		secretSourceConfig := os.Getenv("SECRETS_SOURCE_CONFIG")

		secretsFile := "./secrets_config.json"
		secretsJsonPath := "secrets_list"

		if secretSourceConfig != "" {
			splitedSourceConfig := strings.SplitN(secretSourceConfig, "#", 2)
			secretsFile = splitedSourceConfig[0]
			secretsJsonPath = splitedSourceConfig[1]
		}

		secretsFromSource := readSecretsfromSource(secretsFile, secretsJsonPath)
		secretsEnvList := vaultStruct.getSecretsEnvList(secretsFromSource)

		childEnvs = append(childEnvs, secretsEnvList...)
	}

	return childEnvs
}

func checkCriticalSvcReady(addrToCheck string) {
	retryCnt := 5
	waitTime := 4
	for {
		resp, err := http.Get(addrToCheck)
		if err == nil {
			_, errCoopy := io.Copy(io.Discard, resp.Body)
			if errCoopy != nil {
				Logger.Error("Entrypoint failed to read response body: " + err.Error())
			}
			errClose := resp.Body.Close()
			if errClose != nil {
				Logger.Error("Entrypoint failed to close response body: " + errClose.Error())
			}
		}
		if err == nil && resp.StatusCode == http.StatusOK {
			return
		} else if retryCnt < 1 {
			//stopExtSvcs()
			var reqError string
			if err != nil {
				reqError = err.Error()
			} else {
				reqError = fmt.Sprint(resp.StatusCode)
			}
			Logger.Error("Entrypoint critical svc check not passed: " + reqError)
			os.Exit(1)
		}
		retryCnt--
		time.Sleep(time.Duration(waitTime) * time.Second)
	}
}
