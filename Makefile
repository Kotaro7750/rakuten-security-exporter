PROTO=./rakuten-security-scraper.proto
GO_DIR=./exporter/proto
PYTHON_DIR=./scraper

PYTHON_ARTIFACTS=$(PYTHON_DIR)/rakuten_security_scraper_pb2.py $(PYTHON_DIR)/rakuten_security_scraper_pb2.pyi $(PYTHON_DIR)/rakuten_security_scraper_pb2_grpc.py
GO_ARTIFACTS=$(GO_DIR)/rakuten-security-scraper.pb.go $(GO_DIR)/rakuten-security-scraper_grpc.pb.go

all: $(PYTHON_ARTIFACTS) $(GO_ARTIFACTS)

$(PYTHON_ARTIFACTS): $(PROTO)
	python -m grpc_tools.protoc -I. --python_out=$(PYTHON_DIR) --pyi_out=$(PYTHON_DIR) --grpc_python_out=$(PYTHON_DIR) $^

$(GO_ARTIFACTS): $(PROTO)
	protoc -I=. --go_out=$(GO_DIR) --go_opt=paths=source_relative --go-grpc_out=$(GO_DIR) --go-grpc_opt=paths=source_relative $^
