package domain

type IUseCaseInit interface {
	Init()
}

type IUseCaseDestroy interface {
	Destroy()
}

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

func (p *Pool) Init() {
	for _, v := range p.services {
		if uc, ok := v.(IUseCaseInit); ok {
			uc.Init()
		}
	}
}

func (p *Pool) Destroy() {
	for i := len(p.services) - 1; i >= 0; i-- {
		if uc, ok := p.services[i].(IUseCaseDestroy); ok {
			uc.Destroy()
		}
	}
}
