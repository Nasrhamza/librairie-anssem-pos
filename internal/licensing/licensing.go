package licensing

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/base32"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"unicode"
)

const (
	ProductID      = "Jamel-v1"
	LicenseVersion = 1
	keyPrefix      = "JML1"
)

type Payload struct {
	Version    int    `json:"version"`
	Product    string `json:"product"`
	Customer   string `json:"customer"`
	DeviceCode string `json:"device_code"`
	IssuedAt   string `json:"issued_at"`
}

func DeviceCodeFromSecret(secret []byte) string {
	sum := sha256.Sum256(secret)
	compact := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(sum[:15])
	return DisplayDeviceCode(compact)
}

func NormalizeDeviceCode(value string) string {
	var out strings.Builder
	for _, r := range strings.ToUpper(strings.TrimSpace(value)) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			out.WriteRune(r)
		}
	}
	return out.String()
}

func DisplayDeviceCode(value string) string {
	compact := NormalizeDeviceCode(value)
	if compact == "" {
		return ""
	}
	parts := make([]string, 0, (len(compact)+3)/4)
	for len(compact) > 4 {
		parts = append(parts, compact[:4])
		compact = compact[4:]
	}
	parts = append(parts, compact)
	return strings.Join(parts, "-")
}

func Sign(seed []byte, customer, deviceCode string, issuedAt time.Time) (string, Payload, error) {
	customer = strings.TrimSpace(customer)
	deviceCode = NormalizeDeviceCode(deviceCode)
	if len(seed) != ed25519.SeedSize {
		return "", Payload{}, fmt.Errorf("clé privée invalide")
	}
	if len(customer) < 2 || len(customer) > 120 {
		return "", Payload{}, fmt.Errorf("nom client invalide")
	}
	if len(deviceCode) != 24 {
		return "", Payload{}, fmt.Errorf("code machine invalide")
	}
	payload := Payload{
		Version:    LicenseVersion,
		Product:    ProductID,
		Customer:   customer,
		DeviceCode: deviceCode,
		IssuedAt:   issuedAt.UTC().Format(time.RFC3339),
	}
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return "", Payload{}, err
	}
	privateKey := ed25519.NewKeyFromSeed(seed)
	signature := ed25519.Sign(privateKey, payloadJSON)
	return keyPrefix + "." + base64.RawURLEncoding.EncodeToString(payloadJSON) + "." + base64.RawURLEncoding.EncodeToString(signature), payload, nil
}

func Verify(publicKey []byte, key, expectedDeviceCode string) (Payload, error) {
	if len(publicKey) != ed25519.PublicKeySize {
		return Payload{}, fmt.Errorf("clé publique invalide")
	}
	parts := strings.Split(strings.TrimSpace(key), ".")
	if len(parts) != 3 || parts[0] != keyPrefix {
		return Payload{}, fmt.Errorf("format de licence invalide")
	}
	payloadJSON, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return Payload{}, fmt.Errorf("contenu de licence invalide")
	}
	signature, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil || len(signature) != ed25519.SignatureSize {
		return Payload{}, fmt.Errorf("signature de licence invalide")
	}
	if !ed25519.Verify(ed25519.PublicKey(publicKey), payloadJSON, signature) {
		return Payload{}, fmt.Errorf("licence non authentique")
	}
	var payload Payload
	if err := json.Unmarshal(payloadJSON, &payload); err != nil {
		return Payload{}, fmt.Errorf("données de licence invalides")
	}
	if payload.Version != LicenseVersion || payload.Product != ProductID {
		return Payload{}, fmt.Errorf("licence incompatible avec ce logiciel")
	}
	if strings.TrimSpace(payload.Customer) == "" {
		return Payload{}, fmt.Errorf("client de licence invalide")
	}
	if _, err := time.Parse(time.RFC3339, payload.IssuedAt); err != nil {
		return Payload{}, fmt.Errorf("date de licence invalide")
	}
	expected := NormalizeDeviceCode(expectedDeviceCode)
	if expected == "" || NormalizeDeviceCode(payload.DeviceCode) != expected {
		return Payload{}, fmt.Errorf("cette licence appartient à un autre ordinateur")
	}
	return payload, nil
}
