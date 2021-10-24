//
//  PreferencesAdvanced.m
//  dagger
//
//  Created by midoks on 2021/10/24.
//

#import "PreferencesAdvanced.h"

@interface PreferencesAdvanced ()

@end

@implementation PreferencesAdvanced


-(id)init{
    self = [self initWithNibName:@"PreferencesAdvanced" bundle:nil];
    return self;
}

- (void)viewDidLoad {
    [super viewDidLoad];
}

#pragma mark - MASPreferencesViewController
- (NSString *)viewIdentifier
{
    return @"PreferencesAdvanced";
}

- (NSImage *)toolbarItemImage
{
    return [NSImage imageNamed:NSImageNameAdvanced];
}

- (NSString *)toolbarItemLabel
{
    return @"Advanced";
}

@end
