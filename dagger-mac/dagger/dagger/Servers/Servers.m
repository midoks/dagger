//
//  Servers.m
//  dagger
//
//  Created by midoks on 2021/10/26.
//

#import "Servers.h"
#import "ServersModel.h"

@interface Servers () <NSTableViewDataSource, NSTableViewDelegate>

@property (weak) IBOutlet NSTableView *tableView;

@property (weak) IBOutlet NSTextField *remark;
@property (weak) IBOutlet NSTextField *password;
@property (weak) IBOutlet NSTextField *username;
@property (weak) IBOutlet NSTextField *path;
@property (weak) IBOutlet NSTextField *domain;

@property  NSMutableArray *list;
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
  
    _tableView.delegate = self;
    _tableView.dataSource = self;
    
    [self reloadListData];
    
}

-(IBAction)openConfigDir:(id)sender{
    
    NSString *dir = [[AppCommon appSupportDirURL] path];
    [[NSTask launchedTaskWithLaunchPath:@"/usr/bin/open" arguments:[NSArray arrayWithObjects:dir, nil]] waitUntilExit];
}

-(void)reloadListData
{
    NSString *serverPath = [AppCommon getServerPlist];
    _list = [[NSMutableArray alloc] initWithContentsOfFile:serverPath];
    
    NSLog(@"serverPath:%@", serverPath);
    [_tableView reloadData];
}

#pragma mark tableView
- (NSInteger)numberOfRowsInTableView:(NSTableView *)tableView
{
    return 3;
}

-(id)tableView:(NSTableView *)tableView objectValueForTableColumn:(NSTableColumn *)tableColumn row:(NSInteger)row{
    NSLog(@"dd");
    
    if ([tableColumn.identifier isEqualTo:@"main"]){
        
        return @"111";
    }
    
    if ([tableColumn.identifier isEqualTo:@"status"]){
        
        return [NSImage imageNamed:@"NSMenuOnStateTemplate"];
    }
    
    return [NSImage imageNamed:@"NSMenuOnStateTemplate"];
}




#pragma mark 点击选择框
- (void)tableViewSelectionDidChange:(NSNotification *)notification
{

}
- (IBAction)add:(id)sender
{
    NSString *remark = [NSString stringWithFormat:@"websocket-%ld", [_list count]+1];
    
    [_domain setStringValue:@"wss://domain.xyz"];
    [_path setStringValue:@"ws"];
    [_remark setStringValue:remark];
    
    ServersModel *sm = [[ServersModel alloc] init];
    sm.domain = _domain.stringValue;
    sm.path = _path.stringValue;
    sm.remark = _remark.stringValue;
    
    [_list addObject:sm];
    
    NSLog(@"%lu",(unsigned long)[_list count]);
    [_tableView reloadData];
    
//    [_tableView selectRowIndexes:[[NSIndexSet alloc] initWithIndex:[_list count]-1] byExtendingSelection:YES];
}


@end
