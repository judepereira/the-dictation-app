//go:build darwin

package hotkey

/*
#cgo darwin LDFLAGS: -framework ApplicationServices -framework CoreFoundation
#include <ApplicationServices/ApplicationServices.h>

static CFMachPortRef gTap = NULL;

extern void handleOptionStateChange(int isDown);

// Event tap callback. It listens for modifier flag changes, filters Option key,
// and forwards down/up state to Go. It also re-enables the tap if it gets disabled.
static CGEventRef tapCallback(CGEventTapProxy proxy, CGEventType type, CGEventRef event, void *refcon) {
	if (type == kCGEventTapDisabledByTimeout || type == kCGEventTapDisabledByUserInput) {
		if (gTap != NULL) {
			CGEventTapEnable(gTap, true);
		}
		return event;
	}

	if (type != kCGEventFlagsChanged) {
		return event;
	}

	int64_t keycode = CGEventGetIntegerValueField(event, kCGKeyboardEventKeycode);
	// Left Option = 58, Right Option = 61
	if (keycode == 58 || keycode == 61) {
		CGEventFlags flags = CGEventGetFlags(event);
		int isDown = ((flags & kCGEventFlagMaskAlternate) != 0) ? 1 : 0;
		handleOptionStateChange(isDown);
	}

	return event;
}

// Starts a session event tap for flags-changed events and attaches it to the current thread's run loop.
static int startEventTap() {
	if (gTap != NULL) {
		return 1;
	}
	CGEventMask mask = (1ULL << kCGEventFlagsChanged);
	gTap = CGEventTapCreate(kCGSessionEventTap,
	                        kCGHeadInsertEventTap,
	                        kCGEventTapOptionListenOnly,
	                        mask,
	                        tapCallback,
	                        NULL);
	if (!gTap) {
		return 0;
	}
	CFRunLoopSourceRef src = CFMachPortCreateRunLoopSource(kCFAllocatorDefault, gTap, 0);
	if (!src) {
		CFRelease(gTap);
		gTap = NULL;
		return 0;
	}
	CFRunLoopAddSource(CFRunLoopGetCurrent(), src, kCFRunLoopCommonModes);
	CGEventTapEnable(gTap, true);
	CFRelease(src);
	return 1;
}
*/
import "C"

import (
	"log"
	"runtime"
	"time"
)

var (
	toggleCh        chan<- bool
	stateToggled    bool
	lastOptDown     bool
	lastTap         time.Time
	doubleTapWindow = 500 * time.Millisecond
)

// handleOptionStateChange is called from the C event tap callback with 1 on key down and 0 on key up.
//
//export handleOptionStateChange
func handleOptionStateChange(cIsDown C.int) {
	isDown := cIsDown != 0

	// Detect a "tap" on Option as a down->up transition.
	if lastOptDown && !isDown {
		now := time.Now()
		if !lastTap.IsZero() && now.Sub(lastTap) <= doubleTapWindow {
			// Double tap detected: toggle state and notify.
			lastTap = time.Time{} // reset to avoid triple-triggering
			stateToggled = !stateToggled
			select {
			case toggleCh <- stateToggled:
			default:
				// Avoid blocking if no receiver is ready.
			}
			log.Printf("Double Option tap detected, toggled state to %v", stateToggled)
		} else {
			// First tap; start/refresh the window.
			lastTap = now
		}
	}
	lastOptDown = isDown
}

// Run installs a CGEventTap listening for a double Option-key tap.
// On each detected double-tap it logs the event and sends the toggled state on the provided channel.
// This function blocks the calling goroutine to keep the event run loop alive.
func Run(toggle chan<- bool) {
	toggleCh = toggle

	// Run the event tap on a locked OS thread with its own CFRunLoop.
	go func() {
		runtime.LockOSThread()
		if C.startEventTap() == 0 {
			log.Printf("Failed to start macOS event tap for Option-key detection")
			return
		}
		C.CFRunLoopRun()
	}()

	// Block forever to keep the service alive.
	select {}
}
