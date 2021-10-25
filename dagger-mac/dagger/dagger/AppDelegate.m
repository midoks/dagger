//
//  AppDelegate.m
//  dagger
//
//  Created by midoks on 2021/10/24.
//

#import "AppDelegate.h"
#import "ProxyConfHelper.h"
#import "MASPreferences.h"
#import "Preferences.h"

@interface AppDelegate ()
{
    NSWindowController *_preferenceWindow;
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
        @"LocalSocks5.ListenPort": @"1096",
        @"LocalSocks5.ListenAddress": @"127.0.0.1",
        @"PacServer.BindToLocalhost": @YES,
        @"PacServer.ListenPort":@"1099",
        @"LocalSocks5.Timeout": @"60",
        @"LocalSocks5.EnableUDPRelay": @NO,
        @"LocalSocks5.EnableVerboseMode": @NO,
        @"GFWListURL": @"https://cdn.jsdelivr.net/gh/gfwlist/gfwlist/gfwlist.txt",
        @"AutoConfigureNetworkServices":@YES,
        @"ProxyExceptions": @"127.0.0.1, localhost, 192.168.0.0/16, 10.0.0.0/8, FE80::/64, ::1, FD00::/8",
        
    }];

}

-(void)initPrefencesWindow
{
    
    NSArray *listVC = @[
        [[PreferencesGeneral alloc] init],
        [[PreferencesAdvanced alloc] init],
        [[PreferencesInterfaces alloc] init],
    ];
    
    _preferenceWindow = [[MASPreferencesWindowController alloc] initWithViewControllers:listVC title:@""];
    _preferenceWindow.window.level = NSFloatingWindowLevel;
    
}

#pragma mark Preferences
- (IBAction)showPreferences:(id)sender {
    [_preferenceWindow showWindow:self];
}

- (void)applicationDidFinishLaunching:(NSNotification *)aNotification {
    [self initConfig];
    [self initPrefencesWindow];
    [self setBarStatus];
    
    [ProxyConfHelper install];
}


- (void)applicationWillTerminate:(NSNotification *)aNotification {
}


- (BOOL)applicationSupportsSecureRestorableState:(NSApplication *)app {
    return YES;
}


@end
