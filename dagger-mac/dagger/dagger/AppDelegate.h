//
//  AppDelegate.h
//  dagger
//
//  Created by midoks on 2021/10/24.
//

#import <Cocoa/Cocoa.h>

@interface AppDelegate : NSObject <NSApplicationDelegate>{
    

    NSStatusItem *statusBarItem;
    __weak IBOutlet NSMenu *statusBarItemMenu;
    
}


@end

