//
//  LaunchAgentsUtils.m
//  dagger
//
//  Created by midoks on 2021/10/27.
//

#import "LaunchAgentsUtils.h"
#import "AppCommon.h"
#import "Servers.h"

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
    
    NSUserDefaults *shared = [NSUserDefaults standardUserDefaults];
    NSString *localHttpListenPort = [shared objectForKey:@"LocalHTTP.ListenPort"];
    
    NSString *homeDir = NSHomeDirectory();
    NSString *logFilePath = [NSString  stringWithFormat:@"%@/%@", homeDir, @"Library/Logs/dagger-client-http.log"];
    NSString *appSupportDir = [NSString  stringWithFormat:@"%@/%@", homeDir, APP_SUPPORT_DIR];
    NSString *daggerHttp =[NSString  stringWithFormat:@"%@%@", appSupportDir,@"dagger-client-http"];
    NSString *launchAgentDirPath = [NSString  stringWithFormat:@"%@/%@", homeDir, LAUNCH_AGENT_DIR];
    NSString *plistFilepath = [NSString  stringWithFormat:@"%@/%@", launchAgentDirPath, LAUNCH_AGENT_CONF_HTTP_NAME];
    
    
    
    NSMutableArray *arguments = [@[daggerHttp, @"service"]mutableCopy];
    
    NSMutableArray *list = [Servers serverList];
    NSDictionary *dst = nil;
    for (NSDictionary *i in list){
        if([[i objectForKey:@"status"] isEqualTo:@"on"]){
            dst = i;
            break;
        }
    }
    
    [arguments addObject:@"-p"];
    [arguments addObject:[NSString stringWithFormat:@"localhost:%@",localHttpListenPort]];
    
    if (dst){
        NSString *domain = [dst objectForKey:@"domain"];
        NSString *path = [dst objectForKey:@"path"];
        NSString *wsLink = [NSString stringWithFormat:@"%@/%@",domain, path];
        
        [arguments addObject:@"-w"];
        [arguments addObject:wsLink];
        
        
        NSString *username = [dst objectForKey:@"username"];
        if ([username isNotEqualTo:@""]){
            [arguments addObject:@"-u"];
            [arguments addObject:username];
        }
        
        NSString *password = [dst objectForKey:@"password"];
        if ([username isNotEqualTo:@""]){
            [arguments addObject:@"-m"];
            [arguments addObject:password];
        }
    }
    
    
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
