//
//  PreferencesInterfaces.m
//  dagger
//
//  Created by midoks on 2021/10/24.
//

#import "PreferencesInterfaces.h"

@interface PreferencesInterfaces ()

@end

@implementation PreferencesInterfaces

-(id)init{
    self = [self initWithNibName:@"PreferencesInterfaces" bundle:nil];
    return self;
}

- (void)viewDidLoad {
    [super viewDidLoad];
}

#pragma mark - MASPreferencesViewController
- (NSString *)viewIdentifier
{
    return @"PreferencesInterfaces";
}

- (NSImage *)toolbarItemImage
{
    return [NSImage imageNamed:NSImageNameNetwork];
}

- (NSString *)toolbarItemLabel
{
    return @"Interfaces";
}


@end
