proto-generate:
	@echo "Generating files..."
	protoc --go_out=./pkg/pb --go-grpc_out=./pkg/pb proto/*.proto

clean:
	@echo "Cleaning up..."
	rm -rf *.go *.pb.go *.proto.go *.grpc.pb.go *.proto.txt *.txt.tmp

.PHONY: all clean