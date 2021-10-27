//
//  PreferencesGeneral.m
//  dagger
//
//  Created by midoks on 2021/10/24.
//

#import "PreferencesGeneral.h"

@interface PreferencesGeneral ()

@end

@implementation PreferencesGeneral

-(id)init{
    self = [self initWithNibName:@"PreferencesGeneral" bundle:nil];
    return self;
}

- (void)viewDidLoad {
    [super viewDidLoad];
}

- (IBAction)setLaunchAtLoginAction:(NSButton *)sender {
    [PreferencesGeneral setLaunchAtLogin];
}

#pragma mark - 设置开机自启动
+(void)setLaunchAtLogin {
    NSUserDefaults *shared = [NSUserDefaults standardUserDefaults];
    BOOL launchAtLoginEnable = [shared boolForKey:@"launchAtLogin"];
    
//    NSLog(@"launchAtLoginEnable:%hhd",launchAtLoginEnable);
    
    NSString* launchFolder = [NSString stringWithFormat:@"%@/Library/LaunchAgents",NSHomeDirectory()];
    NSString * bundleID = [[NSBundle mainBundle] objectForInfoDictionaryKey:(NSString *)kCFBundleIdentifierKey];

    NSString* dstLaunchPath = [launchFolder stringByAppendingFormat:@"/%@.plist",bundleID];

    
    NSFileManager* fm = [NSFileManager defaultManager];
    BOOL isDir = NO;
    
    //已经存在启动项中，就不必再创建
//    if ([fm fileExistsAtPath:dstLaunchPath isDirectory:&isDir] && !isDir) {
//       return;
//    }
    
    //下面是一些配置
    NSMutableDictionary* dict = [[NSMutableDictionary alloc] init];
    NSMutableArray* arr = [[NSMutableArray alloc] init];
    [arr addObject:[[NSBundle mainBundle] executablePath]];
    [arr addObject:@"-runMode"];
    [arr addObject:@"autoLaunched"];
    
    if (launchAtLoginEnable){
        [dict setObject:[NSNumber numberWithBool:true] forKey:@"RunAtLoad"];
    } else{
        [dict setObject:[NSNumber numberWithBool:false] forKey:@"RunAtLoad"];
    }
    
    [dict setObject:bundleID forKey:@"Label"];
    [dict setObject:arr forKey:@"ProgramArguments"];
    isDir = NO;
    if (![fm fileExistsAtPath:launchFolder isDirectory:&isDir] && isDir) {
       [fm createDirectoryAtPath:launchFolder withIntermediateDirectories:NO attributes:nil error:nil];
    }

//    NSLog(@"dict : %@",dict);
    [dict writeToFile:dstLaunchPath atomically:NO];
}


#pragma mark - MASPreferencesViewController
- (NSString *)viewIdentifier
{
    return @"PreferencesGeneral";
}

- (NSImage *)toolbarItemImage
{
    return [NSImage imageNamed:NSImageNamePreferencesGeneral];
}

- (NSString *)toolbarItemLabel
{
    return @"General";
}

@end
