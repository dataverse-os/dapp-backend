#!/usr/bin/env bash

export PRIVATE_KEY="$(composedb did:generate-private-key)"
echo $PRIVATE_KEY
export DID="$(composedb did:from-private-key $PRIVATE_KEY)"
echo $DID

jq --arg did "${DID}" '."http-api"."admin-dids"=[$did]' /data/daemon.config.json | sponge /data/daemon.config.json

echo "DID_PRIVATE_KEY=${PRIVATE_KEY}" > /data/.private-key.env