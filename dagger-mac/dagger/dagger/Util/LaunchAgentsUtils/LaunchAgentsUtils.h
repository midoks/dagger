//
//  LaunchAgentsUtils.h
//  dagger
//
//  Created by midoks on 2021/10/27.
//

#import <Foundation/Foundation.h>

NS_ASSUME_NONNULL_BEGIN

#define LAUNCH_AGENT_DIR @"Library/LaunchAgents/"
#define LAUNCH_AGENT_CONF_HTTP_NAME  @"com.midoks.dagger.http.plist"

#define APP_SUPPORT_DIR  @"Library/Application Support/dagger/"

@interface LaunchAgentsUtils : NSObject

+(void)install;
+(BOOL)generateHttpLauchAgentPlist;
+(void)startHttpProxy;
+(void)stopHttpProxy;
@end

NS_ASSUME_NONNULL_END
