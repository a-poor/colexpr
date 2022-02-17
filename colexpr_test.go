package colexpr_test

import (
	"fmt"
	"testing"

	"github.com/a-poor/colexpr"
	"github.com/alecthomas/repr"
)

func TestNothing(t *testing.T) {
	t.Log("Nothing is working...which is good.")
}

func TestParse(t *testing.T) {
	// Create the test cases
	testCasesInt := []struct {
		name  string
		expr  string
		left  int
		right int
		res   int
		op    colexpr.Operator
		err   error
	}{
		{"one-plus-one", "1 + 2", 1, 2, 3, colexpr.OpAdd, nil},
		{"one-plus-three", "2 + 3", 2, 3, 5, colexpr.OpAdd, nil},
		{"four-minus-three", "4 - 3", 4, 3, 1, colexpr.OpSub, nil},
		{"one-minus-three", "1 - 3", 1, 3, -2, colexpr.OpSub, nil},
	}

	// Run the test cases
	for _, tc := range testCasesInt {
		tname := fmt.Sprintf("Int=%s", tc.name)
		t.Run(tname, func(t *testing.T) {
			// Run these tests in parallel
			t.Parallel()

			// Create a parser (this shouldn't fail, so...)
			p := colexpr.NewParser()

			// Parse the expression
			e, err := p.ParseExpression(tc.expr)
			if err != nil { // TODO: Check for expected errors...
				t.Errorf("Error parsing expression: %v", err)
			}

			// Check the parsed results
			if *e.Left.Integer != tc.left {
				t.Errorf("Expected left value to be %d, got %s", tc.left, repr.String(e.Left))
			}
			if *e.Op != tc.op {
				t.Errorf("Expected operator to be %d, got %d", tc.right, e.Op)
			}
			if *e.Right.Integer != tc.right {
				t.Errorf("Expected right value to be %d, got %s", tc.right, repr.String(e.Right))
			}

			// Evaluate the expression
			res, err := e.Evaluate()
			if err != nil {
				t.Errorf("Error evaluating expression: %v", err)
			}

			// Check the result (type and value)
			if tres, ok := res.(int); !ok || tres != tc.res {
				t.Errorf("Expected result to be %d, got %s", tc.res, repr.String(res))
			}
		})

	}
}
