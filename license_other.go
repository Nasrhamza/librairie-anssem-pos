//go:build !windows

package main

import (
	"crypto/rand"
	"fmt"
	"os"
	"path/filepath"

	"librairie-anssem-pos-portable/internal/licensing"
)

func loadOrCreateDeviceCode() (string, error) {
	path := filepath.Join(appDataDir(), "device.key")
	if secret, err := os.ReadFile(path); err == nil {
		if len(secret) != 32 {
			return "", fmt.Errorf("liaison de licence endommagée")
		}
		return licensing.DeviceCodeFromSecret(secret), nil
	} else if !os.IsNotExist(err) {
		return "", err
	}
	secret := make([]byte, 32)
	if _, err := rand.Read(secret); err != nil {
		return "", err
	}
	if err := os.WriteFile(path, secret, 0600); err != nil {
		return "", err
	}
	return licensing.DeviceCodeFromSecret(secret), nil
}
