build:
	go build .; \
	  cd frontend; \
	  npm run build;

install:
	go install .
test:
	go test -v .
