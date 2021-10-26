//
//  PreferencesInterfaces.m
//  dagger
//
//  Created by midoks on 2021/10/24.
//

#import "PreferencesInterfaces.h"
#import "ProxyConfTool.h"

@interface PreferencesInterfaces ()<NSTableViewDelegate,NSTableViewDataSource>

@property (strong) NSArray *netList;
@property (weak) IBOutlet NSTableView *tableView;

@end

@implementation PreferencesInterfaces

-(id)init{
    self = [self initWithNibName:@"PreferencesInterfaces" bundle:nil];
    return self;
}

- (void)viewDidLoad {
    [super viewDidLoad];
    
//    NSUserDefaults *shared = [NSUserDefaults standardUserDefaults];
//    NSString *services = [shared objectForKey:@"Proxy4NetworkServices"];
    
    
    _netList = [ProxyConfTool networkServicesList];
//    NSLog(@"%@",_netList);
    
    _tableView.delegate = self;
    _tableView.dataSource = self;
    
    
    [_tableView reloadData];
}

#pragma mark - NSTableViewDelegate,NSTableViewDataSource

- (NSInteger)numberOfRowsInTableView:(NSTableView *)tableView
{
    return [_netList count];
}

-(id)tableView:(NSTableView *)tableView objectValueForTableColumn:(NSTableColumn *)tableColumn row:(NSInteger)row{
    NSButtonCell *cell = tableColumn.dataCell;
    cell.state = NSControlStateValueOn;
    NSString *tit =  [[_netList objectAtIndex:row] objectForKey:@"userDefinedName"];
    cell.title = tit;
    return cell;
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
