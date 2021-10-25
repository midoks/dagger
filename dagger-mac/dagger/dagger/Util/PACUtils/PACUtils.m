//
//  PACUtils.m
//  dagger
//
//  Created by midoks on 2021/10/25.
//

#import "PACUtils.h"

@implementation PACUtils


+ (void)install {
    BOOL needGenerate = NO;
    NSUserDefaults *shared = [NSUserDefaults standardUserDefaults];
    
    NSString *nowSocks5Address = [shared objectForKey:@"LocalSocks5.ListenAddress"];
    NSString *oldSocks5Address = [shared objectForKey:@"LocalSocks5.ListenAddress.old"];
    
    if ([nowSocks5Address isNotEqualTo:oldSocks5Address]){
        needGenerate = YES;
        [shared setObject:nowSocks5Address forKey:@"LocalSocks5.ListenAddress.old"];
    }
    
    NSString *nowSocks5Port = [shared objectForKey:@"LocalSocks5.ListenPort"];
    NSString *oldSocks5Port = [shared objectForKey:@"LocalSocks5.ListenPort.old"];
    
    if ([nowSocks5Port isNotEqualTo:oldSocks5Port]){
        needGenerate = YES;
        [shared setObject:nowSocks5Address forKey:@"LocalSocks5.ListenPort.old"];
    }
    
    NSFileManager *fm = [NSFileManager defaultManager];
    NSString *pacDir = [NSString stringWithFormat:@"%@/%s", NSHomeDirectory(), PAC_DEFAULT_DIR];
    
    if (![fm fileExistsAtPath:pacDir]){
        needGenerate = YES;
    }
    
    if (needGenerate){NSLog(@"pac file install ...");
        [self GeneratePACFile];
        
    }
}

+(void)GeneratePACFile{
    NSUserDefaults *shared = [NSUserDefaults standardUserDefaults];
    NSFileManager *fm = [NSFileManager defaultManager];
    
    NSString *socks5Address = [shared objectForKey:@"LocalSocks5.ListenAddress"];
    NSString *socks5Port = [shared objectForKey:@"LocalSocks5.ListenPort"];
    
    NSString *pacDir = [NSString stringWithFormat:@"%@/%s", NSHomeDirectory(), PAC_DEFAULT_DIR];
    NSString *pacUserRuleDirPath = [NSString stringWithFormat:@"%@/%s",pacDir, PAC_USER_RULE_PATH];
    NSString *pacGFWDirPath = [NSString stringWithFormat:@"%@/%s",pacDir, PAC_GFW_FILE_PATH];
    NSString *pacGFWJSPath = [NSString stringWithFormat:@"%@/%s",pacDir, PAC_FILE_PATH];
 
    
    
    if (![fm fileExistsAtPath:pacDir]){
        [fm createDirectoryAtURL:[NSURL fileURLWithPath:pacDir] withIntermediateDirectories:YES attributes:nil error:nil];
    }
    
    if (![fm fileExistsAtPath:pacGFWDirPath]){
        NSString *src = [[NSBundle mainBundle] pathForResource:@"gfwlist" ofType:@"txt"];
        [fm copyItemAtPath:src toPath:pacGFWDirPath error:nil];
    }
    
    if (![fm fileExistsAtPath:pacUserRuleDirPath]){
        NSString *src = [[NSBundle mainBundle] pathForResource:@"user-rule" ofType:@"txt"];
        [fm copyItemAtPath:src toPath:pacUserRuleDirPath error:nil];
    }

    NSString *gfwBase64String = [NSString stringWithContentsOfFile:pacGFWDirPath encoding:NSUTF8StringEncoding error:nil];
    NSData *data = [[NSData alloc]initWithBase64EncodedString:gfwBase64String options:NSDataBase64DecodingIgnoreUnknownCharacters];
    NSString *gfw = [[NSString alloc]initWithData:data encoding:NSUTF8StringEncoding];
    NSArray *gfwLine = [gfw componentsSeparatedByString:@"\n"];
    
    
    gfwLine = [gfwLine filteredArrayUsingPredicate:[NSPredicate predicateWithBlock:^BOOL(NSString* object, NSDictionary *bindings) {
        
        if ([object isEqualToString:@""]){
            return NO;
        }
        
        if ([[object substringToIndex:1] isEqualTo:@"!"] || [[object substringToIndex:1] isEqualTo:@"["]){
            return NO;
        }
        return  YES;
    }]];
    
    
   
    
    NSString *userContent = [NSString stringWithContentsOfFile:pacUserRuleDirPath encoding:NSUTF8StringEncoding error:nil];
    NSArray *userLine = [userContent componentsSeparatedByString:@"\n"];
    userLine = [userLine filteredArrayUsingPredicate:[NSPredicate predicateWithBlock:^BOOL(NSString* object, NSDictionary *bindings) {
        if ([object isEqualToString:@""]){
            return NO;
        }
        if ([[object substringToIndex:1] isEqualTo:@"!"] || [[object substringToIndex:1] isEqualTo:@"["]){
            return NO;
        }
        return  YES;
    }]];
    
    
    gfwLine = [gfwLine arrayByAddingObjectsFromArray:userLine];
    NSData *jsonData = [NSJSONSerialization dataWithJSONObject:gfwLine options:NSJSONWritingPrettyPrinted error:nil
    ];
    NSString *gfwJsonStr=[[NSString alloc]initWithData:jsonData encoding:NSUTF8StringEncoding
    ];

    NSString *jsPath = [[NSBundle mainBundle] pathForResource:@"abp" ofType:@"js"];
    NSString *jsContent = [NSString stringWithContentsOfFile:jsPath encoding:NSUTF8StringEncoding error:nil];
    
    jsContent =  [jsContent stringByReplacingOccurrencesOfString:@"__RULES__" withString:gfwJsonStr];
    jsContent = [jsContent stringByReplacingOccurrencesOfString:@"__SOCKS5PORT__" withString:socks5Port];
    jsContent = [jsContent stringByReplacingOccurrencesOfString:@"__SOCKS5ADDR__" withString:socks5Address];
    [jsContent writeToFile:pacGFWJSPath atomically:YES encoding:NSUTF8StringEncoding error:nil];
}


@end
