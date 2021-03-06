#!/usr/bin/env bash

# should have the following environment variables set:
# PORT
# HASH
# ENVIRONMENT
# CLUSTER_ENVIRONMENT
# AMQP_DSN
# LOGGING_DISCORD_WEBHOOK
# DISCORD_TOKENS
# DOCKER_IMAGE_HASH
# CONCURRENT_PROCESSING_LIMIT
# DB_DSN
# LASTFM_KEY
# LASTFM_SECRET
# REDIS_ADDRESS
# REDIS_PASSWORD
# FEATUREFLAG_UNLEASH_URL
# FEATUREFLAG_UNLEASH_INSTANCE_ID
# ERRORTRACKING_RAVEN_DSN
# POLR_BASE_URL
# POLR_API_KEY
# IEXCLOUD_API_SECRET
# DISCORD_API_BASE
# GCLOUD_BUCKET_NAME
# OBJECT_STORAGE_FQDN
# WEVERSE_TOKEN
# INSTAGRAM_SESSION_IDS
# TRELLO_KEY
# TRELLO_TOKEN
# GOOGLE_MAPS_KEY
# DARK_SKY_KEY
# HONEYCOMB_API_KEY
# base64 encoded file content:
# GOOGLE_APPLICATION_CREDENTIALS

template="k8s/manifest.tmpl.yaml"
target="k8s/manifest.yaml"

cp "$template" "$target"
sed -i -e "s|{{PORT}}|$PORT|g" "$target"
sed -i -e "s|{{HASH}}|$HASH|g" "$target"
sed -i -e "s|{{ENVIRONMENT}}|$ENVIRONMENT|g" "$target"
sed -i -e "s|{{AMQP_DSN}}|$AMQP_DSN|g" "$target"
sed -i -e "s|{{LOGGING_DISCORD_WEBHOOK}}|$LOGGING_DISCORD_WEBHOOK|g" "$target"
sed -i -e "s|{{DISCORD_TOKENS}}|$DISCORD_TOKENS|g" "$target"
sed -i -e "s|{{DOCKER_IMAGE_HASH}}|$DOCKER_IMAGE_HASH|g" "$target"
sed -i -e "s|{{CONCURRENT_PROCESSING_LIMIT}}|$CONCURRENT_PROCESSING_LIMIT|g" "$target"
sed -i -e "s|{{DB_DSN}}|$DB_DSN|g" "$target"
sed -i -e "s|{{LASTFM_KEY}}|$LASTFM_KEY|g" "$target"
sed -i -e "s|{{LASTFM_SECRET}}|$LASTFM_SECRET|g" "$target"
sed -i -e "s|{{REDIS_ADDRESS}}|$REDIS_ADDRESS|g" "$target"
sed -i -e "s|{{REDIS_PASSWORD}}|$REDIS_PASSWORD|g" "$target"
sed -i -e "s|{{CLUSTER_ENVIRONMENT}}|$CLUSTER_ENVIRONMENT|g" "$target"
sed -i -e "s|{{FEATUREFLAG_UNLEASH_URL}}|$FEATUREFLAG_UNLEASH_URL|g" "$target"
sed -i -e "s|{{FEATUREFLAG_UNLEASH_INSTANCE_ID}}|$FEATUREFLAG_UNLEASH_INSTANCE_ID|g" "$target"
sed -i -e "s|{{ERRORTRACKING_RAVEN_DSN}}|$ERRORTRACKING_RAVEN_DSN|g" "$target"
sed -i -e "s|{{POLR_BASE_URL}}|$POLR_BASE_URL|g" "$target"
sed -i -e "s|{{POLR_API_KEY}}|$POLR_API_KEY|g" "$target"
sed -i -e "s|{{IEXCLOUD_API_SECRET}}|$IEXCLOUD_API_SECRET|g" "$target"
sed -i -e "s|{{DISCORD_API_BASE}}|$DISCORD_API_BASE|g" "$target"
sed -i -e "s|{{GOOGLE_APPLICATION_CREDENTIALS}}|$GOOGLE_APPLICATION_CREDENTIALS|g" "$target"
sed -i -e "s|{{GCLOUD_BUCKET_NAME}}|$GCLOUD_BUCKET_NAME|g" "$target"
sed -i -e "s|{{OBJECT_STORAGE_FQDN}}|$OBJECT_STORAGE_FQDN|g" "$target"
sed -i -e "s|{{WEVERSE_TOKEN}}|$WEVERSE_TOKEN|g" "$target"
sed -i -e "s|{{INSTAGRAM_SESSION_IDS}}|$INSTAGRAM_SESSION_IDS|g" "$target"
sed -i -e "s|{{TRELLO_KEY}}|$TRELLO_KEY|g" "$target"
sed -i -e "s|{{TRELLO_TOKEN}}|$TRELLO_TOKEN|g" "$target"
sed -i -e "s|{{DARK_SKY_KEY}}|$DARK_SKY_KEY|g" "$target"
sed -i -e "s|{{HONEYCOMB_API_KEY}}|$HONEYCOMB_API_KEY|g" "$target"
sed -i -e "s|{{GOOGLE_MAPS_KEY}}|$GOOGLE_MAPS_KEY|g" "$target"
sed -i -e "s|{{WOLFRAM_APP_ID}}|$WOLFRAM_APP_ID|g" "$target"
