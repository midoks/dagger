//
//  LaunchAgentsUtils.m
//  dagger
//
//  Created by midoks on 2021/10/27.
//

#import "LaunchAgentsUtils.h"
#import "AppCommon.h"

@implementation LaunchAgentsUtils

+(void)install{
    
    NSFileManager *fm = [NSFileManager defaultManager];
    NSString *homeDir = NSHomeDirectory();
    NSString *appSupportDir = [NSString  stringWithFormat:@"%@/%@", homeDir, APP_SUPPORT_DIR];
    NSString *httpProxy = [NSString  stringWithFormat:@"%@/%@", appSupportDir, @"dagger-client-http"];
    if ((![fm fileExistsAtPath:appSupportDir]) || (![fm fileExistsAtPath:httpProxy])) {
        NSString *sh = [NSString stringWithFormat:@"%@/%@", [[NSBundle mainBundle] resourcePath], @"install_dagger_proxy.sh"];
        NSLog(@"run install [%@] script: %@", @"dagger-client-http",sh);
        [AppCommon runSystemCommand:sh];
        NSLog(@"installation [%@] success", @"dagger-client-http");
    }
}

@end
