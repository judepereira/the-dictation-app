//go:build darwin
#import <Cocoa/Cocoa.h>

static NSWindow *gSetupWindow = nil;
static NSProgressIndicator *gSetupProgress = nil;

void ShowSetupWindow() {
  dispatch_async(dispatch_get_main_queue(), ^{
    if (gSetupWindow != nil) {
      [gSetupWindow makeKeyAndOrderFront:nil];
      return;
    }
    gSetupWindow = [[NSWindow alloc] initWithContentRect:NSMakeRect(0,0,420,180)
      styleMask:(NSWindowStyleMaskTitled|NSWindowStyleMaskClosable)
      backing:NSBackingStoreBuffered defer:NO];
    [gSetupWindow setTitle:@"Downloading Model"];

    gSetupProgress = [[NSProgressIndicator alloc] initWithFrame:NSMakeRect(20,80,380,20)];
    [gSetupProgress setIndeterminate:NO];
    [gSetupProgress setMinValue:0];
    [gSetupProgress setMaxValue:100];
    [gSetupProgress setDoubleValue:0];

    [[gSetupWindow contentView] addSubview:gSetupProgress];
    [gSetupWindow center];
    [gSetupWindow makeKeyAndOrderFront:nil];
  });
}

// Determinate/indeterminate switch (for unknown file sizes)
void SetupSetIndeterminate(bool indeterminate) {
  dispatch_async(dispatch_get_main_queue(), ^{
    if (!gSetupProgress) return;
    [gSetupProgress setIndeterminate:indeterminate ? YES : NO];
    if (indeterminate) {
      [gSetupProgress startAnimation:nil];
    } else {
      [gSetupProgress stopAnimation:nil];
    }
  });
}

// Update progress in [0..100]
void SetupUpdateProgress(double percent) {
  dispatch_async(dispatch_get_main_queue(), ^{
    if (!gSetupProgress) return;
    if ([gSetupProgress isIndeterminate]) {
      [gSetupProgress setIndeterminate:NO];
      [gSetupProgress stopAnimation:nil];
    }
    [gSetupProgress setDoubleValue:percent];
  });
}

void CloseSetupWindow() {
  dispatch_async(dispatch_get_main_queue(), ^{
    if (gSetupWindow) {
      [gSetupWindow close];
      gSetupWindow = nil;
      gSetupProgress = nil;
    }
  });
}