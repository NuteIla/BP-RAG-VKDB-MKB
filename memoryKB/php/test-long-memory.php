
<?php
/**
* READ ME: 
* 1. PLEASE INSTALL COMPOSER 
* 2. RUN 'composer require volcengine/volcengine-php-sdk' -- IN SHELL TO CREATE THE VENDOR AND COMPOSER FILES
* 3. PLEASE GO TO SIGN.PHP TO FILL IN YOUR $AK AND $SK FIELDS
 */
ini_set('log_errors', true);
ini_set('error_log', dirname(__FILE__).'/logs/test-long-memory-'.date("Ymd").'.log');

error_log(__FILE__.": IN ************************************************************");
require_once __DIR__ . '/sign.php';


$response = "ERROR";
try {

    error_log(__FILE__.": START");

    createMemory();

    error_log(__FILE__.": END");

}
catch (Exception $e) {
    $response = "EXCEPTION ".$e->getMessage();
    error_log(__FILE__.": {$response}");
}

exit;
function createMemory() {

    $MEMORYKB_NAME = 'wdm_memory_test_1';
    $EVENT_TYPES = [
        [
            "EventType" => "study_demo",
            "Version" => "1",
            "Description" => "The assistant is teaching the user and recording all the assistant’s questions and the user’s answers, along with a rating for each answer.",
            "Properties" => [
                ["PropertyName" => "knowledge_point_name", "PropertyValueType" => "string", "Description" => "The knowledge point collected"],
                ["PropertyName" => "question", "PropertyValueType" => "string", "Description" => "The original question from the assistant"],
                ["PropertyName" => "answer", "PropertyValueType" => "string", "Description" => "The original answer from the user"],
                ["PropertyName" => "is_user_answered", "PropertyValueType" => "bool", "Description" => "Whether the user answered the question"],
                ["PropertyName" => "rating", "PropertyValueType" => "string", "Description" => "Rating of the user’s answer, can be either 'bad' or 'good'"],
                ["PropertyName" => "rating_reasoning", "PropertyValueType" => "string", "Description" => "Reasoning behind the rating"],
                ["PropertyName" => "issues", "PropertyValueType" => "list<string>", "Description" => "Points that need attention, keep up to 3"],
                ["PropertyName" => "rating_score", "PropertyValueType" => "int64", "Description" => "Score for the user’s answer: 0 = bad, 1 = good"],
                ["PropertyName" => "has_been_taught", "PropertyValueType" => "int64", "Description" => "Whether the user has been taught this knowledge point: 1 = yes, 0 = no"],
                ["PropertyName" => "is_answer_good", "PropertyValueType" => "int64", "Description" => "Whether the answer is correct: 1 = yes, 0 = no"],
                ["PropertyName" => "is_answer_bad", "PropertyValueType" => "int64", "Description" => "Whether the answer is wrong: 1 = yes, 0 = no"]
            ],
            "ValidationExpression" => "is_user_answered==True"
        ]
    ];
    $ENTITY_TYPES = [
        [
            "EntityType" => "knowledge_point_demo",
            "Version" => "1",
            "AssociatedEventTypes" => ["study_demo"],
            "Role" => "user",
            "Description" => "Knowledge point",
            "Properties" => [
                ["PropertyName" => "id", "PropertyValueType" => "string", "Description" => "Primary key", "IsPrimaryKey" => true, "UseProvided" => true],
                ["PropertyName" => "knowledge_point_desc", "PropertyValueType" => "string", "Description" => "Description of the knowledge point (required)", "UseProvided" => true],
                ["PropertyName" => "knowledge_point_name", "PropertyValueType" => "string", "Description" => "Name of the knowledge point", "IsPrimaryKey" => true, "UseProvided" => true],
                ["PropertyName" => "difficulty_level", "PropertyValueType" => "list<int64>", "Description" => "Difficulty level from 1 to 5 (1 = easiest, 5 = hardest, required)"],
                ["PropertyName" => "knowledge_start_time", "PropertyValueType" => "int64", "Description" => "Time the knowledge point was first taught (timestamp: YYYYMMDDHHMMSS)"],
                ["PropertyName" => "rating_score_max", "PropertyValueType" => "float32", "Description" => "Maximum rating score"],
                ["PropertyName" => "has_been_taught", "PropertyValueType" => "int64", "Description" => "Whether the user has been taught this knowledge point"],
                ["PropertyName" => "answer_good_count", "PropertyValueType" => "int64", "Description" => "Number of correct answers"],
                ["PropertyName" => "answer_bad_count", "PropertyValueType" => "int64", "Description" => "Number of wrong answers"],
                ["PropertyName" => "rating_score_sum", "PropertyValueType" => "float32", "Description" => "Total rating score"],
                ["PropertyName" => "issues", "PropertyValueType" => "list<string>", "Description" => "Issues needing attention, keep up to 5"],
                ["PropertyName" => "count", "PropertyValueType" => "int64", "Description" => "Count"]
            ]
        ]
    ];
    $data = [
        'CollectionName' => $MEMORYKB_NAME,
        'Description' => "WDM memory test description",
        'BuiltinEventTypes' => ['sys_common'],
        'CustomEventTypeSchemas' => $EVENT_TYPES,
        'CustomEntityTypeSchemas' => $ENTITY_TYPES,
    ];
    error_log(__FILE__.": DATA=".json_encode($data,JSON_PRETTY_PRINT));
    $apiPath = "/api/memory/collection/create";
    $body = json_encode($data, JSON_UNESCAPED_UNICODE);
    $method = 'POST';
    $query = [];
    $header = [];
    try {
        $response = request($apiPath, $method, $query, $header, $body);
        print_r($response->getBody()->getContents());
        error_log(__FILE__.": response=".print_r($response,true));
    } catch (Exception $e) {
        print_r($e->getMessage());
    }
}


?>