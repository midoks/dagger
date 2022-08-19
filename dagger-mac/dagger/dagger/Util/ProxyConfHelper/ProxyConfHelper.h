//
//  ProxyConfHelper.h
//  dagger
//
//  Created by midoks on 2021/10/25.
//

#import <Foundation/Foundation.h>
#import "GCDWebServer.h"
#import "GCDWebServerDataResponse.h"

NS_ASSUME_NONNULL_BEGIN

@interface ProxyConfHelper : NSObject
+ (void)install;
+ (void)enablePACProxy;
+ (void)enableGlobalProxy;
+ (void)disableProxy;
+ (void)enableExternalPACProxy;
+ (void)startMonitorPAC;


+ (void)getCFSpeedTest:(NSString *)domain callback:(void(^)(NSString *)) callback;
+ (void)setCfIpClean:(NSString *)domain callback:(void(^)(NSString *)) callback;
+ (void)setCfIpPreference:(NSString *)domain callback:(void(^)(NSString *)) callback;
@end

NS_ASSUME_NONNULL_END
