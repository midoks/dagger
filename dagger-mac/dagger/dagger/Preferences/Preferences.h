//
//  Preferences.h
//  dagger
//
//  Created by midoks on 2021/10/24.
//

#import <Cocoa/Cocoa.h>

NS_ASSUME_NONNULL_BEGIN

@interface Preferences : NSWindowController


@property (weak) IBOutlet NSToolbar *toolbar;
@property (weak) IBOutlet NSTabView *tabView;


+ (id)Instance;
@end

NS_ASSUME_NONNULL_END
