//
//  Servers.m
//  dagger
//
//  Created by midoks on 2021/10/26.
//

#import "Servers.h"

@interface Servers ()

@end


@implementation Servers

static Servers *_instance = nil;
static dispatch_once_t _instance_once;
+ (id)Instance{
    dispatch_once(&_instance_once, ^{
        _instance = [[Servers alloc] init];
    });
    return _instance;
}

-(id)init{
    self = [super initWithWindowNibName:@"Servers"];
    return self;
}

- (void)windowDidLoad {
    [super windowDidLoad];
    
    // Implement this method to handle any initialization after your window controller's window has been loaded from its nib file.
}

@end
