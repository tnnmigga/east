package domain

type IUseCaseInit interface {
	Init()
}

type IUseCaseDestroy interface {
	Destroy()
}

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

func (p *Pool) Init() {
	for _, v := range p.useCases {
		if uc, ok := v.(IUseCaseInit); ok {
			uc.Init()
		}
	}
}

func (p *Pool) Destroy() {
	for i := len(p.useCases) - 1; i >= 0; i-- {
		if uc, ok := p.useCases[i].(IUseCaseDestroy); ok {
			uc.Destroy()
		}
	}
}
