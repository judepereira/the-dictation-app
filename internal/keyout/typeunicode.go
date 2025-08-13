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
*/
import "C"
import "unicode/utf16"

func EnsureAccessibility() bool { return bool(C.ensureAccessibility()) }

func TypeString(s string) {
	u := utf16.Encode([]rune(s))
	if len(u) == 0 {
		return
	}
	C.type_utf16((*C.UniChar)(&u[0]), C.CFIndex(len(u)))
}
