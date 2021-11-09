package main

import (
	"context"
	"fmt"
	"os"

	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
	"github.com/open-policy-agent/opa/types"
)

func main() {

	module := `
package example.imageimport

default allow = false

allow {
    is_changed
}

allow {
    is_new_version
}

allow {
    count(input.last_10_tags) < 4
    hello("Luis") = "hello, Luis"
}

is_changed {
    input.old_digest != input.new_digest
}

is_new_version {
	input.major_version < 3
	input.minor_version < 2
}
`

	ctx := context.Background()

	query, err := rego.New(
		rego.Query(`
x = data.example.imageimport.allow
`),

		rego.Function1(
			&rego.Function{
				Name: "hello",
				Decl: types.NewFunction(types.Args(types.S), types.S),
			},
			func(_ rego.BuiltinContext, a *ast.Term) (*ast.Term, error) {
				if str, ok := a.Value.(ast.String); ok {
					return ast.StringTerm("hello, " + string(str)), nil
				}
				return nil, nil
			}),

		rego.Module("example.rego", module),
	).PrepareForEval(ctx)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Buu: %v", err)
		os.Exit(1)
	}

	input := map[string]interface{}{
		"event":         "image_pushed",
		"name":          "my_shiny_new_image",
		"tag":           "2.1.001",
		"old_digest":    "123",
		"new_digest":    "123",
		"major_version": 2,
		"minor_version": 2,
		"last_10_tags":  []string{"2.1.001", "2.0.005", "2.0.003"},
	}

	results, err := query.Eval(ctx, rego.EvalInput(input))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Buu: %v", err)
		os.Exit(1)
	}

	fmt.Printf("Results: %v\n", results)

	fmt.Println("Import allowed:", results[0].Bindings["x"])
}
