PROTO_DIR=core/pb
PROTO_FILE=$(PROTO_DIR)/defs.proto
OUT_DIR=$(PROTO_DIR)

all: proto

.PHONY: build

build:
	(cd cli/replica && go build -o ../../bin/replica)
	(cd cli/client  && go build -o ../../bin/client)

proto:
	./etc/protoc-mac-arm64 --go_out=. --go_opt=paths=source_relative \
	       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
	       $(PROTO_FILE)

clean:
	rm -f $(OUT_DIR)/*.go
