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

@property (weak) IBOutlet NSMenuItem *runningStatusMenuItem;
@property (weak) IBOutlet NSMenuItem *toggleRunningMenuItem;

@property (weak) IBOutlet NSMenuItem *autoModeMenuItem;
@property (weak) IBOutlet NSMenuItem *globalModeMenuItem;
@property (weak) IBOutlet NSMenuItem *manualModeMenuItem;



@property (strong) IBOutlet NSWindow *window;
@end

@implementation AppDelegate

#pragma mark menuAction

-(void)updateMainMenu{
    NSUserDefaults *shared = [NSUserDefaults standardUserDefaults];
    BOOL on = [shared boolForKey:@"DaggerOn"];
    if (on) {
        
        [_runningStatusMenuItem setTitle:@"Dagger: On"];
        [_runningStatusMenuItem setImage:[NSImage imageNamed:@"NSStatusAvailable"]];
        
        [_toggleRunningMenuItem setTitle:@"Turn Dagger Off"];
        
    } else {
        
        [_runningStatusMenuItem setTitle:@"Dagger: Off"];
        [_runningStatusMenuItem setImage:[NSImage imageNamed:@"NSStatusNone"]];
        
        [_toggleRunningMenuItem setTitle:@"Turn Dagger On"];
    }
    
    [self updateStatusMenuImage];
}

-(void)updateRunningModeMenu {
    NSUserDefaults *shared = [NSUserDefaults standardUserDefaults];
    NSString *mode = [shared objectForKey:@"DaggerMode"];
    
    [_autoModeMenuItem setState:NSControlStateValueOff];
    [_globalModeMenuItem setState:NSControlStateValueOff];
    [_manualModeMenuItem setState:NSControlStateValueOff];
    
    if ([mode isEqualTo:@"auto"]){
        [_autoModeMenuItem setState:NSControlStateValueOn];
    } else if ([mode isEqualTo:@"global"]){
        [_globalModeMenuItem setState:NSControlStateValueOn];
    } else if ([mode isEqualTo:@"manual"]){
        [_manualModeMenuItem setState:NSControlStateValueOn];
    }
    
    [self updateStatusMenuImage];
}

-(void)updateStatusMenuImage {
    NSUserDefaults *shared = [NSUserDefaults standardUserDefaults];
    BOOL on = [shared boolForKey:@"DaggerOn"];
    NSString *mode = [shared objectForKey:@"DaggerMode"];

    if (on){
        if ([mode isEqualTo:@"auto"]){
            [statusBarItem setImage:[NSImage imageNamed:@"menu_p_icon"]];
            [statusBarItem setAlternateImage:[NSImage imageNamed:@"menu_p_icon"]];
        } else if ([mode isEqualTo:@"global"]){
            [statusBarItem setImage:[NSImage imageNamed:@"menu_g_icon"]];
            [statusBarItem setAlternateImage:[NSImage imageNamed:@"menu_g_icon"]];
        } else if ([mode isEqualTo:@"manual"]){
            [statusBarItem setImage:[NSImage imageNamed:@"menu_m_icon"]];
            [statusBarItem setAlternateImage:[NSImage imageNamed:@"menu_m_icon"]];
        }
        [statusBarItem.image setTemplate:NO];
    } else {
        [statusBarItem setImage:[NSImage imageNamed:@"dagger"]];
        [statusBarItem setAlternateImage:[NSImage imageNamed:@"dagger"]];
        [statusBarItem.image setTemplate:NO];
    }
}

-(void)applyConf{
    NSUserDefaults *shared = [NSUserDefaults standardUserDefaults];
    NSString *mode = [shared objectForKey:@"DaggerMode"];
    BOOL on = [shared boolForKey:@"DaggerOn"];
    
    if (on) {
        if ([mode isEqualTo:@"auto"]){
            [ProxyConfHelper enablePACProxy];
        } else if ([mode isEqualTo:@"global"]){
            [ProxyConfHelper enableGlobalProxy];
        } else if ([mode isEqualTo:@"manual"]){
            [ProxyConfHelper disableProxy];
        }
    } else {
        [ProxyConfHelper disableProxy];
    }
}

- (IBAction)selectPACMode:(id)sender {
    [[NSUserDefaults standardUserDefaults] setObject:@"auto" forKey:@"DaggerMode"];
    [self updateRunningModeMenu];
    [self applyConf];
}

- (IBAction)selectGlobalMode:(id)sender {
    [[NSUserDefaults standardUserDefaults] setObject:@"global" forKey:@"DaggerMode"];
    [self updateRunningModeMenu];
    [self applyConf];
}

- (IBAction)selectManualMode:(id)sender {
    [[NSUserDefaults standardUserDefaults] setObject:@"manual" forKey:@"DaggerMode"];
    [self updateRunningModeMenu];
    [self applyConf];
}

- (IBAction)toggleRunning:(id)sender {
    NSUserDefaults *shared = [NSUserDefaults standardUserDefaults];
    BOOL on = [shared boolForKey:@"DaggerOn"];
    [shared setBool:!on forKey:@"DaggerOn"];
    
    [self updateMainMenu];

}

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
        @"DaggerOn":@NO,
        @"DaggerMode":@"auto",
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
    
    
    
    [self updateMainMenu];
    [self updateRunningModeMenu];
    [self applyConf];
}


- (void)applicationWillTerminate:(NSNotification *)aNotification {
}


- (BOOL)applicationSupportsSecureRestorableState:(NSApplication *)app {
    return YES;
}


@end
