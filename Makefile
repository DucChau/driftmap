.PHONY: build install test clean lint

BINARY := driftmap
GO := go

build:
	$(GO) build -ldflags "-s -w" -o $(BINARY) .

install:
	$(GO) install .

test:
	$(GO) test ./... -v -race

test-short:
	$(GO) test ./... -short

lint:
	@which golangci-lint > /dev/null 2>&1 || (echo "golangci-lint not found; install from https://golangci-lint.run" && exit 1)
	golangci-lint run ./...

clean:
	rm -f $(BINARY)

run-example:
	@echo "==> Capturing baseline..."
	./$(BINARY) capture baseline
	@echo ""
	@echo "==> Listing snapshots..."
	./$(BINARY) list
	@echo ""
	@echo "==> Showing latest snapshot..."
	./$(BINARY) show latest
