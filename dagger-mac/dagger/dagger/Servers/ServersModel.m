//
//  ServersModel.m
//  dagger
//
//  Created by midoks on 2021/10/26.
//

#import "ServersModel.h"

@implementation ServersModel


-(NSMutableDictionary *)setWithValue:(NSString *)remark
                             domain:(NSString *)domain
                               path:(NSString *)path
                           username:(NSString *)username
                           password:(NSString *)password
{
    NSMutableDictionary *serverinfo = [[NSMutableDictionary alloc] init];
    [serverinfo setObject:remark forKey:@"remark"];
    [serverinfo setObject:domain forKey:@"domain"];
    [serverinfo setObject:path forKey:@"path"];
    [serverinfo setObject:username forKey:@"username"];
    [serverinfo setObject:password forKey:@"password"];
    return serverinfo;
}

@end
