module github.com/drifterz13/go-services/internal/task

go 1.16

require (
	github.com/drifterz13/go-services/internal/common v0.0.0-20210830010347-0153a1d358e4
	go.mongodb.org/mongo-driver v1.7.2
	golang.org/x/net v0.0.0-20210825183410-e898025ed96a // indirect
	google.golang.org/genproto v0.0.0-20210828152312-66f60bf46e71 // indirect
	google.golang.org/grpc v1.40.0
	google.golang.org/protobuf v1.27.1
)

replace github.com/drifterz13/go-services/internal/common => ../common
