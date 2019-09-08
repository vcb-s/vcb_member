<?php

namespace App\Http\Controllers;

use App\GlobalVar;

use Illuminate\Support\Facades\Hash;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Collection;
use Illuminate\Http\Request;

class UserController extends Controller
{

    public function show($id)
    {
        return response()->json([
            'id' => $id,
            'hello' => 'world'
        ]);
    }

    /** 获取列表 */
    public function list(Request $request)
    {
        $current = $request->input('current', 1);
        $pageSize = $request->input('pageSize', 20);
        $group = (int)$request->input('group');

        $allUser = DB::table('user');

        if ($group) {
            $allUser->where('user.group', 'like', '%' . (string)$group . '%');
        }

        $allUser
            ->select(
                'user.id',
                'user.retired',
                'user.avast',
                'user.bio',
                'user.nickname',
                'user.job',
                'user.order',
                'user.group',
                // DB::raw('DATE_FORMAT(user.create_at, "%Y-%m-%d") as joinAt'),
            );

        $total = $allUser->count();

        $result = $allUser
            ->offset($pageSize * ($current - 1))
            ->limit($pageSize)
            ->get();

        $result->each(function ($item, $index) {
            $item->group = array_map(function ($id) {
                return (int)$id;
            }, explode(',', $item->group));

            $item->order = (int)$item->order || 0;
        });


        return response()->json([
            'result' => $result,
            'pagination' => [
                'current' => $current,
                'pageSize' => $pageSize,
                'total' => $total
            ]
        ]);
    }

    /** 登录 */
    public function login(Request $request)
    {
        $name = (string)$request->input('name');
        $pw = (string)$request->input('pw');

        if ($name === null || $pw === null) {
            return response()->json(GlobalVar::PARAM_LACK_RESPONSE);
        }

        // 如果登录信息没错，就往 request 写入一个 key为uid value为该用户uid 的头
        $user = DB::table('user')->where('nickname', '=', $name)->get()[0];

        if ($user === null) {
            return response()->json(GlobalVar::PARAM_ERROR_RESPONSE);
        }

        if (Hash::check($user->code)) {
            //
        }

        $request->headers->set('uid', $user->id);

        return response()->json([
            'code' => 0,
            'msg' => '登录成功'
        ]);
    }

    /** SSO关联 */
    public function ssoAuth(Request $request)
    {
        $name = (string)$request->input('name');
        $pw = (string)$request->input('pw');

        if ($name === null || $pw === null) {
            return response()->json(GlobalVar::PARAM_LACK_RESPONSE);
        }

        // 如果登录信息没错，就往 request 写入一个 key为uid value为该用户uid 的头
        $user = DB::table('user')->where('nickname', '=', $name)->get()[0];

        if ($user === null) {
            return response()->json(GlobalVar::PARAM_ERROR_RESPONSE);
        }

        if (Hash::check($user->code)) {
            //
        }

        $request->headers->set('uid', $user->id);

        return response()->json([
            'code' => 0,
            'msg' => '登录成功'
        ]);
    }

    /** 修改用户信息 */
    public function edit(Request $request)
    {
        //
    }
}
