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
        $current = (int)$request->input('current', 1);
        $pageSize = (int)$request->input('pageSize', 20);
        $group = (string)(int)$request->input('group');
        $retired = (bool)$request->input('retired');
        $sticky = (bool)$request->input('sticky');

        $allUser = DB::table('user');

        if ($group) {
            $allUser->where('user.group', 'like', '%' . $group . '%');
        }
        if ($retired) {
            $allUser->where('user.retired', '>', 0);
        }
        if ($sticky) {
            $allUser->where('user.order', '>', 0);
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
            )
            ->orderBy('user.order', 'desc')
            ->orderBy('user.id', 'asc');

        $total = $allUser->count();

        $result = $allUser
            ->offset($pageSize * ($current - 1))
            ->limit($pageSize)
            ->get();

        $result->each(function ($item, $index) {
            $item->group = array_map(function ($id) {
                return (int)$id;
            }, explode(',', $item->group));

            $item->order = (int)$item->order;
            $item->retired = (int)$item->retired;
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
}
