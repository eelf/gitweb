
.PHONY: client server

server/proto client/src/proto: proto/gitweb.proto
	protoc \
	--gogo_out=plugins=grpc:./server \
	--plugin=protoc-gen-ts=./client/node_modules/.bin/protoc-gen-ts \
	--ts_out=service=true:./client/src \
	--js_out=import_style=commonjs,binary:./client/src \
	$^

clean_proto:
	rm -rf server/proto app/src/proto

client:
	cd client && yarn
	cd client && node_modules/webpack/bin/webpack.js watch

server:
	cd server && go run cmd_server/main.go ../client/build
