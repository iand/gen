package gen

import (
	"go/ast"
	"go/build"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
)

// FileSet is a parsed set of Go source files which are assumed to form a package.
type FileSet struct {
	// Dir holds the name of the directory
	Dir string

	// Files contains the absolute path and name of each file in the FileSet.
	Files []string

	// FileSet holds the positions of each token in the set of parsed Go source files.
	FileSet *token.FileSet

	// AstFiles are the parsed versions of each source file.
	AstFiles []*ast.File

	// TypeInfo holds result type information for the source files.
	TypeInfo *types.Info

	// Package holds information about the package formed from the files in the FileSet.
	Package *types.Package
}

const currentDir = "."

// NewFileSet creates a FileSet from a (possibly empty) list of names. If
// multiple names are provided they are assumed to refer to Go source files.
// If a single name is provided that matches a directory then the fileset will
// be initialised to contain the Go source files in that directory. If no
// names are provided then the current working directory is assumed.
func NewFileSet(names []string) (*FileSet, error) {
	// No names supplied so assume current directory
	if len(names) == 0 {
		return FileSetFromDir(currentDir)
	}

	// One name supplied could be a directory or a single file
	// Find out which
	if len(names) == 1 {
		info, err := os.Stat(names[0])
		if err != nil {
			return nil, err
		}
		if info.IsDir() {
			return FileSetFromDir(names[0])
		}
	}

	// Assume names are files
	fs := &FileSet{
		Dir:   filepath.Dir(names[0]),
		Files: names,
	}

	return fs.Parse()

}

// FileSetFromDir creates a FileSet consisting of the Go source files
// in the directory d
func FileSetFromDir(d string) (*FileSet, error) {
	fs := FileSet{
		Dir: d,
	}
	pkg, err := build.Default.ImportDir(d, 0)
	if err != nil {
		return nil, err
	}

	fs.Files = append(fs.Files, pkg.GoFiles...)
	if d == currentDir {
		return fs.Parse()
	}

	for i, f := range fs.Files {
		fs.Files[i] = filepath.Join(d, f)
	}

	return fs.Parse()
}

// Parse verifies whether fs represents a valid, compilable set of Go
// source files and sets the parsed versions of each file in the fileset.
func (fs *FileSet) Parse() (*FileSet, error) {
	var err error

	fs.FileSet = token.NewFileSet()
	for _, f := range fs.Files {
		p, err := parser.ParseFile(fs.FileSet, f, nil, 0)
		if err != nil {
			return nil, err
		}
		fs.AstFiles = append(fs.AstFiles, p)
	}

	config := types.Config{Importer: importer.Default()}
	fs.TypeInfo = &types.Info{
		Types:      make(map[ast.Expr]types.TypeAndValue),
		Defs:       make(map[*ast.Ident]types.Object),
		Uses:       make(map[*ast.Ident]types.Object),
		Implicits:  make(map[ast.Node]types.Object),
		Selections: make(map[*ast.SelectorExpr]*types.Selection),
		Scopes:     make(map[ast.Node]*types.Scope),
	}

	fs.Package, err = config.Check(fs.Dir, fs.FileSet, fs.AstFiles, fs.TypeInfo)
	if err != nil {
		return nil, err
	}

	return fs, nil
}

// Walk traverses all the files in fs invoking v.Visit on each file in turn.
func (fs *FileSet) Walk(v ast.Visitor) {
	for _, astFile := range fs.AstFiles {
		ast.Walk(v, astFile)
	}
}

// Inspect traverses all the files in fs calling f on each file in turn. If f
// returns true Inspect invokes f recursively for each of the non-nil
// children of the file's root node, followed by a call of f(nil).
func (fs *FileSet) Inspect(f func(ast.Node) bool) {
	for _, astFile := range fs.AstFiles {
		ast.Inspect(astFile, f)
	}
}
