//go:build windows

package main

import (
	"bytes"
	"testing"
)

func TestEmbeddedApplicationPayload(t *testing.T) {
	payload, err := payloadFS.ReadFile("payload/" + appExeName)
	if err != nil {
		t.Fatal(err)
	}
	if len(payload) < 1024*1024 {
		t.Fatalf("embedded application is unexpectedly small: %d bytes", len(payload))
	}
	if !bytes.HasPrefix(payload, []byte{'M', 'Z'}) {
		t.Fatal("embedded application is not a Windows executable")
	}
}

func TestEmbeddedWebView2Bootstrapper(t *testing.T) {
	payload, err := payloadFS.ReadFile("payload/MicrosoftEdgeWebview2Setup.exe")
	if err != nil {
		t.Fatal(err)
	}
	if len(payload) < 1024*1024 {
		t.Fatalf("embedded WebView2 bootstrapper is unexpectedly small: %d bytes", len(payload))
	}
	if !bytes.HasPrefix(payload, []byte{'M', 'Z'}) {
		t.Fatal("embedded WebView2 bootstrapper is not a Windows executable")
	}
}
