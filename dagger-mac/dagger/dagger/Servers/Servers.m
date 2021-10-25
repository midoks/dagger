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


-(void)reloadListData
{
    NSString *serverPath = [AppCommon getServerPlist];
    _list = [[NSMutableArray alloc] initWithContentsOfFile:serverPath];
    
    NSLog(@"serverPath:%@", serverPath);
    NSLog(@"ll:%@", [[NSMutableArray alloc] initWithContentsOfFile:serverPath]);
    [_tableView reloadData];
}

- (void)windowDidLoad {
    [super windowDidLoad];
    
    [self reloadListData];
    
}

#pragma mark tableView
- (NSInteger)numberOfRowsInTableView:(NSTableView *)tableView
{
    return [_list count];
}

-(NSView *)tableView:(NSTableView *)tableView viewForTableColumn:(NSTableColumn *)tableColumn row:(NSInteger)row
{
    ServersModel *hnm = [_list objectAtIndex:row];
    NSTableCellView *cell = [tableView makeViewWithIdentifier:tableColumn.identifier owner:self];
    [cell.textField setStringValue:hnm.remark];
    [cell.textField setEditable:NO];
    [cell.textField setDrawsBackground:NO];
    return cell;
}

- (NSDragOperation)tableView:(NSTableView *)tableView validateDrop:(id<NSDraggingInfo>)info proposedRow:(NSInteger)row proposedDropOperation:(NSTableViewDropOperation)dropOperation
{
    if (dropOperation == NSTableViewDropAbove) {
        return NSDragOperationNone;
    }
    return NSDragOperationMove;
}


#pragma mark 点击选择框
- (void)tableViewSelectionDidChange:(NSNotification *)notification
{

}
- (IBAction)add:(id)sender
{
    NSString *remark = [NSString stringWithFormat:@"web-%ld", [_list count]+1];
    
    
    [_domain setStringValue:@"wss://domain.xyz"];
    [_path setStringValue:@"ws"];
    [_remark setStringValue:remark];
    
    
    ServersModel *sm = [[ServersModel alloc] init];
    sm.domain = _domain.stringValue;
    sm.path = _path.stringValue;
    sm.remark = _remark.stringValue;
    
    [_list addObject:sm];
    [self reloadListData];
}


@end
