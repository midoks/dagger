//
//  Preferences.m
//  dagger
//
//  Created by midoks on 2021/10/24.
//

#import "Preferences.h"

#import "PreferencesGeneral.h"
#import "PreferencesAdvanced.h"
#import "PreferencesInterfaces.h"

#import "MASPreferences.h"

@interface Preferences ()

@end

@implementation Preferences

static Preferences *_instance = nil;
static dispatch_once_t _instance_once;
+ (id)Instance{
    dispatch_once(&_instance_once, ^{
        _instance = [[Preferences alloc] init];
    });
    return _instance;
}

-(id)init{
    self = [super init];
    return self;
}

-(NSUserDefaults *)standardUserDefaults {
    return [NSUserDefaults standardUserDefaults];
}

-(void)sync{
    [[self standardUserDefaults] synchronize];
}

#pragma mark - Public Methods
-(NSString *)subAlignXToString:(NSInteger)index{
    NSArray *list = @[@"left",@"center",@"right"];
    return [list objectAtIndex:index];
}

-(NSString *)subAlignYToString:(NSInteger)index{
    NSArray *list = @[@"top",@"center",@"bottom"];
    return [list objectAtIndex:index];
}

-(NSString *)hardwareDecoderOptionToString:(NSInteger)index{
    NSArray *list = @[@"no",@"auto",@"auto-copy"];
    return [list objectAtIndex:index];
}

-(NSString *)rtspTransportationOptionToString:(NSInteger)index{
    NSArray *list = @[@"lavf",@"tcp",@"udp",@"http"];
    return [list objectAtIndex:index];
}

-(NSString *)screenshotFormatToString:(NSInteger)index{
    NSArray *list = @[@"png",@"jpg",@"jpeg",@"ppm",@"pgm",@"pgmyuv",@"tga"];
    return [list objectAtIndex:index];
}


#pragma mark - SET
-(void)setString:(NSString *)value key:(NSString*)key{
    [[self standardUserDefaults] setValue:value forKey:key];
}

-(void)setBool:(BOOL)value key:(NSString*)key{
    [[self standardUserDefaults] setBool:value forKey:key];
}

-(void)setInteger:(NSInteger)value key:(NSString *)key{
    [[self standardUserDefaults] setInteger:value forKey:key];
}


#pragma mark - GET
-(BOOL)boolForKey:(NSString *)key{
    return [[self standardUserDefaults] boolForKey:key];
}

-(float)floatForKey:(NSString *)key{
    return [[self standardUserDefaults] floatForKey:key];
}

-(NSString *)stringForKey:(NSString *)key{
    return [[self standardUserDefaults] stringForKey:key];
}

-(NSInteger)integerForKey:(NSString *)key{
    return [[self standardUserDefaults] integerForKey:key];
}

-(NSString *)colorForKey:(NSString *)key{
    NSData *data = [[self standardUserDefaults] dataForKey:key];
    NSColor *color = [NSKeyedUnarchiver unarchiveObjectWithData:data];
    if (color){
        color = [color colorUsingColorSpace:[NSColorSpace deviceRGBColorSpace]];
        NSString *result = [NSString stringWithFormat:@"%.1f/%.1f/%.1f/%.1f", color.redComponent, color.greenComponent,color.blueComponent,color.alphaComponent];
        return result;
    }
    return @"";
}

@end
