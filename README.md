Public Web URL: https://kqb-stats-parser-5e4m25kezq-uc.a.run.app/

Send a POST request to https://kqb-stats-parser-5e4m25kezq-uc.a.run.app/api
Form Data, Key "images", attach screen shot files. Get a JSON response. 

Run with docker:
docker run \
-p 8080:8080
-e GOOGLE_APPLICATION_CREDENTIALS=/tmp/keys/vision-api.json \
-v $GOOGLE_APPLICATION_CREDENTIALS:/tmp/keys/vision-api.json:ro \
stats-parser:latest