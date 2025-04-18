package utils

import (
	"encoding/json"
	"os"
	"regexp"
	"runtime"

	"github.com/PaesslerAG/jsonpath"
)

func (vaultStruct *vaultStruct) getSecretsEnvList(secretsFromSource *map[string]string) []string {

	// Convert the secrets to a slice of strings.
	secretsEnvList := []string{}

	if len(*secretsFromSource) > 0 {
		secretsEnvMap := make(map[string]string)

		for secretPath, secretEnvVarNamePrefix := range *secretsFromSource {
			getSecret := vaultStruct.readSecretVault(secretPath)
			for secretKey, secretValue := range getSecret {
				finalSecretEnvVar := secretEnvVarNamePrefix + "_" + secretKey
				Logger.Debug("Entrypoint got secret " + finalSecretEnvVar)
				secretsEnvList = append(secretsEnvList, finalSecretEnvVar+"="+secretValue.(string))
				secretsEnvMap[finalSecretEnvVar] = secretValue.(string)
			}
		}

		if runtime.GOOS == "windows" {
			addSecretsToReg(secretsEnvMap)
		}
	}

	return secretsEnvList
}

func readSecretsfromSource(secretsFile string, secretsJsonPath string) *map[string]string {

	// invalidRegexp is a regexp for invalid characters in keys
	var invalidRegexp = regexp.MustCompile(`[^a-zA-Z0-9_]`)

	// Open menifest json from platform inventory folder
	piManifestByte, err := os.ReadFile(secretsFile)
	if err != nil {
		Logger.Error("entrypoint failed read json with secrets: " + err.Error())
		os.Exit(1)
	}

	var secretsStrings = make(map[string]string)

	// Get the path to the secrets field.
	path := "$." + secretsJsonPath

	var manifest interface{}
	// Unmarshal the json file.
	err = json.Unmarshal(piManifestByte, &manifest)
	if err != nil {
		Logger.Error("entrypoint failed to unmarshal json with secrets: " + err.Error())
		os.Exit(1)
	}
	// Get the secrets from the manifest.
	secrets, err := jsonpath.Get(path, manifest)
	if err != nil {
		Logger.Warn("Entrypoint failed with specified path in json: " + err.Error())
	} else {

		for _, value := range secrets.([]interface{}) {
			switch v := value.(type) {
			case string:
				fixedSecretPath := invalidRegexp.ReplaceAllString(v, "_")
				//fixedSecretPath := strings.ReplaceAll(v, "/", "_")
				//fixedSecretPath = strings.ReplaceAll(fixedSecretPath, "-", "_")
				finalSecretEnvVarPrefix := "secret_" + fixedSecretPath
				secretsStrings[v] = finalSecretEnvVarPrefix
			case map[string]interface{}:
				if secretPathString, ok := v["secretname"].(string); ok {
					var finalSecretEnvVarPrefix string

					if secretnameString, ok := v["envvarname"].(string); ok {
						finalSecretEnvVarPrefix = secretnameString
					} else {
						fixedSecretPath := invalidRegexp.ReplaceAllString(secretPathString, "_")
						//fixedSecretPath := strings.ReplaceAll(secretPathString, "/", "_")
						//fixedSecretPath = strings.ReplaceAll(fixedSecretPath, "-", "_")
						finalSecretEnvVarPrefix = "secret_" + fixedSecretPath
					}

					secretsStrings[secretPathString] = finalSecretEnvVarPrefix

				} else {
					Logger.Error("Entrypoint wrong secrets list, secretname does not exists or not string")
					os.Exit(1)
				}
			default:
				Logger.Error("Entrypoint wrong secrets list")
				os.Exit(1)
			}
		}

	}

	return &secretsStrings
}
