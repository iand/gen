package gen

import (
	"go/ast"
	"reflect"
	"sort"
	"testing"
)

func TestEachType(t *testing.T) {
	testCases := []struct {
		src   string
		names []string
		err   error
	}{

		{
			src: `package p
					type X struct{
						a string
					}`,
			names: []string{"X"},
		},

		{
			src: `package p
					type Y interface{
						a()
					}`,
			names: []string{"Y"},
		},

		{
			src: `package p
					type (
						X interface{
							a()
						}
						Y struct {
							a string
						}
						Z int
					)`,
			names: []string{"X", "Y", "Z"},
		},

		{
			src: `package p
					var x = 1`,
			names: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			fs, err := NewFileSetFromTexts(tc.src)
			if !reflect.DeepEqual(err, tc.err) {
				t.Fatalf("got error %v, wanted %v", err, tc.err)
			}

			if err != nil {
				return
			}

			names := []string{}
			fs.EachType(func(ts *ast.TypeSpec) bool {
				names = append(names, ts.Name.Name)
				return true
			})

			sort.Strings(names)
			sort.Strings(tc.names)
			if !reflect.DeepEqual(names, tc.names) {
				t.Errorf("got %+v, wanted %+v", names, tc.names)
			}
		})
	}
}

func TestEachConst(t *testing.T) {
	testCases := []struct {
		src   string
		names []string
		err   error
	}{

		{
			src: `package p
					const X=1`,
			names: []string{"X"},
		},

		{
			src: `package p
					const (
						Y = "yes"
					)`,
			names: []string{"Y"},
		},

		{
			src: `package p
					const (
						X = iota
						Y
						Z
					)`,
			names: []string{"X", "Y", "Z"},
		},

		{
			src: `package p
					var x = 1`,
			names: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			fs, err := NewFileSetFromTexts(tc.src)
			if !reflect.DeepEqual(err, tc.err) {
				t.Fatalf("got error %v, wanted %v", err, tc.err)
			}

			if err != nil {
				return
			}

			names := []string{}
			fs.EachConst(func(vs *ast.ValueSpec) bool {
				names = append(names, vs.Names[0].Name)
				return true
			})

			sort.Strings(names)
			sort.Strings(tc.names)
			if !reflect.DeepEqual(names, tc.names) {
				t.Errorf("got %+v, wanted %+v", names, tc.names)
			}
		})
	}
}

func TestEachFunc(t *testing.T) {
	testCases := []struct {
		src   string
		names []string
		err   error
	}{

		{
			src: `package p
					func X() { }`,
			names: []string{"X"},
		},

		{
			src: `package p
					func X(a string) { }`,
			names: []string{"X"},
		},

		{
			src: `package p
					func X(string) bool{ return false }`,
			names: []string{"X"},
		},

		{
			src: `package p
			      type A struct{}
				  func (a A) X(string) bool{ return false }`,
			names: []string{"X"},
		},
	}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			fs, err := NewFileSetFromTexts(tc.src)
			if !reflect.DeepEqual(err, tc.err) {
				t.Fatalf("got error %v, wanted %v", err, tc.err)
			}

			if err != nil {
				return
			}

			names := []string{}
			fs.EachFunc(func(fd *ast.FuncDecl) bool {
				names = append(names, fd.Name.Name)
				return true
			})

			sort.Strings(names)
			sort.Strings(tc.names)
			if !reflect.DeepEqual(names, tc.names) {
				t.Errorf("got %+v, wanted %+v", names, tc.names)
			}
		})
	}
}
