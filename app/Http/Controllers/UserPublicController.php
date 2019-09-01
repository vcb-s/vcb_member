<?php

namespace App\Http\Controllers;

class UserPublicController extends Controller
{
    public function index($id)
    {
        return response()->json([
            'id' => $id,
            'hello' => 'world'
        ]);
    }
}
