//
//  Preferences.h
//  dagger
//
//  Created by midoks on 2021/10/24.
//

#import <Cocoa/Cocoa.h>

#import "PreferencesGeneral.h"
#import "PreferencesAdvanced.h"
#import "PreferencesInterfaces.h"

NS_ASSUME_NONNULL_BEGIN

@interface Preferences : NSObject

+ (id)Instance;

#pragma mark - Public Methods
-(NSString *)subAlignXToString:(NSInteger)index;
-(NSString *)subAlignYToString:(NSInteger)index;
-(NSString *)hardwareDecoderOptionToString:(NSInteger)index;
-(NSString *)rtspTransportationOptionToString:(NSInteger)index;
-(NSString *)screenshotFormatToString:(NSInteger)index;

#pragma mark - SET
-(void)setString:(NSString *)value key:(NSString*)key;
-(void)setBool:(BOOL)value key:(NSString*)key;
-(void)setInteger:(NSInteger)value key:(NSString *)key;

#pragma mark - GET
-(BOOL)boolForKey:(NSString *)key;
-(float)floatForKey:(NSString *)key;
-(NSInteger)integerForKey:(NSString *)key;
-(NSString *)stringForKey:(NSString *)key;
-(NSString *)colorForKey:(NSString *)key;

@end

NS_ASSUME_NONNULL_END
