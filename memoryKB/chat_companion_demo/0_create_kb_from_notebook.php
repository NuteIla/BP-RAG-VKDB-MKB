<?php
require_once __DIR__ . '/sign.php';
$memory_kb_name = 'chat_companion_123';


$payload = [
    "CollectionName" => $memory_kb_name,
    "Description" => "A memory base for recording and analyzing user profile with the chat companion",
    'BuiltinEventTypes' => ["sys_event_v1", "sys_profile_collect_v1"],
    'BuiltinEntityTypes' => ["sys_profile_v1"]
];

$apiPath = '/api/memory/collection/create';
$method = 'POST';
$query = [];
$header = [];
$body = json_encode($payload, JSON_UNESCAPED_UNICODE);

try {
    $response = request($apiPath, $method, $query, $header, $body);
    print_r($response->getBody()->getContents());
} catch (Exception $e) {
    echo $e->getMessage();
}
