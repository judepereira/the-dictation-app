// internal/ui/setup_darwin.go
//go:build darwin

package ui

/*
#cgo LDFLAGS: -framework Cocoa
void ShowSetupWindow();
*/
import "C"

func showSetupNative() { C.ShowSetupWindow() }
