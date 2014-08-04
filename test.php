<?php

function postdata($data)
{
    return join('&', array_map(function($key, $value) {
        return rawurlencode($key).'='.rawurlencode($value);
    }, array_keys($data), array_values($data)));
}

$client = stream_socket_client("unix:///tmp/fetch.sock");
$clientId = "hIsPlzwzep7EaELFurZH-sPe";
$clientSecret = "ncyhVpTBDW4xOu2tbGjDpvQQcQ_WLhos";
$permissions = "manage_project:pure-it-staging";

$now = microtime(true);
fwrite($client, json_encode([
    'jsonrpc' => '2.0',
    'method' => 'fetch.Fetch',
    'params' => [[
        'Url' => "https://auth.sphere.io/oauth/token",
        "Method" => "POST",
        'Body' => base64_encode(postdata([
            'grant_type' => 'client_credentials',
            'scope' => $permissions,
        ])),
        "Headers" => [
            "Content-Type" => ["application/x-www-form-urlencoded"],
            "Authorization" => ["Basic ".base64_encode("$clientId:$clientSecret")],
        ],
    ]],
    'id' => 1
]));

echo $response = fgets($client);
echo microtime(true) - $now, "\n";

$data = json_decode($response, true);
$body = json_decode(base64_decode($data['result']['body']), true);

for ($i = 0; $i < 1000; $i++) {
    $now = microtime(true);
    fwrite($client, json_encode([
        'jsonrpc' => '2.0',
        'method' => 'fetch.Fetch',
        'params' => [[
            'url' => 'https://api.sphere.io',
            'headers' => [
                'Authorization' => ["Bearer {$body['access_token']}"],
            ]
        ]]
    ]));
    $rawResponse = fgets($client);
    echo microtime(true) - $now, "\n";

    sleep(1);
}
