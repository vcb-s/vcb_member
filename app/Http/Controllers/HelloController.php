<?php

namespace App\Http\Controllers;

use Illuminate\Http\Request;

class HelloController extends Controller
{
    public function index(Request $request)
    {
        $allHeader = $request->headers->all();
        $request->headers->set('uid', '1238');
        return response()->json($allHeader);
    }
}
