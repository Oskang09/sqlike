package rsql

type tokenType int

// types :
const (
	Operator = iota
	String
	Or
	And
	Numeric
	Text
	Group
	Whitespace
)

// Token L
type Token struct {
	Type        int
	Value       string
	Lexeme      []byte
	TC          int
	StartLine   int
	StartColumn int
	EndLine     int
	EndColumn   int
}