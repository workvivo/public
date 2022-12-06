<?php

namespace App\Http\Controllers;

use App\Http\Controllers\Controller;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Auth;
use Firebase\JWT\JWT;
use Firebase\JWT\Key;
use App\Models\User;
use Illuminate\Support\Facades\Cache;


class Social extends Controller
{
    
    public function social(Request $request)
    {
        $WVPublicKeySource = getenv("WORKVIVO_APP_URL_PUBLIC_KEYS").'?keyid=';
        $keyid = $request->keyid;
        /* look up the public key */
        $publicKey = file_get_contents($WVPublicKeySource.$keyid); 
        
        /* 
        confirm the JWT with the public key, a good library will return an error 1
        if the JWT has expired do check your Library, recommend looking 
        here https://jwt.io/libraries for the fewatures of your Library
         */
        try {
            $decoded = JWT::decode($request->jwt, new Key(json_decode($publicKey, true)[$keyid], 'RS256'));
         }
         catch (\Exception $e) {
            dd('oops, an error ', $e->getMessage());
         }
         /* 
         At this point the JWT expiry shoud have been confirmed by the Library, now we can check other items;
         - state 
            - this acts as a nonce confirm it does not already exist in cache
         - iss
            - this should be as you expect, for such as;
                - {customer}.workvivo.co 
                - {customer}.workvivo.com
                - {customer}.workvivo.io
        - workvivo_id / third_party_id / email
            - confirm these are bound to the same user
        - Valid Relay
            - you could validate the Relay too, below we do not.
         */

         /*state check for a nonce
          - 0 means nothing in cache in this implementation.
          - when cache is same as the state value this means it already exists and so we should block
         */
         if ($this->getState($decoded->state) !== 0) {
             abort(403, 'nonce failed');
            }
        
        /* iss check
         - ensure the host issuing the request is as expected
        */
        if (getenv('WORKVIVO_ISSUE_HOST') !== $decoded->iss ) {
            abort(403, 'issue host failed');
        }

        /* workvivo_id / third_party_id / email
         - the email and id must match a record otherwise it fails
        */
         $theUserId = User::whereEmail($decoded->email)->whereId($decoded->third_party_id)->firstOrFail(['id']);

        /*
            ALL GOOD! Let's login and redirect
        */
         Auth::loginUsingId($theUserId->id);

        return redirect(urldecode($request->relay),302,[],true);
        
    }

    public function getState($statekey)
    {
        return Cache::get($statekey, function () use ($statekey) {
            Cache::put($statekey, 300); // add now to cache
            return 0;
        });
    }

    }
