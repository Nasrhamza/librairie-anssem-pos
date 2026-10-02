package licensing

import (
	"crypto/ed25519"
	"crypto/rand"
	"strings"
	"testing"
	"time"
)

func TestSignedLicenseIsBoundToDevice(t *testing.T) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	deviceCode := DeviceCodeFromSecret([]byte("licensed computer secret"))
	key, signed, err := Sign(privateKey.Seed(), "Client Test", deviceCode, time.Date(2026, 9, 2, 12, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatal(err)
	}
	verified, err := Verify(publicKey, key, deviceCode)
	if err != nil {
		t.Fatal(err)
	}
	if verified.Customer != signed.Customer || verified.DeviceCode != NormalizeDeviceCode(deviceCode) {
		t.Fatalf("unexpected payload: %#v", verified)
	}
	otherDevice := DeviceCodeFromSecret([]byte("different computer secret"))
	if _, err := Verify(publicKey, key, otherDevice); err == nil || !strings.Contains(err.Error(), "autre ordinateur") {
		t.Fatalf("expected device mismatch, got %v", err)
	}
}

func TestTamperedLicenseIsRejected(t *testing.T) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	deviceCode := DeviceCodeFromSecret([]byte("device"))
	key, _, err := Sign(privateKey.Seed(), "Client Test", deviceCode, time.Now())
	if err != nil {
		t.Fatal(err)
	}
	tampered := key[:len(key)-1] + "A"
	if tampered == key {
		tampered = key[:len(key)-1] + "B"
	}
	if _, err := Verify(publicKey, tampered, deviceCode); err == nil {
		t.Fatal("tampered license was accepted")
	}
}
