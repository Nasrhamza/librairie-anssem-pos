//go:build windows

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"syscall"
	"unsafe"

	webview "github.com/jchv/go-webview2"
)

var (
	desktopUser32      = syscall.NewLazyDLL("user32.dll")
	desktopShowWindow  = desktopUser32.NewProc("ShowWindow")
	desktopMessageBoxW = desktopUser32.NewProc("MessageBoxW")
)

func runDesktopWindow(url string) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	dataPath := filepath.Join(appDataDir(), "WebView2")
	if err := os.MkdirAll(dataPath, 0755); err != nil {
		return fmt.Errorf("préparation du profil desktop: %w", err)
	}
	w := webview.NewWithOptions(webview.WebViewOptions{
		Debug:     false,
		AutoFocus: true,
		DataPath:  dataPath,
		WindowOptions: webview.WindowOptions{
			Title:  "Jamel v1 — Librairie Anssem",
			Width:  1440,
			Height: 900,
			IconId: 2,
			Center: true,
		},
	})
	if w == nil {
		return fmt.Errorf("Microsoft WebView2 Runtime est absent ou indisponible. Relancez Jamel-v1-Setup.exe pour l'installer")
	}
	defer w.Destroy()
	w.SetSize(1050, 680, webview.HintMin)
	w.Navigate(url)
	desktopShowWindow.Call(uintptr(w.Window()), 3) // SW_MAXIMIZE
	w.Run()
	return nil
}

func showDesktopError(title, message string) {
	titlePtr, _ := syscall.UTF16PtrFromString(title)
	messagePtr, _ := syscall.UTF16PtrFromString(message)
	desktopMessageBoxW.Call(0, uintptr(unsafe.Pointer(messagePtr)), uintptr(unsafe.Pointer(titlePtr)), 0x10)
}
