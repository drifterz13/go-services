.PHONY: task_proto
task_proto:
	protoc --proto_path=proto \
	--go_out=internal/common/genproto/task --go_opt=paths=source_relative \
	--go-grpc_out=internal/common/genproto/task --go-grpc_opt=paths=source_relative \
	task.proto

.PHONY: user_proto
user_proto:
	protoc --proto_path=proto \
	--go_out=internal/common/genproto/user --go_opt=paths=source_relative \
	--go-grpc_out=internal/common/genproto/user --go-grpc_opt=paths=source_relative \
	user.proto
