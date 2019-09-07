<?php

/*
|--------------------------------------------------------------------------
| Application Routes
|--------------------------------------------------------------------------
|
| Here is where you can register all of the routes for an application.
| It is a breeze. Simply tell Lumen the URIs it should respond to
| and give it the Closure to call when that URI is requested.
|
*/

$router->get(
    '/',
    [
        'as' => 'home',
        'uses' => 'HelloController@index',
    ]
);
$router->get(
    '/home',
    [
        'uses' => 'HelloController@index',
    ]
);

// $router->get(
//     '/login',
//     [
//         'as' => 'login',
//         'middleware' => ['tokenIssue'],
//         'uses' => 'UserController@login',
//     ]
// );

$router->get(
    '/user/list',
    [
        'uses' => 'UserController@list',
    ]
);

$router->get(
    '/group/list',
    [
        'uses' => 'GroupController@list',
    ]
);
