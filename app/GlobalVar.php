<?php

namespace App;

class GlobalVar
{
    const LACK_AUTH_ERR_RESPONSE = [
        'code' => -1,
        'msg' => '认证信息缺失',
    ];
    const INVALID_AUTH_ERR_RESPONSE = [
        'code' => -2,
        'msg' => '认证信息无效',
    ];
    const EXPIRED_AUTH_ERR_RESPONSE = [
        'code' => -3,
        'msg' => '认证信息过期',
    ];
    const LACK_REFRESHTOKEN_ERR_RESPONSE = [
        'code' => -4,
        'msg' => '本地登录信息缺失',
    ];
    const INVALID_REFRESHTOKEN_ERR_RESPONSE = [
        'code' => -5,
        'msg' => '本地登录信息无效',
    ];
    const EXPIRED_REFRESHTOKEN_ERR_RESPONSE = [
        'code' => -6,
        'msg' => '本地登录信息过期',
    ];
    const BLACKLIST_REFRESHTOKEN_ERR_RESPONSE = [
        'code' => -7,
        'msg' => '本地登录信息封禁',
    ];

    /**
     * 常规token过期时间
     *
     * 30分钟
     */
    const ACCESS_TOKEN_EXP_TIME = 60 * 30;
    /**
     * RefreshToken过期时间
     *
     * 60天
     */
    const REFRESH_TOKEN_EXP_TIME = 60 * 60 * 24 * 60;
    /**
     * RefreshToken重签判定时间
     */
    const REFRESH_TOKEN_RE_ISSUE_TIME = 60 * 60 * 24 * 60 / 2;
}
