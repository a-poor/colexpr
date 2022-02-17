package colexpr

import (
	"errors"
	"fmt"

	"github.com/alecthomas/participle/v2"
)

var (
	ErrNoOpPassed      = errors.New("no operator passed")
	ErrTooManyOps      = errors.New("too many operators passed")
	ErrLeftNotInt      = errors.New("left value is not an int")
	ErrRightNotInt     = errors.New("right value is not an int")
	ErrUnknownOperator = errors.New("unknown operator")
)

type Parser struct {
	parser *participle.Parser
}

func NewParser() *Parser {
	return &Parser{
		participle.MustBuild(&Expr{}),
	}
}

func (p *Parser) ParseExpression(s string) (*Expr, error) {
	e := &Expr{}
	if err := p.parser.ParseString("", s, e); err != nil {
		return nil, err
	}
	return e, nil
}

type Expr struct {
	Left  *Value    `parse:"@@"`
	Op    *Operator `parse:"@( \"*\" | \"/\" | \"+\" | \"-\" )"`
	Right *Value    `parse:"@@"`
}

func (e *Expr) Evaluate() (interface{}, error) {
	if e.Left.isInt() && e.Right.isInt() {
		return e.EvaluateInt()
	}
	return e.EvaluateFloat(), nil
}

func (e *Expr) EvaluateInt() (int, error) {
	// Check types
	if !e.Left.isInt() {
		return 0, ErrLeftNotInt
	}
	if !e.Right.isInt() {
		return 0, ErrRightNotInt
	}

	a, b := e.Left.getInt(), e.Right.getInt()
	res := e.Op.applyInt(a, b)
	return res, nil
}

func (e *Expr) EvaluateFloat() float64 {
	a, b := e.Left.getFloat(), e.Right.getFloat()
	res := e.Op.applyFloat(a, b)
	return res
}

type Value struct {
	Integer *int     `parse:"  @Int"`
	Float   *float64 `parse:"| @Float"`
}

func (v *Value) isInt() bool {
	return v.Integer != nil
}

func (v *Value) getInt() int {
	if v.Integer == nil {
		panic("not an int")
	}
	return *v.Integer
}

func (v *Value) getFloat() float64 {
	if v.Integer != nil {
		return float64(*v.Integer)
	}
	return *v.Float
}

type Operator int

func (op Operator) applyInt(l, r int) int {
	switch op {
	case OpMul:
		return l * r
	case OpDiv:
		return l / r
	case OpAdd:
		return l + r
	case OpSub:
		return l - r
	}
	panic("unreachable")
}

func (op Operator) applyFloat(l, r float64) float64 {
	switch op {
	case OpMul:
		return l * r
	case OpDiv:
		return l / r
	case OpAdd:
		return l + r
	case OpSub:
		return l - r
	}
	panic("unreachable")
}

const (
	OpMul Operator = iota
	OpDiv
	OpAdd
	OpSub
)

var opMap = map[string]Operator{
	"*": OpMul,
	"/": OpDiv,
	"+": OpAdd,
	"-": OpSub,
}

func (o *Operator) Capture(s []string) error {
	// Validate the input
	if len(s) == 0 {
		return ErrNoOpPassed
	}
	if len(s) > 1 {
		return ErrTooManyOps
	}

	// Pull out the operator
	c := s[0]

	// Lookup the operator
	op, ok := opMap[c]
	if !ok {
		return fmt.Errorf("unknown operator %q: %w", c, ErrUnknownOperator)
	}

	// Success!
	*o = op
	return nil
}
