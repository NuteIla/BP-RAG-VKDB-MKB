<?php
require_once __DIR__ . '/sign.php';

// Function to generate current millisecond timestamp
function current_milli_time()
{
    return intval(round(microtime(true) * 1000));
}
$memory_kb_name = 'chat_companion_123';
$user_id='1234';
$chat_companion_id='1234567';

$payload = [
    "collection_name" => $memory_kb_name,
    "query" => "Guess who ended up winning the match?",
    "limit" => 5,
    "filter" => [
        "user_id" => $user_id,
        "memory_type" => ["sys_profile_v1"]
    ]
];

$apiPath = '/api/memory/search';
$method = 'POST';
$query = [];
$header = [
    "X-Request-ID" => "search-profile-" . current_milli_time()
];
$body = json_encode($payload, JSON_UNESCAPED_UNICODE);

try {
    $response = request($apiPath, $method, $query, $header, $body);
    $responseBody = $response->getBody()->getContents();
    $result = json_decode($responseBody, true);

    if ($response->getStatusCode() == 200) {
        if ($result && isset($result['code']) && $result['code'] == 0) {
            echo "User profile query successful: " . json_encode($result, JSON_UNESCAPED_UNICODE | JSON_PRETTY_PRINT) . "\n";
        } else {
            echo "User profile query failed: " . json_encode($result, JSON_UNESCAPED_UNICODE | JSON_PRETTY_PRINT) . "\n";
        }
    } else {
        echo "API request error: HTTP " . $response->getStatusCode() . " - " . $responseBody . "\n";
    }
} catch (Exception $e) {
    echo "Error: " . $e->getMessage() . "\n";
}