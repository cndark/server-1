package mdata

// ============================================================================

var (
	m_registry = make(map[string]*Reg) // [module name]
)

// ============================================================================

type Reg struct {
	Name string

	LoadData   func() interface{}
	DataLoaded func()

	SaveAsync func()
	Save      func()

	NewModuleData func() interface{}
}

// ============================================================================

func Register(reg *Reg) {
	m_registry[reg.Name] = reg
}
