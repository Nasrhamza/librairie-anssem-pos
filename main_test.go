package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestPBKDF2SHA256KnownVector(t *testing.T) {
	got := hex.EncodeToString(pbkdf2SHA256([]byte("password"), []byte("salt"), 1, 32))
	want := "120fb6cffcf8b32c43e7225256c4f837a86548c92ccc35480805987cb70be17b"
	if got != want {
		t.Fatalf("PBKDF2 result = %s, want %s", got, want)
	}
}

func TestDigits(t *testing.T) {
	got := digits("+216 97-981-347")
	if got != "21697981347" {
		t.Fatalf("digits() = %q", got)
	}
}

func TestDemoDataHasFiveCartsAfterEnsure(t *testing.T) {
	s := &Store{data: demoData()}
	s.ensureCarts()
	if len(s.data.Carts) != 5 {
		t.Fatalf("expected 5 carts, got %d", len(s.data.Carts))
	}
	for i, c := range s.data.Carts {
		if c.Slot != i+1 {
			t.Fatalf("cart slot %d = %d", i, c.Slot)
		}
	}
}

func TestAtomicSaveWritesValidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "data.json")
	s := &Store{path: path, data: demoData()}
	s.ensureCarts()
	if err := s.saveLocked(); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var d StoreData
	if err := json.Unmarshal(b, &d); err != nil {
		t.Fatalf("saved file is not valid JSON: %v", err)
	}
	if len(d.Products) == 0 || len(d.Carts) != 5 {
		t.Fatalf("saved data incomplete")
	}
}

func TestLegacyNaje7DataMigratesToLibrairieAnssem(t *testing.T) {
	root := t.TempDir()
	t.Setenv("LOCALAPPDATA", root)
	legacyDir := filepath.Join(root, "MaktabetNaje7")
	if err := os.MkdirAll(legacyDir, 0755); err != nil {
		t.Fatal(err)
	}
	legacyData := []byte(`{"products":[{"id":99,"name":"Produit conservé"}]}`)
	legacyAuth := []byte(`{"username":"admin"}`)
	if err := os.WriteFile(filepath.Join(legacyDir, "data.json"), legacyData, 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(legacyDir, "auth.json"), legacyAuth, 0600); err != nil {
		t.Fatal(err)
	}
	newDir := appDataDir()
	if filepath.Base(newDir) != "LibrairieAnssem" {
		t.Fatalf("new data directory = %q", newDir)
	}
	for name, want := range map[string][]byte{"data.json": legacyData, "auth.json": legacyAuth} {
		got, err := os.ReadFile(filepath.Join(newDir, name))
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(got, want) {
			t.Fatalf("migrated %s does not match legacy file", name)
		}
	}
}

func TestConnectionDefaults(t *testing.T) {
	d := demoData()
	if d.Settings.Conn500Price < 0 || d.Settings.Conn1000Price < 0 {
		t.Fatal("connection price cannot be negative")
	}
	if len(d.SIMs) == 0 || d.SIMs[0].QuotaMB != 1000 {
		t.Fatal("default SIM quota must be 1000 MB")
	}
	if !validThemeColor(d.Settings.Theme) {
		t.Fatalf("invalid default theme color: %q", d.Settings.Theme)
	}
}

func TestValidThemeColor(t *testing.T) {
	for _, color := range []string{"#1b805d", "#000000", "#abcdef"} {
		if !validThemeColor(color) {
			t.Fatalf("expected %s to be valid", color)
		}
	}
	for _, color := range []string{"1b805d", "#abcd", "#gg0000", ""} {
		if validThemeColor(color) {
			t.Fatalf("expected %q to be invalid", color)
		}
	}
}

func TestSIMUsageForDay(t *testing.T) {
	transfers := []Transfer{
		{SIMID: 1, Time: "2026-08-31T09:00:00+02:00", PackageMB: 500},
		{SIMID: 1, Time: "2026-08-31T10:00:00+02:00", PackageMB: 500},
		{SIMID: 1, Time: "2026-09-01T09:00:00+02:00", PackageMB: 1000},
		{SIMID: 2, Time: "2026-08-31T11:00:00+02:00", PackageMB: 1000},
	}
	if got := simUsageForDay(transfers, 1, "2026-08-31"); got != 1000 {
		t.Fatalf("usage = %d, want 1000", got)
	}
}

func TestDailySIMReportPDF(t *testing.T) {
	d := demoData()
	d.Transfers = []Transfer{{
		ID: 1, Time: "2026-08-31T09:30:00+02:00", SIMID: 1,
		SIMLabel: "SIM 01", DestPhone: "97981347", PackageMB: 500, Price: 2,
	}}
	pdf := buildDailySIMReportPDF(d, "2026-08-31")
	if !bytes.HasPrefix(pdf, []byte("%PDF-1.4")) {
		t.Fatal("daily report is not a PDF")
	}
	for _, want := range []string{"Rapport ventes SIM", "97981347", "TOTAL: 2.000 DT"} {
		if !bytes.Contains(pdf, []byte(want)) {
			t.Fatalf("daily report missing %q", want)
		}
	}
}

func TestSaleInvoicePDF(t *testing.T) {
	sale := Sale{
		ID: 7, Time: "2026-08-31T10:00:00+02:00", Payment: "cash",
		Received: 10, Change: 2, Total: 8,
		Items: []SaleItem{{ProductID: 1, Name: "Cahier 96 pages", Qty: 2, UnitPrice: 4, Total: 8}},
	}
	profile := defaultBusinessProfile()
	pdf := buildSaleInvoicePDF(sale, profile, storeLogoPNG)
	if !bytes.HasPrefix(pdf, []byte("%PDF-1.4")) {
		t.Fatal("invoice is not a PDF")
	}
	for _, want := range []string{"FACTURE CLIENT", "LIBRAIRIE ANSSEM", storeWhatsApp, storeEmail, storeTaxID, "Cahier 96 pages", "8.000 DT", "/Subtype /Image"} {
		if !bytes.Contains(pdf, []byte(want)) {
			t.Fatalf("invoice missing %q", want)
		}
	}
}

func TestBusinessAndDigitalServiceDefaults(t *testing.T) {
	d := StoreData{Products: demoData().Products, Settings: Settings{}}
	s := &Store{data: d}
	s.ensureDefaults()
	if s.data.Settings.Business.Name != storeDisplayName {
		t.Fatalf("business name = %q", s.data.Settings.Business.Name)
	}
	if len(s.data.DigitalServices) < 8 {
		t.Fatalf("expected default digital services, got %d", len(s.data.DigitalServices))
	}
}

func TestNormalizeDigitalServiceRequiresFullWebURL(t *testing.T) {
	valid := DigitalService{Name: "Portail", URL: "https://example.com/page", Icon: "PT"}
	if err := normalizeDigitalService(&valid); err != nil {
		t.Fatal(err)
	}
	invalid := DigitalService{Name: "Portail", URL: "javascript:alert(1)", Icon: "PT"}
	if err := normalizeDigitalService(&invalid); err == nil {
		t.Fatal("javascript URL must be rejected")
	}
}

func TestSaleInvoiceUsesCustomizedBusiness(t *testing.T) {
	sale := Sale{ID: 9, Time: "2026-09-02T10:00:00+02:00", Payment: "cash", Total: 1,
		Items: []SaleItem{{Name: "Test", Qty: 1, UnitPrice: 1, Total: 1}}}
	profile := defaultBusinessProfile()
	profile.Name = "Librairie Client Test"
	profile.Email = "client@example.com"
	pdf := buildSaleInvoicePDF(sale, profile, storeLogoPNG)
	for _, want := range []string{"LIBRAIRIE CLIENT TEST", "client@example.com"} {
		if !bytes.Contains(pdf, []byte(want)) {
			t.Fatalf("customized invoice missing %q", want)
		}
	}
}

func TestPinnedProductPersists(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "data.json")
	d := demoData()
	d.Products[0].Pinned = true
	s := &Store{path: path, data: d}
	s.ensureCarts()
	if err := s.saveLocked(); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var got StoreData
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatal(err)
	}
	if !got.Products[0].Pinned {
		t.Fatal("pinned product was not persisted")
	}
}
