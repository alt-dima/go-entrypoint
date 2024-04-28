package utils

import vault "github.com/hashicorp/vault/api"

type vaultStruct struct {
	vaultClient *vault.Client
}
