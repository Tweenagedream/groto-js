

Build the image in the Docker directory and run docker from the main directory:

docker build docker/. -t proto-js-compiler
docker run -v protos:/protos count_pb
go run main.go