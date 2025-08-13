// internal/ui/setup.m
//go:build darwin
#import <Cocoa/Cocoa.h>

void ShowSetupWindow() {
  dispatch_async(dispatch_get_main_queue(), ^{
    NSWindow *w = [[NSWindow alloc] initWithContentRect:NSMakeRect(0,0,420,180)
      styleMask:(NSWindowStyleMaskTitled|NSWindowStyleMaskClosable)
      backing:NSBackingStoreBuffered defer:NO];
    [w setTitle:@"First-time Setup"];
    NSProgressIndicator *p = [[NSProgressIndicator alloc] initWithFrame:NSMakeRect(20,80,380,20)];
    [p setIndeterminate:NO]; [p setMinValue:0]; [p setMaxValue:100];
    [[w contentView] addSubview:p];
    [w center]; [w makeKeyAndOrderFront:nil];
  });
}

