package domain

type Pool struct {
	services []any
}

func NewPool(maxCaseIndex int) *Pool {
	return &Pool{
		services: make([]any, maxCaseIndex),
	}
}

func (p *Pool) PutCase(caseIndex int, useCase any) {
	p.services[caseIndex] = useCase
}

func (p *Pool) GetCase(caseIndex int) any {
	return p.services[caseIndex]
}
