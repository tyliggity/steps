#!/bin/bash

echo -e "$TEAMS_CONTENT"  > /test.txt

CONTENT=""
while read p; do
    CONTENT="$CONTENT<br />$p"
done </test.txt

#curl -H "Content-Type: application/json" -d "{\"title\": \"$TEAMS_TITLE\", \"text\": \"$CONTENT\"}" $TEAMS_WEBHOOK


curl --request POST \
  --url "$TEAMS_WEBHOOK" \
  --header 'Accept: application/json' \
  --header 'Content-Type: application/json' \
  --data "{
        \"@type\": \"MessageCard\",
        \"@context\": \"http://schema.org/extensions\",
        \"themeColor\": \"0076D7\",
        \"markdown\": \"false\",
        \"summary\": \"Notification from StackPulse\",  
        \"title\": \"$TEAMS_TITLE\", 
        \"text\": \"$CONTENT\", 
        \"sections\": [{
        \"facts\":[
            $TEAMS_TABLE_DATA
            ]}],
            \"potentialAction\": [{
                \"@type\": \"ActionCard\",
                \"name\": \"Details\",
            \"actions\": [{
                \"@type\": \"OpenUri\",
                \"name\": \"$TEAMS_URL_TITLE\",
                \"targets\": [{ \"os\": \"default\", \"uri\": \"$TEAMS_URL\" }]
            }]}]
}"