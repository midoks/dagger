//
//  Preferences.m
//  dagger
//
//  Created by midoks on 2021/10/24.
//

#import "Preferences.h"

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
    self = [self initWithWindowNibName:@"Preferences"];
    return self;
}

-(NSUserDefaults *)standardUserDefaults {
    return [NSUserDefaults standardUserDefaults];
}

-(void)sync{
    [[self standardUserDefaults] synchronize];
}


-(IBAction)toolbarAction:(NSToolbarItem *)sender {
    [_tabView selectTabViewItemWithIdentifier:[sender itemIdentifier]];
}


- (void)windowDidLoad {
    [super windowDidLoad];
    
    NSLog(@"Preferences");
}

@end
