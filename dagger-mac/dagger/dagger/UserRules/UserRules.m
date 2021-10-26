//
//  UserRules.m
//  dagger
//
//  Created by midoks on 2021/10/26.
//

#import "UserRules.h"
#import "PACUtils.h"
#import "ProxyConfHelper.h"

@interface UserRules ()

@property (weak) IBOutlet NSTextView *userRulesView;

@end

@implementation UserRules

static UserRules *_instance = nil;
static dispatch_once_t _instance_once;
+ (id)Instance{
    dispatch_once(&_instance_once, ^{
        _instance = [[UserRules alloc] init];
    });
    return _instance;
}

-(id)init{
    self = [super initWithWindowNibName:@"UserRules"];
    return self;
}

-(void)applyConf{
    NSUserDefaults *shared = [NSUserDefaults standardUserDefaults];
    NSString *mode = [shared objectForKey:@"DaggerMode"];
    BOOL on = [shared boolForKey:@"DaggerOn"];
    
    if (on) {
        if ([mode isEqualTo:@"auto"]){
            [ProxyConfHelper enablePACProxy];
        } else if ([mode isEqualTo:@"global"]){
            [ProxyConfHelper enableGlobalProxy];
        } else if ([mode isEqualTo:@"manual"]){
            [ProxyConfHelper disableProxy];
        }
    } else {
        [ProxyConfHelper disableProxy];
    }
}

- (IBAction)btnOK:(id)sender
{
    NSString *pacDir = [NSString stringWithFormat:@"%@/%s", NSHomeDirectory(), PAC_DEFAULT_DIR];
    NSString *pacUserRuleDirPath = [NSString stringWithFormat:@"%@/%s",pacDir, PAC_USER_RULE_PATH];
    NSString *ur = [_userRulesView string];
    
    [ur writeToFile:pacUserRuleDirPath atomically:YES encoding:NSUTF8StringEncoding error:nil];
    [PACUtils GeneratePACFile];
    
    [ProxyConfHelper disableProxy];
    [self applyConf];
}

- (IBAction)btnCancel:(id)sender
{
    [self.window close];
}


- (void)windowDidLoad {
    [super windowDidLoad];
    NSString *pacDir = [NSString stringWithFormat:@"%@/%s", NSHomeDirectory(), PAC_DEFAULT_DIR];
    NSString *pacUserRuleDirPath = [NSString stringWithFormat:@"%@/%s",pacDir, PAC_USER_RULE_PATH];
    NSString *userRuleTextContent = [NSString stringWithContentsOfFile:pacUserRuleDirPath encoding:NSUTF8StringEncoding error:nil];
    
    [_userRulesView setString:userRuleTextContent];
}

@end
