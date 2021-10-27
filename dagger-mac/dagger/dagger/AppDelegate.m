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
#import "PACUtils.h"
#import "Servers.h"
#import "UserRules.h"
#import "LaunchAgentsUtils.h"

@interface AppDelegate () <NSUserNotificationCenterDelegate>
{
    NSWindowController *_preferenceWindow;
    UserRules *_userRuleWindow;
    Servers *_serverConf;
}

@property (weak) IBOutlet NSMenuItem *runningStatusMenuItem;
@property (weak) IBOutlet NSMenuItem *toggleRunningMenuItem;

@property (weak) IBOutlet NSMenuItem *autoModeMenuItem;
@property (weak) IBOutlet NSMenuItem *globalModeMenuItem;
@property (weak) IBOutlet NSMenuItem *manualModeMenuItem;

@property (weak) IBOutlet NSMenuItem *serverMenuItem;
@property (weak) IBOutlet NSMenuItem *serverBeginSeparatorMenuItem;
@property (weak) IBOutlet NSMenuItem *serverEndSeparatorMenuItem;




@property (strong) IBOutlet NSWindow *window;
@end

@implementation AppDelegate

#pragma mark 用户通知中心
- (BOOL)userNotificationCenter:(NSUserNotificationCenter *)center shouldPresentNotification:(NSUserNotification *)notification
{
    return YES;
}

-(void)Toast:(NSString *)content
{
    [[NSUserNotificationCenter defaultUserNotificationCenter] removeAllDeliveredNotifications];
    for (NSUserNotification *notify in [[NSUserNotificationCenter defaultUserNotificationCenter] scheduledNotifications])
    {
        [[NSUserNotificationCenter defaultUserNotificationCenter] removeScheduledNotification:notify];
    }
    
    
    NSUserNotification *notification = [[NSUserNotification alloc] init];
    notification.title = @"user notification";
    notification.informativeText = content;
    
    //设置通知的代理
    [[NSUserNotificationCenter defaultUserNotificationCenter] setDelegate:self];
    [[NSUserNotificationCenter defaultUserNotificationCenter] scheduleNotification:notification];
}

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

-(void)updateServersMenu
{
    NSMenu *menu = _serverMenuItem.submenu;
    NSInteger bIndex = [menu indexOfItem:_serverBeginSeparatorMenuItem]+1;
    NSInteger eIndex = [menu indexOfItem:_serverEndSeparatorMenuItem]-1;

    for (NSInteger mIndex = eIndex ; mIndex>=bIndex; mIndex--) {
        [menu removeItemAtIndex:mIndex];
    }
    
    NSMutableArray *slist  = [Servers serverList];
    for (NSInteger i=0; i<[slist count]; i++) {
        NSMutableDictionary *row = [slist objectAtIndex:i];
        NSMenuItem *item = [[NSMenuItem alloc] init];
        [item setTitle:[row valueForKey:@"remark"]];
        [item setEnabled:YES];
        item.tag = i;
        item.state = [[row valueForKey:@"status"] isEqualTo:@"on"]?NSControlStateValueOn:NSControlStateValueOff;
        
        [item setAction:@selector(selectServer:)];
        [menu insertItem:item atIndex:bIndex];
    }
}

-(void)selectServer:(NSMenuItem*)sender{
//    NSLog(@"selectServer");
    [Servers set:sender.tag value:@"on" forKey:@"status"];
    [self updateServersMenu];
    
    [[NSNotificationCenter defaultCenter] postNotificationName:@"changeConfigList" object:nil userInfo:nil];
}

-(void)applyConf{
    NSUserDefaults *shared = [NSUserDefaults standardUserDefaults];
    NSString *mode = [shared objectForKey:@"DaggerMode"];
    BOOL on = [shared boolForKey:@"DaggerOn"];
    
    if (on) {
        if ([mode isEqualTo:@"auto"]){
            [LaunchAgentsUtils startHttpProxy];
            [ProxyConfHelper enablePACProxy];
        } else if ([mode isEqualTo:@"global"]){
            [LaunchAgentsUtils startHttpProxy];
            [ProxyConfHelper enableGlobalProxy];
        } else if ([mode isEqualTo:@"manual"]){
            [ProxyConfHelper disableProxy];
            [LaunchAgentsUtils stopHttpProxy];
        }
    } else {
        [ProxyConfHelper disableProxy];
        [LaunchAgentsUtils stopHttpProxy];
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
    [self applyConf];
}

- (IBAction)showLog:(id)sender {
//    NSWorkspace *ws = [NSWorkspace sharedWorkspace];
    
//    NSURL *appUrl = [ws URLForApplicationWithBundleIdentifier:@"com.apple.Console"];
//    NSArray  *logArr = [NSArray arrayWithObjects:@"~/Library/Logs/dagger-client-http.log",nil];
//
//    NSMutableDictionary* dict = [[NSMutableDictionary alloc] init];
//    [dict setObject:logArr forKey:NSWorkspaceLaunchConfigurationArguments];
//
//    [ws launchApplicationAtURL:appUrl
//                       options:NSWorkspaceLaunchDefault
//                 configuration:dict
//                         error:nil];
    
    NSString *logPath = [NSString stringWithFormat:@"%@/%@",NSHomeDirectory(), @"Library/Logs/dagger-client-http.log"];
    
    [[NSTask launchedTaskWithLaunchPath:@"/usr/bin/open" arguments:[NSArray arrayWithObjects:logPath, nil]] waitUntilExit];
    
}

- (IBAction)updateGFWList:(NSMenuItem *)sender {
    [PACUtils UpdatePACFromGFWList:^{
        [self Toast:@"updated gfw file ok"];
    } fail:^{
        [self Toast:@"updated gfw file fail"];
    }];
}

- (IBAction)editUserRulesForPAC:(NSMenuItem *)sender {
    
    [_userRuleWindow showWindow:nil];
    [NSApp activateIgnoringOtherApps:YES];
    [_userRuleWindow.window makeKeyAndOrderFront:sender];
}

#pragma mark 设置界面UI
-(void)setBarStatus
{
    statusBarItem = [[NSStatusBar systemStatusBar] statusItemWithLength:24.0];
    statusBarItem.image = [NSImage imageNamed:@"dagger"];
    statusBarItem.alternateImage = [NSImage imageNamed:@"dagger"];
    statusBarItem.menu = statusBarItemMenu;
    statusBarItem.toolTip = @"dagger";
    [statusBarItem setHighlightMode:YES];
}

-(void)initConfig{
    NSUserDefaults *shared = [NSUserDefaults standardUserDefaults];
    
    
    [shared registerDefaults:@{
        @"launchAtLogin":@NO,
        @"DaggerOn":@NO,
        @"DaggerMode":@"auto",
        @"LocalSocks5.ListenPort": @"1096",
        @"LocalSocks5.ListenAddress": @"127.0.0.1",
        @"LocalHTTP.ListenAddress": @"127.0.0.1",
        @"LocalHTTP.ListenPort": @"1097",
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

-(void)initWindow
{
    
    NSArray *listVC = @[
        [[PreferencesGeneral alloc] init],
        [[PreferencesAdvanced alloc] init],
        [[PreferencesInterfaces alloc] init],
    ];
    
    _preferenceWindow = [[MASPreferencesWindowController alloc] initWithViewControllers:listVC title:@""];
    _preferenceWindow.window.level = NSFloatingWindowLevel;

    _serverConf = [[Servers alloc] init];
    _serverConf.window.level = NSFloatingWindowLevel;
    _userRuleWindow = [[UserRules alloc] init];
}

#pragma mark Servers
- (IBAction)setServers:(id)sender {
    [_serverConf showWindow:nil];
}

#pragma mark Preferences
- (IBAction)showPreferences:(id)sender {
    [_preferenceWindow showWindow:nil];
}

-(void)regNotifEvent{
    // drag event
    [[NSNotificationCenter defaultCenter] addObserver:self selector:@selector(changeConfigListEvent) name:@"changeConfigList" object:nil];
}

-(void)changeConfigListEvent{
    [ProxyConfHelper disableProxy];
    [LaunchAgentsUtils stopHttpProxy];
    
    [self updateServersMenu];
    [self applyConf];
    
    [self Toast:@"Update succeeded"];
}

- (void)applicationDidFinishLaunching:(NSNotification *)aNotification {
    [self initConfig];
    [self initWindow];
    [self setBarStatus];
    [self regNotifEvent];
    
    [ProxyConfHelper install];
    [PACUtils install];
    [LaunchAgentsUtils install];
    
    [PreferencesGeneral setLaunchAtLogin];
    
    [self updateMainMenu];
    [self updateRunningModeMenu];
    [self updateServersMenu];
    [self applyConf];
}


- (void)applicationWillTerminate:(NSNotification *)aNotification {
    [ProxyConfHelper disableProxy];
    [LaunchAgentsUtils stopHttpProxy];
}


- (BOOL)applicationSupportsSecureRestorableState:(NSApplication *)app {
    return YES;
}


@end
