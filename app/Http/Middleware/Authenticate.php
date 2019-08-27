<?php

namespace App\Http\Middleware;

use Jose\Component\Core\JWK;
use Jose\Easy\Build;

use Closure;

class Authenticate
{
    public function __construct()
    {
        //
    }
    /**
     * Handle an incoming request.
     *
     * @param  \Illuminate\Http\Request  $request
     * @param  \Closure  $next
     * @return mixed
     */
    public function handle($request, Closure $next)
    {
        return $next($request->headers->add([
            'uid' => '123'
        ]));
    }
}
