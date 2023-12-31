pbc:
	protoc --proto_path=./ --gofast_out=.  pb/*.proto

corepb:
	protoc --proto_path=./ --gofast_out=source_relative:.  core/eventbus/*.proto
	protoc --proto_path=./ --gofast_out=source_relative:.  core/infra/nats/*.proto

ba:
	cd game && go build -o ../bin/game main.go 
	cd gateway && go build -o ../bin/gateway main.go 

run:
	tools/shell/run.sh
	
stop:
	tools/shell/stop.sh

client:
	node tools/fakecli/main.js

br: ba run

rerun: stop br