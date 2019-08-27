<?php

namespace App\Http\Controllers;

class UserController extends Controller
{
    public function __construct()
    {
        $this->middleware('auth');
    }
    public function show($id)
    {
        return response()->json([
            'id' => $id,
            'hello' => 'world'
        ]);
    }
}
