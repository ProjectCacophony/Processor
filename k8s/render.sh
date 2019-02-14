#!/usr/bin/env bash

# should have the following environment variables set:
# PORT
# ENVIRONMENT
# AMQP_DSN
# LOGGING_DISCORD_WEBHOOK
# DISCORD_TOKEN_<BOT ID>
# DOCKER_IMAGE_HASH

template="k8s/manifest.tmpl.yaml"
target="k8s/manifest.yaml"

cp "$template" "$target"
sed -i -e "s|{{PORT}}|$PORT|g" "$target"
sed -i -e "s|{{ENVIRONMENT}}|$ENVIRONMENT|g" "$target"
sed -i -e "s|{{AMQP_DSN}}|$AMQP_DSN|g" "$target"
sed -i -e "s|{{LOGGING_DISCORD_WEBHOOK}}|$LOGGING_DISCORD_WEBHOOK|g" "$target"
sed -i -e "s|{{DISCORD_TOKEN_545676292154982400}}|$DISCORD_TOKEN_545676292154982400|g" "$target"
sed -i -e "s|{{DOCKER_IMAGE_HASH}}|$DOCKER_IMAGE_HASH|g" "$target"
