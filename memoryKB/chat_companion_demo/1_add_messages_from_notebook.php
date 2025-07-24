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
    "session_id" => "chat_companion_memory_session_uid_" . current_milli_time(),
    "messages" => [
        [
            "role" => "user",
            "content" => " I really love watching competitions. There's another one tomorrow!"
        ],
        [
            "role" => "user",
            "content" => "Athlete A has been training for six whole months just for this tennis match with Athelete B."
        ],
        [
            "role" => "assistant",
            "content" => "That's so inspiring! Then it's pretty much a sure win this time, right?"
        ],
        [
            "role" => "user",
            "content" => "I think so too — Athlete B barely prepared at all. I always enjoy watching tennis matches"
        ],
        [
            "role" => "assistant",
            "content" => "Haha, let's look forward to Xiao A's performance together!"
        ]
    ],
    "metadata" => [
        "default_user_id" => $user_id,
        "default_assistant_id" => $chat_companion_id,
        "time" => current_milli_time() - 120000, // Session start time
        // "group_id" => "group_chat_001" //special for group chat case. 
    ]
];

$apiPath = '/api/memory/messages/add';
$method = 'POST';
$query = [];
$header = [
    "X-Request-ID" => "add-msg-" . current_milli_time()
];
$body = json_encode($payload, JSON_UNESCAPED_UNICODE);

try {
    $response = request($apiPath, $method, $query, $header, $body);
    $responseBody = $response->getBody()->getContents();
    $result = json_decode($responseBody, true);

    if ($response->getStatusCode() == 200) {
        if ($result && isset($result['code']) && $result['code'] == 0) {
            echo "Message added successfully: " . json_encode($result, JSON_UNESCAPED_UNICODE | JSON_PRETTY_PRINT) . "\n";
        } else {
            echo "Message addition failed: " . json_encode($result, JSON_UNESCAPED_UNICODE | JSON_PRETTY_PRINT) . "\n";
        }
    } else {
        echo "API request error: HTTP " . $response->getStatusCode() . " - " . $responseBody . "\n";
    }
} catch (Exception $e) {
    echo "Error: " . $e->getMessage() . "\n";
}
