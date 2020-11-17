build:
	go build .; \
	  cd frontend; \
	  npm run build

install:
	go install .; \
	cd frontend; \
	npm ci

test:
	go test -v .
