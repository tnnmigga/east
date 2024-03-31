pbc:
	@python tools/generator/pbc.py source=$(source) include=vendor/github.com/gogo/protobuf

corepb:
	protoc --proto_path=./ --gofast_out=source_relative:.  core/eventbus/*.proto
	protoc --proto_path=./ --gofast_out=source_relative:.  core/infra/link/*.proto

ba:
	cd game && go build -o ../bin/game main.go 
	cd gateway && go build -o ../bin/gateway main.go 
	cd login && go build -o ../bin/login main.go

run:
	tools/shell/run.sh
	
stop:
	tools/shell/stop.sh

client:
	node tools/fakecli/main.js

br: ba run

rerun: stop br

worker:
	@python tools/generator/worker.py name=$(name)