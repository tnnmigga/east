pbc:
	protoc --proto_path=./ --gofast_out=.  pb/*.proto

corepb:
	protoc --proto_path=./ --gofast_out=.  core/pb/*.proto
