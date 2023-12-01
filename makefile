pbc:
	protoc --proto_path=./ --gofast_out=.  pb/*.proto

corepb:
	protoc --proto_path=./ --gofast_out=source_relative:.  core/pb/*.proto

ba:
	cd game && go build -o ../bin/game main.go 