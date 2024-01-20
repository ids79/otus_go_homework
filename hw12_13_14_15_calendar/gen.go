package gen

//go:generate rm -rf internal/pb
//go:generate mkdir -p internal/pb
//go:generate protoc --proto_path=api --go_out=internal/pb --go-grpc_out=internal/pb api/EventService.proto
