<?php

namespace App\Http\Controllers;

use Illuminate\Support\Facades\DB;
use Illuminate\Support\Collection;
use Illuminate\Http\Request;
use App\GlobalVar;

class GroupController extends Controller
{

    /** 获取列表 */
    public function list(Request $request)
    {
        $result = DB::table('group')->get();

        return response()->json([
            'code' => 0,
            'result' => $result,
        ]);
    }
}
