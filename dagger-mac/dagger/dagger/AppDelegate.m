//
//  AppDelegate.m
//  dagger
//
//  Created by midoks on 2021/10/24.
//

#import "AppDelegate.h"
#import "Preferences/Preferences.h"

@interface AppDelegate ()
{
    Preferences *_preferenceWindow;
}


@property (strong) IBOutlet NSWindow *window;
@end

@implementation AppDelegate

#pragma mark 设置界面UI
-(void)setBarStatus
{
    statusBarItem = [[NSStatusBar systemStatusBar] statusItemWithLength:23.0];
    statusBarItem.image = [NSImage imageNamed:@"dagger"];
    statusBarItem.alternateImage = [NSImage imageNamed:@"dagger"];
    statusBarItem.menu = statusBarItemMenu;
    statusBarItem.toolTip = @"dagger";
    [statusBarItem setHighlightMode:YES];
}

-(void)initConfig{
    NSUserDefaults *shared = [NSUserDefaults standardUserDefaults];
    
    [shared registerDefaults:@{
        @"LocalSocks5.ListenPort":@1096,
        @"LocalSocks5.ListenAddress": @"127.0.0.1",
        @"GFWListURL": @"https://cdn.jsdelivr.net/gh/gfwlist/gfwlist/gfwlist.txt",
        @"ProxyExceptions": @"127.0.0.1, localhost, 192.168.0.0/16, 10.0.0.0/8, FE80::/64, ::1, FD00::/8",
    }];
    
    
    [shared setBool:NO forKey:@"_UIConstraintBasedLayoutLogUnsatisfiable"];
}

-(void)initPrefencesWindow
{
    _preferenceWindow = [[Preferences alloc] init];
    _preferenceWindow.window.level = NSFloatingWindowLevel;
}

#pragma mark Preferences
- (IBAction)showPreferences:(id)sender {

    [_preferenceWindow close];
    [_preferenceWindow showWindow:self];
    
    [NSApp activateIgnoringOtherApps:YES];
    [_preferenceWindow.window makeKeyAndOrderFront:sender];
    NSLog(@"ww:%@",_preferenceWindow);
}

- (void)applicationDidFinishLaunching:(NSNotification *)aNotification {
    [self initConfig];
    [self initPrefencesWindow];
    [self setBarStatus];
}


- (void)applicationWillTerminate:(NSNotification *)aNotification {
}


- (BOOL)applicationSupportsSecureRestorableState:(NSApplication *)app {
    return YES;
}


@end
