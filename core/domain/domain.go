package domain

type Pool struct {
	services []any
}

func NewPool(maxCaseIndex int) *Pool {
	return &Pool{
		services: make([]any, maxCaseIndex),
	}
}

func (p *Pool) PutImpl(caseIndex int, useCase any) {
	p.services[caseIndex] = useCase
}

func (p *Pool) GetImpl(caseIndex int) any {
	return p.services[caseIndex]
}
