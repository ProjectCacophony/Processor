#!/usr/bin/env bash

# should have the following environment variables set:
# PORT
# ENVIRONMENT
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

template="k8s/manifest.tmpl.yaml"
target="k8s/manifest.yaml"

cp "$template" "$target"
sed -i -e "s|{{PORT}}|$PORT|g" "$target"
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
