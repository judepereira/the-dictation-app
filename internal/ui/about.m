//go:build darwin
#import <Cocoa/Cocoa.h>

void ShowAboutWindow() {
  dispatch_async(dispatch_get_main_queue(), ^{
    NSAlert *alert = [[NSAlert alloc] init];
    alert.messageText = @"The Dictation App";
    alert.informativeText = @"Version 1.0.0\nSupport: https://judepereira.com";
    [alert addButtonWithTitle:@"Close"];
    [alert runModal];
  });
}

