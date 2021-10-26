//
//  ServersModel.h
//  dagger
//
//  Created by midoks on 2021/10/26.
//

#import <Foundation/Foundation.h>

NS_ASSUME_NONNULL_BEGIN

@interface ServersModel : NSObject


@property (nonatomic, strong) NSString *remark;
@property (nonatomic, strong) NSString *domain;
@property (nonatomic, strong) NSString *path;
@property (nonatomic, strong) NSString *username;
@property (nonatomic, strong) NSString *password;


-(NSMutableDictionary *)setWithValue:(NSString *)remark
                             domain:(NSString *)domain
                               path:(NSString *)path
                           username:(NSString *)username
                            password:(NSString *)password;
@end

NS_ASSUME_NONNULL_END
