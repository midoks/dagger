//
//  AppDelegate.swift
//  dagger
//
//  Created by midoks on 2021/10/22.
//

import Cocoa
import Carbon

@NSApplicationMain
class AppDelegate: NSObject, NSApplicationDelegate,NSUserNotificationCenterDelegate {

    @IBOutlet var window: NSWindow!
    
    @IBOutlet var statusMenu: NSMenu!
    var statusItem: NSStatusItem!


    func applicationDidFinishLaunching(_ aNotification: Notification) {
        // Insert code here to initialize your application
        
        statusItem = NSStatusBar.system.statusItem(withLength: 11)
        let image : NSImage = NSImage(named: "dagger")!
        image.isTemplate = true
        statusItem.image = image
        statusItem.toolTip = "dagger";
        statusItem.menu = statusMenu
    }

    func applicationWillTerminate(_ aNotification: Notification) {
        // Insert code here to tear down your application
    }

    func applicationSupportsSecureRestorableState(_ app: NSApplication) -> Bool {
        return true
    }


}

