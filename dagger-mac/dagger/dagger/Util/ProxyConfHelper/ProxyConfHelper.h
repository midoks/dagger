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
@end

NS_ASSUME_NONNULL_END
