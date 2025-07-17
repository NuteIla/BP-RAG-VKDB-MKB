<?php
require_once __DIR__ . '/sign.php';

$MEMORYKB_NAME = 'bp_memorykb_demo12';

$message_english_study = [
    ["role" => "assistant", "content" => "The video just now mentioned pizza, taco, hotdog, and avocado. Which food do you like the best and why?"],
    ["role" => "user", "content" => "I like to eat taco."],
    ["role" => "assistant", "content" => "Great! Tacos are yummy. Have you tried pizza before?"],
    ["role" => "user", "content" => "I don’t know."],
    ["role" => "assistant", "content" => "Please answer in English. Can you tell me in English?"],
    ["role" => "user", "content" => "I don’t understand the question."],
    ["role" => "assistant", "content" => "Let me explain. It means 'Have you eaten pizza before?' Now answer in English, ok?"],
    ["role" => "user", "content" => "No."],
    ["role" => "assistant", "content" => "Well done! You did great. Do you like hotdog?"],
    ["role" => "user", "content" => "Yes."],
    ["role" => "assistant", "content" => "Amazing! You're really good at English. Keep learning the following content."],
    ["role" => "assistant", "content" => "Now try to translate this Chinese sentence into English."],
    ["role" => "user", "content" => "She is angry to her brother."],
    ["role" => "assistant", "content" => "Almost right! The correct preposition is “with”. It's 'She is angry with her brother'.\nHow do you say: 妈妈对弟弟生气了？"],
    ["role" => "user", "content" => "Mom is angry with my brother."],
    ["role" => "assistant", "content" => "Great job! You're doing so well. How do you say: 我奶奶对我生气了？"],
    ["role" => "user", "content" => "Grandmother is angry with me."],
    ["role" => "assistant", "content" => "Fantastic! You're amazing. Keep learning the following content."],
    ["role" => "assistant", "content" => "Now try to translate this Chinese sentence into English.\n他们对迟来的巴士感到生气。How do you say this in English?"],
    ["role" => "user", "content" => "I don't know."],
    ["role" => "assistant", "content" => "Don’t be afraid to make mistakes. The correct answer is: They are angry about the late bus.\nHow do you say: 他们对坏掉的椅子感到生气？"],
    ["role" => "user", "content" => "They are angry about broken chair."],
    ["role" => "assistant", "content" => "Think again — the correct answer is: They are angry about the broken chair.\nHow do you say: 妈妈对吵闹的音乐感到生气？"],
    ["role" => "user", "content" => "Mom is angry about music."],
    ["role" => "assistant", "content" => "Don’t worry about mistakes. The correct answer is: Mom is angry about the noisy music.\nHow do you say: 爸爸对没洗的碗感到生气？"],
    ["role" => "user", "content" => "Dad is angry about 没洗的碗。"],
    ["role" => "assistant", "content" => "Please answer in English?"],
    ["role" => "user", "content" => "How do you say 没洗的碗 in English?"],
    ["role" => "assistant", "content" => "You can say 'unwashed dishes'. Now try translating: 爸爸对没洗的碗感到生气。How would you say it?"],
    ["role" => "user", "content" => "Dad is angry about unwashed dishes."],
    ["role" => "assistant", "content" => "Awesome! How do you say: 爷爷对没关的电视感到生气？"],
    ["role" => "user", "content" => "Grandpa is angry about the working TV."],
    ["role" => "assistant", "content" => "Great job, you've already used the target language point correctly twice. Keep it up!"],
];

$payload = [
    "collection_name" => $MEMORYKB_NAME,
    "session_id" => "message0",
    "messages" => $message_english_study,
    "metadata" => [
        "default_user_id" => "user1",
        "default_assistant_id" => "assistant1",
        "time" => 1747611158000
    ],
    "entities" => [
        [
            "entity_type" => "knowledge_point_demo",
            "entity_scope" => [
                ["id" => "1", "knowledge_point_name" => "taco", "knowledge_point_desc" => "Vocabulary"],
                ["id" => "2", "knowledge_point_name" => "pizza", "knowledge_point_desc" => "Vocabulary"]
            ]
        ]
    ]
];

$apiPath = '/api/memory/messages/add';
$method = 'POST';
$query = [];
$header = [];
$body = json_encode($payload, JSON_UNESCAPED_UNICODE);

try {
    $response = request($apiPath, $method, $query, $header, $body);
    echo $response->getBody();
} catch (Exception $e) {
    echo $e->getMessage();
}