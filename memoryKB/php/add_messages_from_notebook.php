<?php
require_once __DIR__ . '/sign.php';

// 生成当前毫秒时间戳的函数
function current_milli_time() {
    return intval(round(microtime(true) * 1000));
}

$payload = [
    "collection_name" => "english_learning_memory",
    "session_id" => "english_session_studentA_kp_weather_" . current_milli_time(),
    "messages" => [
        [
            "role" => "assistant",
            "content" => "Okay, let's talk about weather. How would you say \"今天天气多云但很暖和\" in English?"
        ],
        [
            "role" => "user",
            "content" => "Today is cloudy but very warm."
        ],
        [
            "role" => "assistant",
            "content" => "That's a good start! You can also say, \"It's cloudy but warm today\". Your sentence structure was correct. Great job!"
        ]
    ],
    "metadata" => [
        "default_user_id" => "student_A",
        "default_assistant_id" => "tutor_001",
        "time" => current_milli_time() - 120000, // 会话开始时间
        "group_id" => "english_class_101"
    ],
    "entities" => [
        [
            "entity_type" => "english_knowledge_point",
            "entity_scope" => [
                [ // 确保这些字段在 english_knowledge_point schema 中定义为 UseProvided: true 或 IsPrimaryKey: true
                    "id" => 201, // 知识点ID
                    "knowledge_point_name" => "weather related vocabulary",
                    "good_evaluation_criteria" => "Student is able to describe weather conditions"
                ]
            ]
        ]
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
            echo "消息添加成功: " . json_encode($result, JSON_UNESCAPED_UNICODE | JSON_PRETTY_PRINT) . "\n";
        } else {
            echo "消息添加失败: " . json_encode($result, JSON_UNESCAPED_UNICODE | JSON_PRETTY_PRINT) . "\n";
        }
    } else {
        echo "API请求错误: HTTP " . $response->getStatusCode() . " - " . $responseBody . "\n";
    }
} catch (Exception $e) {
    echo "错误: " . $e->getMessage() . "\n";
}
?>
