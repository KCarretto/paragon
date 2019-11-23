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

# add a credentials
curl localhost/api/v1/targets/addCredential -d '{"targetID": 12884901889, "principal": "root", "secret": "password"}'

# make some jerbs
curl localhost/api/v1/jobs/create -d '{"name": "testJob1", "content": "original1", "tags": []}'
curl localhost/api/v1/jobs/create -d '{"name": "testJob2", "content": "original2", "tags": ["testTag2", "testTag3"]}'
curl localhost/api/v1/jobs/create -d '{"name": "testJob3", "content": "original4", "tags": ["testTag4"]}'

# queue some jerbs
curl localhost/api/v1/jobs/queue -d '{"id": 4294967298}'
curl localhost/api/v1/jobs/queue -d '{"id": 4294967299}'

# claim le task
curl localhost/api/v1/tasks/claim -d '{"id": 17179869187}'
