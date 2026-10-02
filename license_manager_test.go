//go:build windows

package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"librairie-anssem-pos-portable/internal/licensing"
)

func TestLicenseActivationIsUniquePerInstallation(t *testing.T) {
	seed, err := os.ReadFile(filepath.Join("owner-tools", "license-generator", "private", "jamel_ed25519_seed.key"))
	if err != nil {
		t.Fatal(err)
	}

	firstRoot := filepath.Join(t.TempDir(), "first")
	t.Setenv("LOCALAPPDATA", firstRoot)
	first, err := newLicenseManager()
	if err != nil {
		t.Fatal(err)
	}
	firstStatus := first.status()
	if firstStatus.Activated || firstStatus.DeviceCode == "" {
		t.Fatalf("unexpected initial status: %#v", firstStatus)
	}
	key, _, err := licensing.Sign(seed, "Client Autorisé", firstStatus.DeviceCode, time.Now())
	if err != nil {
		t.Fatal(err)
	}
	activated, err := first.activate(key)
	if err != nil {
		t.Fatal(err)
	}
	if !activated.Activated || activated.Customer != "Client Autorisé" {
		t.Fatalf("unexpected activated status: %#v", activated)
	}

	secondRoot := filepath.Join(t.TempDir(), "second")
	t.Setenv("LOCALAPPDATA", secondRoot)
	second, err := newLicenseManager()
	if err != nil {
		t.Fatal(err)
	}
	if second.deviceCode == first.deviceCode {
		t.Fatal("separate installations received the same device code")
	}
	if _, err := second.activate(key); err == nil {
		t.Fatal("license from the first installation activated the second installation")
	}
}

func TestDeviceBindingRoundTrip(t *testing.T) {
	t.Setenv("LOCALAPPDATA", t.TempDir())
	first, err := loadOrCreateDeviceCode()
	if err != nil {
		t.Fatal(err)
	}
	second, err := loadOrCreateDeviceCode()
	if err != nil {
		t.Fatal(err)
	}
	if first != second {
		t.Fatalf("device code changed between launches: %q != %q", first, second)
	}
}
