//
//  PreferencesGeneral.m
//  dagger
//
//  Created by midoks on 2021/10/24.
//

#import "PreferencesGeneral.h"

@interface PreferencesGeneral ()

@end

@implementation PreferencesGeneral

-(id)init{
    self = [self initWithNibName:@"PreferencesGeneral" bundle:nil];
    return self;
}

- (void)viewDidLoad {
    [super viewDidLoad];
    // Do view setup here.
}

#pragma mark - MASPreferencesViewController
- (NSString *)viewIdentifier
{
    return @"PreferencesGeneral";
}

- (NSImage *)toolbarItemImage
{
    return [NSImage imageNamed:NSImageNamePreferencesGeneral];
}

- (NSString *)toolbarItemLabel
{
    return @"General";
}

@end
