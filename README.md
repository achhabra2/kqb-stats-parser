Run with docker:
docker run \
-p 8080:8080
-e GOOGLE_APPLICATION_CREDENTIALS=/tmp/keys/vision-api.json \
-v $GOOGLE_APPLICATION_CREDENTIALS:/tmp/keys/vision-api.json:ro \
stats-parser:latest