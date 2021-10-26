//
//  Servers.h
//  dagger
//
//  Created by midoks on 2021/10/26.
//

#import <Cocoa/Cocoa.h>
#import "AppCommon.h"

NS_ASSUME_NONNULL_BEGIN

@interface Servers : NSWindowController

+ (id)Instance;
+(NSMutableArray *)serverList;
+(void)set:(NSInteger )index value:(NSString *)value forKey:(NSString *)key;
@end

NS_ASSUME_NONNULL_END
