package tasks

// TODO: Add python support via --python_out=.

//go:generate protoc -I=./ -I=../../ -I=${GOPATH}/pkg/mod/github.com/gogo/googleapis@v1.3.0/ -I=${GOPATH}/pkg/mod/github.com/gogo/protobuf@v1.3.1/ --gogoslick_out=plugins=grpc,Mgoogle/api/annotations.proto=github.com/gogo/googleapis/google/api:. --grpc-gateway_out=allow_patch_feature=false,Mgoogle/api/annotations.proto=github.com/gogo/googleapis/google/api:. tasks.proto
