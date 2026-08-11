//go:build ignore



package main

import (
	"bytes"
	"fmt"
	"go/parser"
	"go/printer"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	fset := token.NewFileSet()
	err := filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if d.Name() == ".git" || d.Name() == "node_modules" {
				return fs.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}
		if filepath.Base(path) == "remove_comments.go" {
			return nil
		}

		
		f, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			return nil 
		}

		var buf bytes.Buffer
		
		err = printer.Fprint(&buf, fset, f)
		if err != nil {
			fmt.Printf("Error printing %s: %v\n", path, err)
			return nil
		}

		
		original, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		
		if !bytes.Equal(original, buf.Bytes()) {
			err = os.WriteFile(path, buf.Bytes(), 0644)
			if err != nil {
				fmt.Printf("Error writing %s: %v\n", path, err)
				return nil
			}
			fmt.Println("Stripped comments from:", path)
		}
		return nil
	})
	if err != nil {
		fmt.Println("Error:", err)
	}
}

