module github.com/drifterz13/go-services/internal/user

go 1.16

require (
	github.com/drifterz13/go-services/internal/common v0.0.0-20210830010347-0153a1d358e4
	google.golang.org/grpc v1.40.0
	google.golang.org/protobuf v1.27.1
)

replace github.com/drifterz13/go-services/internal/common => "../common"