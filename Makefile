pbc:
	@python nett/scripts/generator/pbc.py source=$(source)

corepb:
	protoc --proto_path=./ --gofast_out=source_relative:.  core/eventbus/*.proto
	protoc --proto_path=./ --gofast_out=source_relative:.  core/infra/link/*.proto

ba:
	cd game && go build -o ../bin/game main.go 
	cd gateway && go build -o ../bin/gateway main.go 
	cd door && go build -o ../bin/door main.go

run:
	tools/shell/run.sh
	
stop:
	tools/shell/stop.sh

wscli:
	@node tools/fakecli/wscli.js

tcpcli:
	@node tools/fakecli/tcpcli.js

br: ba run

rerun: stop br

worker:
	@python tools/generator/worker.py name=$(name)

init:
	git clone git@github.com:tnnmigga/nett.git
	cd nett/scripts/fakecli/ && npm install
	go mod vendor

.PHONY: pbc corepb ba run stop wscli tcpcli br rerun worker init