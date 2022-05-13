package etherscan

type (
	Moduler interface {
		GetName() string
	}

	Module struct {
		client *apiClient
		Name   string
	}
)

func NewModule(client *apiClient, name string) *Module {
	return &Module{
		client,
		name,
	}
}

func (m *Module) GetName() string {
	return m.Name
}
