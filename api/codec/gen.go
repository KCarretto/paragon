package codec

// TODO: Add python support via --python_out=.

//go:generate protoc -I=./ -I=../../ -I=${GOPATH}/pkg/mod/github.com/gogo/googleapis@v1.3.0/ -I=${GOPATH}/pkg/mod/github.com/gogo/protobuf@v1.3.1/ --gogoslick_out=plugins=grpc,Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types,Mgoogle/api/annotations.proto=github.com/gogo/googleapis/google/api:. agent.proto server.proto
