<?php

namespace App\Http\Controllers;

class UserController extends Controller
{
    public function index($id)
    {
        return response()->json([
            'id' => $id,
            'hello' => 'world'
        ]);
    }
}
