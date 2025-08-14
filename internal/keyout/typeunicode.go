// internal/keyout/typeunicode.go
//go:build darwin

package keyout

/*
#cgo LDFLAGS: -framework Cocoa -framework ApplicationServices
#include <ApplicationServices/ApplicationServices.h>
#include <CoreFoundation/CoreFoundation.h>

static bool ensureAccessibility() {
  const void *keys[] = { kAXTrustedCheckOptionPrompt };
  const void *vals[] = { kCFBooleanTrue };
  CFDictionaryRef opts = CFDictionaryCreate(kCFAllocatorDefault, keys, vals, 1,
                                            &kCFTypeDictionaryKeyCallBacks,
                                            &kCFTypeDictionaryValueCallBacks);
  bool trusted = AXIsProcessTrustedWithOptions(opts);
  CFRelease(opts);
  return trusted;
}

static void type_utf16(const UniChar *buf, CFIndex len) {
  CGEventRef down = CGEventCreateKeyboardEvent(NULL, 0, true);
  CGEventKeyboardSetUnicodeString(down, len, buf);
  CGEventPost(kCGHIDEventTap, down);
  CFRelease(down);

  CGEventRef up = CGEventCreateKeyboardEvent(NULL, 0, false);
  CGEventKeyboardSetUnicodeString(up, len, buf);
  CGEventPost(kCGHIDEventTap, up);
  CFRelease(up);
}

// Backspace key virtual code on macOS
#define KEY_BACKSPACE ((CGKeyCode)0x33)

static void press_key(CGKeyCode code) {
  CGEventRef down = CGEventCreateKeyboardEvent(NULL, code, true);
  CGEventPost(kCGHIDEventTap, down);
  CFRelease(down);

  CGEventRef up = CGEventCreateKeyboardEvent(NULL, code, false);
  CGEventPost(kCGHIDEventTap, up);
  CFRelease(up);
}

static void backspace_n(int n) {
  for (int i = 0; i < n; i++) {
    press_key(KEY_BACKSPACE);
  }
}
*/
import "C"
import (
	"sync"
	"unicode/utf16"
)

func EnsureAccessibility() bool { return bool(C.ensureAccessibility()) }

var (
	mu        sync.Mutex
	prevRunes int
)

func TypeString(s string) {
	runes := []rune(s)

	mu.Lock()
	// Remove previously typed content
	if prevRunes > 0 {
		C.backspace_n(C.int(prevRunes))
	}

	// Type the current full transcript
	if len(runes) > 0 {
		u := utf16.Encode(runes)
		C.type_utf16((*C.UniChar)(&u[0]), C.CFIndex(len(u)))
	}

	prevRunes = len(runes)
	mu.Unlock()
}

func Stop() {
	mu.Lock()
	prevRunes = 0
	mu.Unlock()
}

func Start() {
	Stop()
}
