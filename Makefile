SHELL=bash
.SHELLFLAGS=-euo pipefail -c

run: example
	./example

run-wazero: 
	mkdir -p output && echo 'from host'>output/from-host.txt
	wazero \
		run \
		-cachedir=.cache \
		-mount=$$PWD/lib/python3.13:/usr/local/lib/python3.13:ro \
		-mount=$$PWD/lib/protobuf-5.26.1/google:/usr/local/lib/google:ro\
		-mount=$$PWD/output:/output \
		-mount=$$PWD/s:/s \
		python.wasm s/main.py aaa bbb

example: python.wasm lib *.py *.go go.*
	go build

python.wasm lib:
	rm -rf python.wasm lib python-3.12.1-wasi_sdk-20.zip
	wget https://github.com/brettcannon/cpython-wasi-build/releases/download/v3.12.1/python-3.12.1-wasi_sdk-20.zip
	unzip python-3.12.1-wasi_sdk-20.zip

wazero:
	wget https://github.com/tetratelabs/wazero/releases/download/v1.6.0/wazero_1.6.0_linux_amd64.tar.gz
	tar xf wazero_1.6.0_linux_amd64.tar.gz wazero

docker-build:
	docker build --progress=plain --tag=python-wazero .

docker-build-no-cache:
	docker build --no-cache --progress=plain .

docker-run: docker-build
	docker run python-wazero

clean:
	rm -rf example python*wasi*.zip python.wasm lib .cache wazero*

install:
	docker exec jobico-control-plane mkdir -p /data/volumes/pv1/python/prg
	docker exec jobico-control-plane curl -o /data/volumes/pv1/python/python-3.13.0a5-wasi_sdk-20.zip -LJ https://github.com/brettcannon/cpython-wasi-build/releases/download/v3.13.0a5/python-3.13.0a5-wasi_sdk-20.zip
	docker exec jobico-control-plane python3 -m zipfile -e /data/volumes/pv1/python/python-3.13.0a5-wasi_sdk-20.zip /data/volumes/pv1/python/
	docker exec jobico-control-plane rm /data/volumes/pv1/python/python-3.13.0a5-wasi_sdk-20.zip
	docker cp sdk/ jobico-control-plane:/data/volumes/pv1/python

install-file:
	docker cp test/hello/main.py jobico-control-plane:/data/volumes/pv1/python/prg
	
install-job:
	kubectl apply -f test/k8s-greet.yml

uninstall-job:
	kubectl delete -f test/k8s-greet.yml
