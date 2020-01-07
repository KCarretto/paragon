#!/bin/bash

# make some tags
curl localhost/api/v1/tags/create -d '{"name": "testTag1"}'
curl localhost/api/v1/tags/create -d '{"name": "testTag2"}'
curl localhost/api/v1/tags/create -d '{"name": "testTag3"}'
curl localhost/api/v1/tags/create -d '{"name": "testTag4"}'
curl localhost/api/v1/tags/create -d '{"name": "testTag5"}'


# make some targets
curl localhost/api/v1/targets/create -d '{"name": "testTarget1", "primaryIP": "0.0.0.1", "tags": ["testTag1", "testTag2"]}'
curl localhost/api/v1/targets/create -d '{"name": "testTarget2", "primaryIP": "0.0.0.2", "tags": ["testTag2", "testTag3"]}'
curl localhost/api/v1/targets/create -d '{"name": "testTarget3", "primaryIP": "0.0.0.3", "tags": ["testTag3", "testTag4"]}'
curl localhost/api/v1/targets/create -d '{"name": "testTarget4", "primaryIP": "0.0.0.4", "tags": ["testTag4", "testTag5"]}'

# update a target
curl localhost/api/v1/targets/setTargetFields -d '{"id": 17179869188, "name": "newTestTarget4"}'

# add a credentials
curl localhost/api/v1/targets/addCredential -d '{"targetID": 17179869185, "principal": "root", "secret": "password"}'

# fail a credential
curl localhost/api/v1/credentials/fail -d '{"id": 1 }'

# make some jerbs
curl localhost/api/v1/job_templates/create -d '{"name": "testJob1", "content": "$original1", "tags": []}'
curl localhost/api/v1/job_templates/create -d '{"name": "testJob2", "content": "original2", "tags": ["testTag2", "testTag3"]}'
curl localhost/api/v1/job_templates/create -d '{"name": "testJob3", "content": "original4", "tags": ["testTag4"]}'

# queue some jerbs
curl localhost/api/v1/job_templates/queue -d '{"id": 8589934594, "jobName": "Job 1", "parameters": "{\"original\": \"new\"}"}'

# claim le task
curl localhost/api/v1/tasks/claim -d '{"id": 21474836481}'
