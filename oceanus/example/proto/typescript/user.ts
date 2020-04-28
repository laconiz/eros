namespace proto {

    // 用户ID
    export type UserID = number;

    // 性别
    export enum Gender {
        Invalid = 0, // 未设置
        Male = 1, // 男
        Female = 2, // 女
    }

    // 用户
    export class User {
        ID : UserID; // 用户ID
        Name : string; // 名称
        Avatar : string; // 头像
        Gender : Gender; // 性别
    }

    // 用户信息消息
    export class ProfileACK {
        User : User; // 基础信息
        Phone : string; // 手机
    }

    // 修改用户信息消息
    export class ProfileModifyACK {
        Name : string; // 昵称
        Phone : string; // 手机
        Avatar : string; // 头像
    }

    Register("ProfileACK", ProfileACK);
    Register("ProfileModifyACK", ProfileModifyACK);
}