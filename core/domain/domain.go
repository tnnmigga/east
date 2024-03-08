package domain

type Pool struct {
	useCases []any
}

func NewPool(maxCaseIndex int) *Pool {
	return &Pool{
		useCases: make([]any, maxCaseIndex),
	}
}

func (p *Pool) PutCase(caseIndex int, useCase any) {
	p.useCases[caseIndex] = useCase
}

func (p *Pool) GetCase(caseIndex int) any {
	return p.useCases[caseIndex]
}
