package utils

import (
	"context"
	"os"

	vault "github.com/hashicorp/vault/api"
	auth "github.com/hashicorp/vault/api/auth/approle"
)

func (VaultStruct *vaultStruct) initVaultClient(vaultAddr string) {
	config := vault.DefaultConfig()
	config.Address = vaultAddr

	client, err := vault.NewClient(config)
	if err != nil {
		Logger.Error("Entrypoint failed create Vault client: " + err.Error())
		os.Exit(1)
	}

	appRoleAuth, err := auth.NewAppRoleAuth(
		os.Getenv("VAULT_APPROLE_RID"),
		&auth.SecretID{FromEnv: "VAULT_APPROLE_SID"},
	)
	if err != nil {
		Logger.Error("Entrypoint failed to create AppRoleAuth: " + err.Error())
		os.Exit(1)
	}

	authInfo, err := client.Auth().Login(context.Background(), appRoleAuth)
	if err != nil {
		Logger.Error("Entrypoint failed to login to Vault with appRoleAuth: " + err.Error())
		os.Exit(1)
	}
	if authInfo == nil {
		Logger.Error("Entrypoint failed empty authInfo: " + err.Error())
		os.Exit(1)
	}

	VaultStruct.vaultClient = client
}

func (vaultStruct *vaultStruct) readSecretVault(secretPath string) map[string]interface{} {

	secretResp, err := vaultStruct.vaultClient.KVv1("secret").Get(context.Background(), secretPath) //client.Read(context.Background(), "/secret/"+secretName)
	if err != nil {
		Logger.Error("Entrypoint failed to get secret from Vault: " + err.Error())
		os.Exit(1)
	}

	data := secretResp.Data

	return data
}
