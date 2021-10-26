//
//  PACUtils.h
//  dagger
//
//  Created by midoks on 2021/10/25.
//

#import <Foundation/Foundation.h>

NS_ASSUME_NONNULL_BEGIN

#define PAC_DEFAULT_DIR ".dagger"
#define PAC_USER_RULE_PATH "user-rule.txt"
#define PAC_FILE_PATH "gfwlist.js"
#define PAC_GFW_FILE_PATH "gfwlist.txt"

@interface PACUtils : NSObject

+ (void)install;
+(void)UpdatePACFromGFWList;
@end

NS_ASSUME_NONNULL_END
