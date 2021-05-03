package strace

import (
	"go/ast"
	"go/token"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"unicode"
)

type Strace struct {
	FilePath string
	FileLine int
	FuncName string
}

func (s Strace) Filename() string {
	return filepath.Base(s.FilePath)
}

func (s Strace) Tabs() int {
	file, err := ioutil.ReadFile(s.FilePath)
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(file), "\n")
	line := strings.TrimRightFunc(lines[s.FileLine-1], unicode.IsSpace)
	lenWithoutTabs := len(strings.TrimSpace(line))

	return len(line) - lenWithoutTabs
}

//
// func (s Strace) Replace(t golden.TestingTB, value interface{}) {
// 	openFile, err := lockedfile.OpenFile(s.FilePath, os.O_RDWR, 0666)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer func() {
// 		if err := openFile.Close(); err != nil {
// 			t.Logf(err.Error())
// 		}
// 	}()
//
// 	if err := openFile.SetDeadline(time.Now().Add(time.Second)); err != nil {
// 		t.Logf(err.Error())
// 	}
//
// 	fileSet := token.NewFileSet()
// 	node, err := parser.ParseFile(fileSet, s.FilePath, openFile, parser.AllErrors)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
//
// 	// if err := ast.Print(fileSet, node); err != nil {
// 	// 	t.Fatal(err)
// 	// }
//
// 	res := findNode(fileSet, node, s.FuncName, s.FilePath, s.FileLine)
// 	if res == nil {
// 		t.Fatal("couldn't find the right node to replace the value")
// 	}
//
// 	if res.Value == value.(string) {
// 		return
// 	}
// 	res.Value = value.(string)
//
// 	if err := openFile.Truncate(0); err != nil {
// 		t.Fatal(err)
// 	}
//
// 	if _, err := openFile.Seek(0, io.SeekStart); err != nil {
// 		t.Fatal(err)
// 	}
//
// 	if err := printer.Fprint(openFile, fileSet, node); err != nil {
// 		t.Fatal(err)
// 	}
//
// 	if err := openFile.Close(); err != nil {
// 		t.Fatal(err)
// 	}
// }

func findNode(
	f *token.FileSet, node ast.Node, name string, file string, line int,
) (res *ast.BasicLit) {
	ast.Inspect(
		node, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.CallExpr:
				if isFunc(f, x.Fun, name, file, line) {
					if len(x.Args) != 1 {
						return true
					}
					l, ok := x.Args[0].(*ast.BasicLit)
					if !ok {
						return true
					}

					res = l
					return false
				}
			}

			return true
		},
	)

	return res
}

func isFunc(f *token.FileSet, node ast.Node, name string, file string, line int) bool {
	// log.Println(name, file, line, node)

	selectorExpr, ok := node.(*ast.Ident)
	if !ok {
		return false
	}

	if selectorExpr.Name != name {
		return false
	}

	position := f.Position(selectorExpr.NamePos)
	return position.Filename == file && position.Line == line
}

func NewStrace(skip int) Strace {
	_, filepath, fileline, ok := runtime.Caller(skip)
	if !ok {
		panic("Couldn't get the caller information")
	}

	return Strace{
		FilePath: filepath,
		FileLine: fileline,
		FuncName: funcname(skip),
	}
}

func funcname(skip int) string {
	pc, _, _, ok := runtime.Caller(skip)
	if !ok {
		panic("Couldn't get the caller information")
	}
	forPC := runtime.FuncForPC(pc)
	functionPath := forPC.Name()
	// Next four lines are required to use GCCGO function naming conventions.
	// For Ex:  github_com_docker_libkv_store_mock.WatchTree.pN39_github_com_docker_libkv_store_mock.Mock
	// uses interface information unlike golang github.com/docker/libkv/store/mock.(*Mock).WatchTree
	// With GCCGO we need to remove interface information starting from pN<dd>.
	re := regexp.MustCompile("\\.pN\\d+_")
	if re.MatchString(functionPath) {
		functionPath = re.Split(functionPath, -1)[0]
	}

	parts := strings.Split(functionPath, ".")
	funcname := parts[len(parts)-1]

	return funcname
}
