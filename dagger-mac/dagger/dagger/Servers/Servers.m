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
    
    _list = [[NSMutableArray alloc] init];
    
    [self reloadListData];
    
}


- (IBAction)add:(id)sender
{
    NSString *pathplist = [AppCommon getServerPlist];
    NSString *remark = [NSString stringWithFormat:@"websocket-%ld", [_list count]+1];
    
    [_domain setStringValue:@"wss://domain.xyz"];
    [_path setStringValue:@"ws"];
    [_remark setStringValue:remark];
    
    NSMutableDictionary *info = [[NSMutableDictionary alloc] init];
    [info setObject:_remark.stringValue forKey:@"remark"];
    [info setObject:_domain.stringValue forKey:@"domain"];
    [info setObject:_path.stringValue forKey:@"path"];
    [info setObject:@"" forKey:@"username"];
    [info setObject:@"" forKey:@"password"];
    [info setObject:@"off" forKey:@"status"];
    
    [_list addObject:info];
    
    [_tableView reloadData];
    [_list writeToFile:pathplist atomically:YES];
}

- (IBAction)remove:(id)sender
{
    if ([_list count]==1){
        return;
    }
    NSString *pathplist = [AppCommon getServerPlist];
    NSInteger row = [_tableView selectedRow];
    [_list removeObjectAtIndex:row];
    [_tableView reloadData];
    [_list writeToFile:pathplist atomically:YES];
}

-(IBAction)openConfigDir:(id)sender{
    
    NSString *dir = [[AppCommon appSupportDirURL] path];
    [[NSTask launchedTaskWithLaunchPath:@"/usr/bin/open" arguments:[NSArray arrayWithObjects:dir, nil]] waitUntilExit];
}

- (IBAction)btnOK:(id)sender
{
    NSString *pathplist = [AppCommon getServerPlist];
    NSInteger row = [_tableView selectedRow];
    
    [[_list objectAtIndex:row] setObject:_domain.stringValue forKey:@"domain"];
    [[_list objectAtIndex:row] setObject:_path.stringValue forKey:@"path"];
    [[_list objectAtIndex:row] setObject:_remark.stringValue forKey:@"remark"];
    [[_list objectAtIndex:row] setObject:_username.stringValue forKey:@"username"];
    [[_list objectAtIndex:row] setObject:_password.stringValue forKey:@"password"];
    
    [_list writeToFile:pathplist atomically:YES];
    [_tableView reloadData];
}

- (IBAction)btnCancel:(id)sender
{
    [self.window close];
}


-(void)reloadListData
{
    NSString *serverPath = [AppCommon getServerPlist];
    NSFileManager *fm = [NSFileManager defaultManager];
    
    if([fm fileExistsAtPath:serverPath]){
        _list = [[NSMutableArray alloc] initWithContentsOfFile:serverPath];
    }
    
    [_tableView reloadData];
}

#pragma mark tableView
- (NSInteger)numberOfRowsInTableView:(NSTableView *)tableView
{
    return [_list count];
}

-(id)tableView:(NSTableView *)tableView objectValueForTableColumn:(NSTableColumn *)tableColumn row:(NSInteger)row{
    
    NSMutableDictionary *sm = [_list objectAtIndex:row];
    
    
    if ([tableColumn.identifier isEqualTo:@"main"]){
        return [sm valueForKey:@"remark"];
    }
    
    if (([tableColumn.identifier isEqualTo:@"status"]) && ([[sm objectForKey:@"status"] isEqualTo:@"on"])) {
        return [NSImage imageNamed:@"NSMenuOnStateTemplate"];
    }
    
    return [NSImage imageNamed:@""];
}

#pragma mark 点击选择框
- (void)tableViewSelectionDidChange:(NSNotification *)notification
{
    
    NSInteger row = [_tableView selectedRow];
    NSMutableDictionary *sm = [_list objectAtIndex:row];
    
    [_domain setStringValue:[sm objectForKey:@"domain"]];
    [_path setStringValue:[sm objectForKey:@"path"]];
    [_remark setStringValue:[sm objectForKey:@"remark"]];
    [_username setStringValue:[sm objectForKey:@"username"]];
    [_password setStringValue:[sm objectForKey:@"password"]];
}


@end
