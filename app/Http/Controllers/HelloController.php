<?php

namespace App\Http\Controllers;

use Illuminate\Http\Request;
use App\GlobalVar;

class HelloController extends Controller
{
    public function index(Request $request)
    {
        $allHeader = $request->headers->all();
        $request->headers->set('uid', '1238');
        $allHeader['uuid1'] = GlobalVar::UUID();
        $allHeader['uuid2'] = GlobalVar::UUID();
        $allHeader['uuid3'] = GlobalVar::UUID();
        return response()->json($allHeader);
    }
}
