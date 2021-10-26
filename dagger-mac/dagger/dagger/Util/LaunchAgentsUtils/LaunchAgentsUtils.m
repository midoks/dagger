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
    
    [self generateHttpLauchAgentPlist];
}


+(void)startHttpProxy{
    
    NSString *sh = [NSString stringWithFormat:@"%@/%@", [[NSBundle mainBundle] resourcePath], @"start_dagger_proxy.sh"];
    NSLog(@"start [%@]",sh);
    [AppCommon runSystemCommand:sh];
    NSLog(@"start [%@] end", @"start_dagger_proxy.sh");
}

+(void)stopHttpProxy{
    
    NSString *sh = [NSString stringWithFormat:@"%@/%@", [[NSBundle mainBundle] resourcePath], @"stop_dagger_proxy.sh"];
    NSLog(@"stop [%@]",sh);
    [AppCommon runSystemCommand:sh];
    NSLog(@"stop [%@] end", @"start_dagger_proxy.sh");
}

+(BOOL)generateHttpLauchAgentPlist{
    
    NSString *homeDir = NSHomeDirectory();
    NSString *logFilePath = [NSString  stringWithFormat:@"%@/%@", homeDir, @"Library/Logs/dagger-client-http.log"];
    NSString *appSupportDir = [NSString  stringWithFormat:@"%@/%@", homeDir, APP_SUPPORT_DIR];
    NSString *daggerHttp =[NSString  stringWithFormat:@"%@%@", appSupportDir,@"dagger-client-http"];
    NSString *launchAgentDirPath = [NSString  stringWithFormat:@"%@/%@", homeDir, LAUNCH_AGENT_DIR];
    NSString *plistFilepath = [NSString  stringWithFormat:@"%@/%@", launchAgentDirPath, LAUNCH_AGENT_CONF_HTTP_NAME];
    
    NSArray *arguments = @[daggerHttp, @"service", @"-p", @"localhost:1098"];
    
    NSDictionary *info = @{
        @"Label":@"com.midoks.dagger.http",
        @"WorkingDirectory":appSupportDir,
        @"StandardOutPath": logFilePath,
        @"StandardErrorPath": logFilePath,
        @"ProgramArguments": arguments,
    };
    
    [info writeToFile:plistFilepath atomically:YES];
    
    return YES;
}

@end
