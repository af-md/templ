package generator

import (
	"bytes"
	"testing"

	"github.com/a-h/templ/parser/v2"
	"github.com/google/go-cmp/cmp"
)

func TestGeneratorSourceMap(t *testing.T) {
	w := new(bytes.Buffer)
	g := generator{
		w:         NewRangeWriter(w),
		sourceMap: parser.NewSourceMap(),
	}
	exp := parser.GoExpression{
		Expression: parser.Expression{
			Value: "line1\nline2",
		},
	}
	err := g.writeGoExpression(exp)
	if err != nil {
		t.Fatalf("failed to write Go expression: %v", err)
	}
	expected := parser.NewPosition(0, 0, 0)

	actual, ok := g.sourceMap.TargetPositionFromSource(0, 0)
	if !ok {
		t.Errorf("failed to get matching target")
	}
	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf("unexpected target:\n%v", diff)
	}
}

func TestRewriteExpression(t *testing.T) {
	tests := []struct {
		name     string
		input    parser.Expression
		expected parser.Expression
	}{
		{
			name: "expressions that don't include templ.Context() are not changed",
			input: parser.Expression{
				Value: "hello",
				Range: parser.Range{
					From: parser.Position{Index: 0, Line: 0, Col: 0},
					To:   parser.Position{Index: 5, Line: 0, Col: 5},
				},
			},
			expected: parser.Expression{
				Value: "hello",
				Range: parser.Range{
					From: parser.Position{Index: 0, Line: 0, Col: 0},
					To:   parser.Position{Index: 5, Line: 0, Col: 5},
				},
			},
		},
		{
			name: `templ.Context() is changed to templ_7745c5c3_Ctx`,
			input: parser.Expression{
				Value: "templ.Context()",
				Range: parser.Range{
					From: parser.Position{Index: 0, Line: 0, Col: 0},
					To:   parser.Position{Index: 15, Line: 0, Col: 15},
				},
			},
			expected: parser.Expression{
				Value: "templ_7745c5c3_Ctx",
				Range: parser.Range{
					From: parser.Position{Index: 0, Line: 0, Col: 0},
					To:   parser.Position{Index: 18, Line: 0, Col: 18},
				},
			},
		},
		{
			name: `templ.Context() can be used to extract a value`,
			input: parser.Expression{
				Value: `templ.Context().Value("abc")`,
				Range: parser.Range{
					From: parser.Position{Index: 0, Line: 0, Col: 0},
					To:   parser.Position{Index: 28, Line: 0, Col: 28},
				},
			},
			expected: parser.Expression{
				Value: `templ_7745c5c3_Ctx.Value("abc")`,
				Range: parser.Range{
					From: parser.Position{Index: 0, Line: 0, Col: 0},
					To:   parser.Position{Index: 31, Line: 0, Col: 31},
				},
			},
		},
		{
			name: `templ.Context() is not replaced inside standard strings`,
			input: parser.Expression{
				Value: `"templ.Context()"`,
				Range: parser.Range{
					From: parser.Position{Index: 0, Line: 0, Col: 0},
					To:   parser.Position{Index: 17, Line: 0, Col: 17},
				},
			},
			expected: parser.Expression{
				Value: `"templ.Context()"`,
				Range: parser.Range{
					From: parser.Position{Index: 0, Line: 0, Col: 0},
					To:   parser.Position{Index: 17, Line: 0, Col: 17},
				},
			},
		},
		{
			name: `templ.Context() is not replaced inside backtick strings`,
			input: parser.Expression{
				Value: "`templ.Context()`",
				Range: parser.Range{
					From: parser.Position{Index: 0, Line: 0, Col: 0},
					To:   parser.Position{Index: 17, Line: 0, Col: 17},
				},
			},
			expected: parser.Expression{
				Value: "`templ.Context()`",
				Range: parser.Range{
					From: parser.Position{Index: 0, Line: 0, Col: 0},
					To:   parser.Position{Index: 17, Line: 0, Col: 17},
				},
			},
		},
		{
			name: `templ.Context() doesn't need to be at the start of the expression`,
			input: parser.Expression{
				Value: ` value := templ.Context().Value("abc").(string) `,
				Range: parser.Range{
					From: parser.Position{Index: 0, Line: 0, Col: 0},
					To:   parser.Position{Index: 17, Line: 0, Col: 17},
				},
			},
			expected: parser.Expression{
				Value: ` value := templ_7745c5c3_Ctx.Value("abc").(string) `,
				Range: parser.Range{
					From: parser.Position{Index: 0, Line: 0, Col: 0},
					To:   parser.Position{Index: 51, Line: 0, Col: 51},
				},
			},
		},
		{
			name: `expressions can be multiline`,
			input: parser.Expression{
				Value: `if true {
	value = templ.Context().Value("abc").(string)
}`,
				Range: parser.Range{
					From: parser.Position{Index: 0, Line: 0, Col: 0},
					To:   parser.Position{Index: 58, Line: 2, Col: 1},
				},
			},
			expected: parser.Expression{
				Value: `if true {
	value = templ_7745c5c3_Ctx.Value("abc").(string)
}`,
				Range: parser.Range{
					From: parser.Position{Index: 0, Line: 0, Col: 0},
					To:   parser.Position{Index: 61, Line: 2, Col: 1},
				},
			},
		},
		{
			name: `expressions with invalid standard strings are returned as-is`,
			input: parser.Expression{
				Value: `templ.Context().Value("ab`,
				Range: parser.Range{
					From: parser.Position{Index: 0, Line: 0, Col: 0},
					To:   parser.Position{Index: 25, Line: 0, Col: 25},
				},
			},
			expected: parser.Expression{
				Value: `templ.Context().Value("ab`,
				Range: parser.Range{
					From: parser.Position{Index: 0, Line: 0, Col: 0},
					To:   parser.Position{Index: 25, Line: 0, Col: 25},
				},
			},
		},
		{
			name: `expressions with invalid backtick strings are returned as-is`,
			input: parser.Expression{
				Value: "templ.Context().Value(`ab",
				Range: parser.Range{
					From: parser.Position{Index: 0, Line: 0, Col: 0},
					To:   parser.Position{Index: 25, Line: 0, Col: 25},
				},
			},
			expected: parser.Expression{
				Value: "templ.Context().Value(`ab",
				Range: parser.Range{
					From: parser.Position{Index: 0, Line: 0, Col: 0},
					To:   parser.Position{Index: 25, Line: 0, Col: 25},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := rewriteExpression(test.input)
			if diff := cmp.Diff(test.expected, actual); diff != "" {
				t.Error(diff)
			}
		})
	}
}
