// Package gen provides support for writing go generate commands
package gen

import (
	"go/ast"
	"text/template"
)

// TemplateType locates the type named ty in the fileset fs and passes its
// definition to template t. If format is true then the output of the template
// will be passed through go fmt.
func TemplateType(ty string, fs FileSet, t *template.Template, format bool) error {
	tm := &templater{}

	for _, af := range fs.AstFiles {
		if af != nil {
			ast.Inspect(af, tm.inspect)
			if tm.err != nil {
				return tm.err
			}
		}
	}

	return nil
}

type templater struct {
	err error
}

func (t *templater) inspect(node ast.Node) bool {
	// typeDecl, ok := node.(*ast.TypeSpec)
	// if !ok || typeDecl.Name.Name != f.typeName {
	// 	// We only care about type declarations.
	// 	return true
	// }
	return false
}
