//
//  AppCommon.m
//  dagger
//
//  Created by midoks on 2021/10/26.
//

#import "AppCommon.h"

@implementation AppCommon

#pragma mark 获取上一级目录
+(NSString *)getDirName:(NSString *)dirname
{
    NSArray *i = [dirname componentsSeparatedByString:@"/"];
    NSMutableArray *ii = [[NSMutableArray alloc] initWithArray:i];
    [ii removeLastObject];
    NSString *r = [ii componentsJoinedByString:@"/"];
    return r;
}

#pragma mark 获取运行根目录
+ (NSString *)getRootDir
{
    char path[1024];
    unsigned size = 1024;
    
    _NSGetExecutablePath(path, &size);
    path[size] = '\0';
    
    NSString *str = [NSString stringWithFormat:@"%s", path];
    str = [self getDirName:str];
    str = [self getDirName:str];
    str = [self getDirName:str];
    str = [self getDirName:str];
    
//    str = [NSString stringWithFormat:@"%@/mdserver/", str];
    str = [NSString stringWithFormat:@"/Applications/dagger/%@", @""];
    return str;
}

#pragma mark - 创建目录
+(void)createDirIfNoExist:(NSURL *)url{
    NSFileManager *fm = [NSFileManager defaultManager];
    NSString *path = [url path];
    if (![fm fileExistsAtPath:path]){
        [fm createDirectoryAtURL:url withIntermediateDirectories:YES attributes:nil error:nil];
    }
}

#pragma mark - 获取App支持的目录,不存在就自动创建
+(NSURL *)appSupportDirURL {
    NSFileManager *fm = [NSFileManager defaultManager];
    NSArray<NSURL *> *asPath = [fm URLsForDirectory:NSApplicationSupportDirectory inDomains:NSUserDomainMask];
    NSString *bundleID = [NSBundle mainBundle].bundleIdentifier;
    NSURL * appAsUrl = [asPath.firstObject URLByAppendingPathComponent:bundleID];
    
    [self createDirIfNoExist:appAsUrl];
    return appAsUrl;
}


#pragma mark 如果你希望调用系统命
+ (void)runSystemCommand:(NSString *)cmd
{
    [[NSTask launchedTaskWithLaunchPath:@"/bin/sh"
                              arguments:[NSArray arrayWithObjects:@"-c", cmd, nil]]
     waitUntilExit];
}

+(NSString*)getServerPlist {
    NSFileManager *fm = [NSFileManager defaultManager];
    NSURL *dirUrl = [AppCommon appSupportDirURL];
    NSURL *writeURL = [dirUrl URLByAppendingPathComponent:@"server.plist"];
    
    if (![fm fileExistsAtPath:[writeURL path]]){
        NSString *serverList = [[NSBundle mainBundle] pathForResource:@"server" ofType:@"plist"];
        NSString *content = [NSString stringWithContentsOfFile:serverList encoding:NSUTF8StringEncoding error:nil];
        [content writeToURL:writeURL atomically:YES encoding:NSUTF8StringEncoding error:nil];
    }
    return  [writeURL path];
}
@end
