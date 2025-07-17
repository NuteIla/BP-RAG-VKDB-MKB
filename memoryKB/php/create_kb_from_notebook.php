<?php
require_once __DIR__ . '/sign.php';
$memory_kb_name = 'english_learning_memory23';
$CustomEventTypeShema = [
    [
        "EventType" => "english_study",
        "Description" => "Record the Q&A and scoring between the assistant and the student in an English study session",
        "Properties" => [
            [
                "PropertyName" => "knowledge_point_name",
                "PropertyValueType" => "string",
                "Description" => "The name of the knowledge point involved in the current conversation"
            ],
            [
                "PropertyName" => "question",
                "PropertyValueType" => "string",
                "Description" => "The question raised by the assistant"
            ],
            [
                "PropertyName" => "answer",
                "PropertyValueType" => "string",
                "Description" => "The student's answer"
            ],
            [
                "PropertyName" => "rating_score",
                "PropertyValueType" => "float32",
                "Description" => "The numerical score for the student's answer, full score is 10"
            ]
        ]
    ]
];
$CustomEntityTypeSchemas = [
    [
        "EntityType" => "english_knowledge_point",
        "AssociatedEventTypes" => ["english_study"],
        "Description" => "Track students' learning progress on specific English knowledge points",
        "Properties" => [
            [
                "PropertyName" => "id",
                "PropertyValueType" => "int64",
                "Description" => "Unique primary key ID of the knowledge point",
                "IsPrimaryKey" => true,
                "UseProvided" => true
            ],
            [
                "PropertyName" => "knowledge_point_name",
                "PropertyValueType" => "string",
                "Description" => "Specific name of the knowledge point (e.g., 'weather-related vocabulary')",
                "UseProvided" => true
            ],
            [
                "PropertyName" => "good_evaluation_criteria",
                "PropertyValueType" => "string",
                "Description" => "Criteria for judging an answer as 'good'",
                "UseProvided" => true
            ],
            [
                "PropertyName" => "rating_score_max",
                "PropertyValueType" => "float32",
                "Description" => "The highest numerical score obtained on this knowledge point",
                "AggregateExpression" => [
                    "Op" => "MAX",
                    "EventType" => "english_study",
                    "EventPropertyName" => "rating_score"
                ]
            ]
        ]
    ]
];
$payload = [
    "CollectionName" => $memory_kb_name,
    "Description" => "A memory base for recording and analyzing students' English learning process",
    "CustomEventTypeSchemas" => $CustomEventTypeShema,
    "CustomEntityTypeSchemas" => $CustomEntityTypeSchemas
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
