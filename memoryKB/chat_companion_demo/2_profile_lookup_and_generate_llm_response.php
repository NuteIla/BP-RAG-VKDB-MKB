<?php
require_once __DIR__ . '/sign.php';

// Function to generate current millisecond timestamp
function current_milli_time()
{
    return intval(round(microtime(true) * 1000));
}

$memory_kb_name = 'chat_companion_123';
$user_id = '1234';
$chat_companion_id = '1234567';

// Step 1: Query user profile memory first
$search_payload = [
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
    "X-Request-ID" => "search-memory-" . current_milli_time()
];
$body = json_encode($search_payload, JSON_UNESCAPED_UNICODE);

try {
    // Query memory
    $response = request($apiPath, $method, $query, $header, $body);
    $responseBody = $response->getBody()->getContents();
    $result = json_decode($responseBody, true);

    $memories = [];
    if ($response->getStatusCode() == 200 && $result && isset($result['code']) && $result['code'] == 0) {
        // Extract memory_info from all results
        if (isset($result['data']['result_list'])) {
            foreach ($result['data']['result_list'] as $item) {
                if (isset($item['memory_info'])) {
                    $memories[] = $item['memory_info'];
                }
            }
        }
        echo "Memory query successful, found " . count($memories) . " memories\n";
    } else {
        echo "Memory query failed: " . json_encode($result, JSON_UNESCAPED_UNICODE | JSON_PRETTY_PRINT) . "\n";
    }

    // Step 2: Use memory to generate System Prompt
    $memories_text = implode("\n", $memories);
    $system_prompt = "You are an AI assistant with excellent memory who can remember historical conversations with users.\n" .
                    "Please refer to the background information below and naturally continue the conversation in a conversational manner, as if you really remember what happened before.\n\n" .
                    "[Background Information]\n" . $memories_text;

    // Step 3: Call ModelARK LLM API to generate response
    $ark_payload = [
        "model" => "doubao-seed-1-6-250615",
        "messages" => [
            [
                "role" => "system",
                "content" => $system_prompt
            ],
            [
                "role" => "user",
                "content" => "Guess who ended up winning the match?"
            ]
        ]
    ];

    // Note: You need to modify this to ARK API configuration
    // You need to adjust the following parameters according to the actual ARK API configuration
    $ark_api_path = '/v1/chat/completions'; // ARK API path
    $ark_method = 'POST';
    $ark_query = [];
    $ark_header = [
        "X-Request-ID" => "llm-generate-" . current_milli_time(),
        "Authorization" => "Bearer your_ark_api_key" // Replace with your ARK API Key
    ];
    $ark_body = json_encode($ark_payload, JSON_UNESCAPED_UNICODE);

    echo "\nGenerating LLM response...\n";
    echo "System Prompt: " . $system_prompt . "\n\n";

    // This needs to be called according to the actual ARK API configuration
    // Since ARK API may require different authentication methods, you may need to modify the request function or create a new request function
    
    // Temporary output, when actually using, you need to configure the correct ARK API call
    echo "[Simulated LLM Response] Based on the memory information, I remember you mentioned earlier that Athlete A trained for six months, while Athlete B barely prepared. Considering this situation, I guess Athlete A should have won the match! Did you watch the game? Was the result as we expected?\n";
    
    /*
    // Actual ARK LLM API call code (requires correct authentication information configuration)
    $ark_response = request($ark_api_path, $ark_method, $ark_query, $ark_header, $ark_body);
    $ark_response_body = $ark_response->getBody()->getContents();
    $ark_result = json_decode($ark_response_body, true);
    
    if ($ark_response->getStatusCode() == 200) {
        if (isset($ark_result['choices'][0]['message']['content'])) {
            echo "LLM Response: " . $ark_result['choices'][0]['message']['content'] . "\n";
        } else {
            echo "LLM response parsing failed: " . json_encode($ark_result, JSON_UNESCAPED_UNICODE | JSON_PRETTY_PRINT) . "\n";
        }
    } else {
        echo "LLM API request error: HTTP " . $ark_response->getStatusCode() . " - " . $ark_response_body . "\n";
    }
    */

} catch (Exception $e) {
    echo "Error: " . $e->getMessage() . "\n";
}