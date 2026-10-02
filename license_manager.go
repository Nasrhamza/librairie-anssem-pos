package main

import (
	"embed"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"librairie-anssem-pos-portable/internal/licensing"
)

//go:embed licensing/public.key
var licensingAssets embed.FS

type LicenseStatus struct {
	Required   bool   `json:"required"`
	Activated  bool   `json:"activated"`
	DeviceCode string `json:"device_code"`
	Customer   string `json:"customer,omitempty"`
	IssuedAt   string `json:"issued_at,omitempty"`
	Error      string `json:"error,omitempty"`
}

type LicenseManager struct {
	mu         sync.Mutex
	path       string
	publicKey  []byte
	deviceCode string
}

func newLicenseManager() (*LicenseManager, error) {
	publicKey, err := licensingAssets.ReadFile("licensing/public.key")
	if err != nil {
		return nil, fmt.Errorf("clé publique absente: %w", err)
	}
	deviceCode, err := loadOrCreateDeviceCode()
	if err != nil {
		return nil, err
	}
	return &LicenseManager{
		path:       filepath.Join(appDataDir(), "license.key"),
		publicKey:  publicKey,
		deviceCode: deviceCode,
	}, nil
}

func (l *LicenseManager) status() LicenseStatus {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.statusLocked()
}

func (l *LicenseManager) statusLocked() LicenseStatus {
	status := LicenseStatus{Required: true, DeviceCode: l.deviceCode}
	key, err := os.ReadFile(l.path)
	if os.IsNotExist(err) {
		return status
	}
	if err != nil {
		status.Error = "Licence illisible"
		return status
	}
	payload, err := licensing.Verify(l.publicKey, string(key), l.deviceCode)
	if err != nil {
		status.Error = err.Error()
		return status
	}
	status.Activated = true
	status.Customer = payload.Customer
	status.IssuedAt = payload.IssuedAt
	return status
}

func (l *LicenseManager) activate(key string) (LicenseStatus, error) {
	key = strings.TrimSpace(key)
	if len(key) < 80 || len(key) > 4096 {
		return l.status(), fmt.Errorf("clé produit invalide")
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	payload, err := licensing.Verify(l.publicKey, key, l.deviceCode)
	if err != nil {
		return l.statusLocked(), err
	}
	if err := os.WriteFile(l.path, []byte(key+"\r\n"), 0600); err != nil {
		return l.statusLocked(), fmt.Errorf("impossible d'enregistrer la licence: %w", err)
	}
	return LicenseStatus{
		Required: true, Activated: true, DeviceCode: l.deviceCode,
		Customer: payload.Customer, IssuedAt: payload.IssuedAt,
	}, nil
}

func (l *LicenseManager) protect(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") &&
			!strings.HasPrefix(r.URL.Path, "/api/license/") &&
			!strings.HasPrefix(r.URL.Path, "/api/public/") &&
			!strings.HasPrefix(r.URL.Path, "/api/branding/") && r.URL.Path != "/api/health" &&
			!l.status().Activated {
			errJSON(w, http.StatusForbidden, "Activation du logiciel requise")
			return
		}
		next.ServeHTTP(w, r)
	})
}
