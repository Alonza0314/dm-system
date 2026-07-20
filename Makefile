.PHONY: backend frontend tidy lint test run clean docker

BACKEND_SRC := $(shell find backend -name "*.go")
FRONTEND_SRC := $(shell find frontend -type f ! -path "frontend/dist/*" ! -path "frontend/node_modules/*")
FRONTEND_STAMP := build/frontend/.stamp

all: backend frontend

build/dm: $(BACKEND_SRC)
	@echo "[+] Building backend..."
	mkdir -p build
	cd backend && go build -o ../build/dm .
	@echo "[✔] Backend build finished"

build/frontend: $(FRONTEND_SRC)
	@echo "[+] Installing frontend deps..."
	cd frontend && yarn install

	@echo "[+] Building frontend..."
	cd frontend && yarn build

	@echo "[✔] Frontend build finished"
	@mkdir -p build/frontend
	@cp -r frontend/dist/. build/frontend/
	@touch $(FRONTEND_STAMP)

backend:
	@if [ -f build/dm ]; then \
		if [ -z "$$(find backend -name '*.go' -newer build/dm)" ]; then \
			echo "[✔] backend is up-to-date, no build needed"; \
			exit 0; \
		fi; \
	fi; \
	$(MAKE) build/dm

frontend:
	@if [ -f $(FRONTEND_STAMP) ]; then \
		if [ -z "$$(find frontend -type f -newer $(FRONTEND_STAMP))" ]; then \
			echo "[✔] frontend is up-to-date, no build needed"; \
			exit 0; \
		fi; \
	fi; \
	$(MAKE) build/frontend

tidy:
	@echo "[+] Tidying backend dependencies..."
	cd backend && go mod tidy
	@echo "[✔] Backend dependencies tidied"

lint:
	@echo "[+] Linting backend..."
	cd backend && golangci-lint run
	@echo "[✔] Backend linting finished"

test:
	cd backend && go test ./... -v

run:
	./build/dm -c config.yaml

clean:
	rm -rf build

docker:
	./docker/build.sh

dockertest:
	./docker/build.sh test
