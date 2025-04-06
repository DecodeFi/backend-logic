BUILD_PATH='./build/'
CMD='./cmd'

build_all: build_server build_chisel

build_server:
	go build -o ${BUILD_PATH} ${CMD}/server/...

build_chisel:
	go build -o ${BUILD_PATH} ${CMD}/chisel/...
