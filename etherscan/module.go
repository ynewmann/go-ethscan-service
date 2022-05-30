package etherscan

type (
	Moduler interface {
		GetName() string
	}

	Module struct {
		api  *api
		Name string
	}
)

func NewModule(api *api, name string) *Module {
	return &Module{
		api,
		name,
	}
}

func (m *Module) GetName() string {
	return m.Name
}
