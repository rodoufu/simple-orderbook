test:
	go test -v --cover ./...

build:
	cd cmd/simplebook && go build

run: build
	./cmd/simplebook/simplebook

clean:
	rm cmd/simplebook/simplebook || true
	docker image rm github.com/rodoufu/simple-orderbook:latest || true

build_docker:
	docker build -t github.com/rodoufu/simple-orderbook:latest .

pwd=$(shell pwd)
run_docker: build_docker
	docker run --rm --name simplebook -v $(pwd)/input_file.csv:/app/input_file.csv -it github.com/rodoufu/simple-orderbook:latest
