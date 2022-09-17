pwd=$(shell pwd)
input_file=input_file.csv

test:
	go test -v --cover ./...

build: test
	cd cmd/simplebook && go build

run: build
	./cmd/simplebook/simplebook $(input_file)

clean:
	rm cmd/simplebook/simplebook || true
	docker image rm github.com/rodoufu/simple-orderbook:latest || true

build_docker:
	docker build -t github.com/rodoufu/simple-orderbook:latest .

run_docker: build_docker
	docker run --rm --name simplebook -v $(pwd)/$(input_file):/app/$(input_file) -it github.com/rodoufu/simple-orderbook:latest $(input_file)
