pbc:
	protoc --proto_path=./ --go_out=.  pb/*.proto

corepb:
	protoc --proto_path=./ --gofast_out=source_relative:.  core/pb/*.proto
