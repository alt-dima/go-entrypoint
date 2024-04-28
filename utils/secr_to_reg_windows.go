package utils

import (
	"os"

	"golang.org/x/sys/windows/registry"
)

func addSecretsToReg(secrets map[string]string) {
	k, _, err := registry.CreateKey(registry.LOCAL_MACHINE, `System\CurrentControlSet\Control\Session Manager\Environment`, registry.CREATE_SUB_KEY|registry.SET_VALUE)
	if err != nil {
		Logger.Error("Entrypoint failed to create a key in registry: " + err.Error())
		os.Exit(1)
	}
	for name, secret := range secrets {
		if err := k.SetStringValue(name, secret); err != nil {
			Logger.Error("Entrypoint failed to set kv in registry: " + err.Error())
			os.Exit(1)
		}
	}
	if err := k.Close(); err != nil {
		Logger.Error("Entrypoint failed to close registry: " + err.Error())
		os.Exit(1)
	}
}
