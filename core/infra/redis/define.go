package redis

type Pipeline struct {

}

type op interface {
	key() string
	value() any
}

type OpGet struct {
	Key string
	Value any
}

func (op *OpGet) key() string {
	return op.Key
}

func (op *OpGet) value() any {
	return op.Value
}

type OpSet struct {
	Key   string
	Value any
}

func (op *OpSet) key() string {
	return op.Key
}

func (op *OpSet) value() any {
	return op.Value
}

type Del struct {
	Key string
}

func (op *Del) key() string {
	return op.Key
}

func (op *Del) value() any {
	return nil
}