pbc:
	protoc --proto_path=./ --gofast_out=.  pb/*.proto

corepb:
	protoc --proto_path=./ --gofast_out=source_relative:.  core/eventbus/*.proto

ba:
	cd game && go build -o ../bin/game main.go 
	cd gateway && go build -o ../bin/gateway main.go 

rungateway:
	cd gateway && ../bin/gateway

rungame:
	cd game && ../bin/game

runclient:
	node tools/fakecli/main.js