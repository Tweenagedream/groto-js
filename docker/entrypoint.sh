#!/bin/sh

mkdir -p /protos/js /protos/static
echo "Closure Protos in protos/js, compiled protos in protos/static"
/root/bin/protoc --js_out=library=$1,binary:/protos/js/ \
	--proto_path /protos/ \
	/protos/*.proto
echo "protos/js completed"
entrypoint=$(grep proto /protos/js/*.js | head -n1 | cut -d"'" -f2)
java -jar closure*.jar \
	--entry_point="$entrypoint" \
	--js_output_file /protos/static/compiled.js \
	--dependency_mode STRICT \
	--env BROWSER \
	protobuf*/js/ \
	closure-library/ \
	/protos/js/*.js
echo "proto/static completed"