// internal/ui/about_darwin.go
//go:build darwin

package ui

/*
#cgo LDFLAGS: -framework Cocoa
void ShowAboutWindow();
*/
import "C"

func showAboutNative() { C.ShowAboutWindow() }
