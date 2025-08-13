// internal/ui/about.m
//go:build darwin
#import <Cocoa/Cocoa.h>

void ShowAboutWindow() {
  dispatch_async(dispatch_get_main_queue(), ^{
    NSAlert *alert = [[NSAlert alloc] init];
    alert.messageText = @"DictationApp";
    alert.informativeText = @"Version 0.1.0\nSupport: https://example.com/support";
    [alert addButtonWithTitle:@"OK"];
    [alert runModal];
  });
}

