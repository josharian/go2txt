package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"log"
	"os"

	"github.com/josharian/go2txt/edit"
)

func main() {
	// Read stdin into buf and simultaneously parse into an AST
	fset := token.NewFileSet()
	buf := new(bytes.Buffer)
	r := io.TeeReader(os.Stdin, buf)
	f, err := parser.ParseFile(fset, "", r, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}
	b := buf.Bytes()

	// Walk the AST, looking for []T{...} and rewriting into slice T{...}.
	// This uses a copy of GOROOT/src/cmd/go/internal/edit.
	ed := edit.NewBuffer(b)
	ast.Inspect(f, func(n ast.Node) bool {
		at, ok := n.(*ast.ArrayType)
		if !ok {
			return true
		}
		if at.Len != nil { // ignore arrays like [10] or [...]; just slices for now
			return true
		}
		lb := fset.Position(at.Lbrack).Offset
		rb := lb + bytes.IndexByte(b[lb:], ']') + 1 // ] is guaranteed to exist, otherwise there would have been a parse failure
		ed.Replace(lb, rb, "slice ")
		return true
	})

	// Dump the result to stdout
	fmt.Println(ed.String())
}
