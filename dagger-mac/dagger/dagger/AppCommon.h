//
//  AppCommon.h
//  dagger
//
//  Created by midoks on 2021/10/26.
//

#import <Foundation/Foundation.h>
#import <Cocoa/Cocoa.h>
#include <mach-o/dyld.h>

NS_ASSUME_NONNULL_BEGIN

@interface AppCommon : NSObject



#pragma mark - 创建目录
+(void)createDirIfNoExist:(NSURL *)url;
#pragma mark - 获取App支持的目录,不存在就自动创建
+(NSURL *)appSupportDirURL;

+(NSString*)getServerPlist;

+ (void)runSystemCommand:(NSString *)cmd;

@end

NS_ASSUME_NONNULL_END
