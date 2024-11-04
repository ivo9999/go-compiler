package compiler

type SymbolScope string

const (
	GlobalScope  SymbolScope = "GLOBAL"
	LocalScope   SymbolScope = "LOCAL"
	BuiltinScope SymbolScope = "BUILTIN"
)

type Symbol struct {
	Name  string
	Scope SymbolScope
	Index int
}

type SymbolTable struct {
	numDefinitions int
	store          map[string]Symbol
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{store: make(map[string]Symbol)}
}

func (s *SymbolTable) Define(name string) Symbol {
	symbol := Symbol{
		Name:  name,
		Scope: GlobalScope,
		Index: s.numDefinitions,
	}
	s.store[name] = symbol
	s.numDefinitions++
	return symbol
}

func (s *SymbolTable) Resolve(name string) (Symbol, bool) {
	symbol, ok := s.store[name]
	if !ok {
		return Symbol{}, false
	}
	return symbol, true
}
