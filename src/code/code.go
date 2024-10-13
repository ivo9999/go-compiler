package code

type Instructions []byte

type Opcode byte

const (
	OpConstant Opcode = iota
)

type Definition struct {
	Name         string
	OprandWidths []int
}

var definitions = map[Opcode]Definition{
	OpConstant: {"OpConstant", []int{2}},
}

func Lookup(op byte) Definition {
	def, ok := definitions[Opcode(op)]
	if !ok {
		return Definition{}
	}
	return def
}
