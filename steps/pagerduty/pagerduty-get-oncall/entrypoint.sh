#!/bin/bash

curl -s --location --request GET 'https://api.pagerduty.com/oncalls' \
--header "Authorization: Token token=$PAGERDUTY_TOKEN" \
--header 'Accept: application/vnd.pagerduty+json;version=2'