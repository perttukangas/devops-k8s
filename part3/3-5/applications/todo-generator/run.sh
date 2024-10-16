#!/bin/sh

echo Running...

URL=$(curl -s -I https://en.wikipedia.org/wiki/Special:Random | grep -i location | awk '{print $2}' | tr -d '\r')
JSON_PAYLOAD=$(jq -n --arg url "$URL" '{todo: ("Read " + $url)}')

echo $URL
echo $JSON_PAYLOAD

curl -X POST -H "Content-Type: application/json" -d "$JSON_PAYLOAD" http://todo-backend-svc:2346

echo Done.