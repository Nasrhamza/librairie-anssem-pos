//go:build windows

package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"time"
	"unsafe"

	"librairie-anssem-pos-portable/internal/licensing"
)

const (
	cryptProtectUIForbidden  = 0x1
	cryptProtectLocalMachine = 0x4
)

var (
	licenseCrypt32            = syscall.NewLazyDLL("crypt32.dll")
	licenseKernel32           = syscall.NewLazyDLL("kernel32.dll")
	licenseCryptProtectData   = licenseCrypt32.NewProc("CryptProtectData")
	licenseCryptUnprotectData = licenseCrypt32.NewProc("CryptUnprotectData")
	licenseLocalFree          = licenseKernel32.NewProc("LocalFree")
)

type licenseDataBlob struct {
	Size uint32
	Data *byte
}

var deviceBindingEntropy = []byte("Jamel-v1|LibrairieAnssem|device-binding|v1")

func loadOrCreateDeviceCode() (string, error) {
	path := filepath.Join(appDataDir(), "device.key")
	if encoded, err := os.ReadFile(path); err == nil {
		protected, err := base64.RawStdEncoding.DecodeString(string(encoded))
		if err == nil {
			secret, decryptErr := unprotectDeviceSecret(protected)
			if decryptErr == nil && len(secret) == 32 {
				return licensing.DeviceCodeFromSecret(secret), nil
			}
		}
		_ = os.Rename(path, fmt.Sprintf("%s.invalid-%d", path, time.Now().Unix()))
	} else if !os.IsNotExist(err) {
		return "", fmt.Errorf("lecture de la liaison machine: %w", err)
	}

	secret := make([]byte, 32)
	if _, err := rand.Read(secret); err != nil {
		return "", fmt.Errorf("création du code machine: %w", err)
	}
	protected, err := protectDeviceSecret(secret)
	if err != nil {
		return "", err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, []byte(base64.RawStdEncoding.EncodeToString(protected)), 0600); err != nil {
		return "", err
	}
	if err := os.Rename(tmp, path); err != nil {
		_ = os.Remove(tmp)
		return "", err
	}
	return licensing.DeviceCodeFromSecret(secret), nil
}

func protectDeviceSecret(secret []byte) ([]byte, error) {
	in := newLicenseBlob(secret)
	entropy := newLicenseBlob(deviceBindingEntropy)
	var out licenseDataBlob
	description, _ := syscall.UTF16PtrFromString("Jamel v1 device binding")
	ok, _, callErr := licenseCryptProtectData.Call(
		uintptr(unsafe.Pointer(&in)), uintptr(unsafe.Pointer(description)), uintptr(unsafe.Pointer(&entropy)),
		0, 0, cryptProtectUIForbidden|cryptProtectLocalMachine, uintptr(unsafe.Pointer(&out)),
	)
	if ok == 0 {
		return nil, fmt.Errorf("protection Windows du code machine: %w", callErr)
	}
	defer licenseLocalFree.Call(uintptr(unsafe.Pointer(out.Data)))
	return copyLicenseBlob(out), nil
}

func unprotectDeviceSecret(protected []byte) ([]byte, error) {
	in := newLicenseBlob(protected)
	entropy := newLicenseBlob(deviceBindingEntropy)
	var out licenseDataBlob
	ok, _, callErr := licenseCryptUnprotectData.Call(
		uintptr(unsafe.Pointer(&in)), 0, uintptr(unsafe.Pointer(&entropy)),
		0, 0, cryptProtectUIForbidden, uintptr(unsafe.Pointer(&out)),
	)
	if ok == 0 {
		return nil, callErr
	}
	defer licenseLocalFree.Call(uintptr(unsafe.Pointer(out.Data)))
	return copyLicenseBlob(out), nil
}

func newLicenseBlob(data []byte) licenseDataBlob {
	if len(data) == 0 {
		return licenseDataBlob{}
	}
	return licenseDataBlob{Size: uint32(len(data)), Data: &data[0]}
}

func copyLicenseBlob(blob licenseDataBlob) []byte {
	if blob.Size == 0 || blob.Data == nil {
		return nil
	}
	return append([]byte(nil), unsafe.Slice(blob.Data, int(blob.Size))...)
}
