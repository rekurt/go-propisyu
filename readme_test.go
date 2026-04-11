package propisyu

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
	"testing"
	"unicode"
)

// TestReadmeCoversPublicAPI walks the package AST and asserts that every
// exported top-level declaration (func, type, const, var) is mentioned in
// both README.md and README_EN.md. This catches documentation drift the
// moment a new public API is added, so the READMEs can never get out of
// sync the way they did before v0.4.0 — where Ordinal, Money,
// DecimalToWordsPrecision, and the three Currency presets shipped in the
// b0f007a commit without any README mention and stayed undocumented for
// several releases.
//
// The check uses go/ast rather than regex so it correctly handles block
// forms like `var ( CurrencyRUB = ... )` and `const ( GenderMasculine = 1 )`.
func TestReadmeCoversPublicAPI(t *testing.T) {
	t.Parallel()

	names := collectExportedNames(t)
	if len(names) == 0 {
		t.Fatal("no exported names found — parser is misconfigured")
	}

	ru := readReadme(t, "README.md")
	en := readReadme(t, "README_EN.md")

	var missingRu, missingEn []string
	for _, name := range names {
		if !strings.Contains(ru, name) {
			missingRu = append(missingRu, name)
		}
		if !strings.Contains(en, name) {
			missingEn = append(missingEn, name)
		}
	}

	if len(missingRu) > 0 {
		t.Errorf(
			"README.md is missing %d public API name(s): %v",
			len(missingRu), missingRu,
		)
	}
	if len(missingEn) > 0 {
		t.Errorf(
			"README_EN.md is missing %d public API name(s): %v",
			len(missingEn), missingEn,
		)
	}
}

// TestReadmesInHeaderParity asserts that README.md and README_EN.md have
// the same number of level-2 and level-3 markdown headers. Translations
// drift silently when one language gets a new section and the other
// doesn't — this test catches that as soon as it happens.
func TestReadmesInHeaderParity(t *testing.T) {
	t.Parallel()

	ru := readReadme(t, "README.md")
	en := readReadme(t, "README_EN.md")

	levels := []struct {
		label  string
		prefix string
	}{
		{"## (h2)", "## "},
		{"### (h3)", "### "},
	}

	for _, level := range levels {
		ruN := countLinesWithPrefix(ru, level.prefix)
		enN := countLinesWithPrefix(en, level.prefix)
		if ruN != enN {
			t.Errorf(
				"README header parity %s: README.md has %d, README_EN.md has %d",
				level.label, ruN, enN,
			)
		}
	}
}

func collectExportedNames(t *testing.T) []string {
	t.Helper()

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, ".", func(info os.FileInfo) bool {
		return !strings.HasSuffix(info.Name(), "_test.go")
	}, 0)
	if err != nil {
		t.Fatalf("parse package: %v", err)
	}

	var names []string
	seen := map[string]bool{}
	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			for _, decl := range file.Decls {
				collectFromDecl(decl, &names, seen)
			}
		}
	}
	return names
}

func collectFromDecl(decl ast.Decl, dst *[]string, seen map[string]bool) {
	switch d := decl.(type) {
	case *ast.FuncDecl:
		if d.Recv != nil {
			return // skip methods
		}
		addIfExported(d.Name.Name, dst, seen)
	case *ast.GenDecl:
		for _, spec := range d.Specs {
			switch s := spec.(type) {
			case *ast.TypeSpec:
				addIfExported(s.Name.Name, dst, seen)
			case *ast.ValueSpec:
				for _, n := range s.Names {
					addIfExported(n.Name, dst, seen)
				}
			}
		}
	}
}

func addIfExported(name string, dst *[]string, seen map[string]bool) {
	if name == "" || seen[name] {
		return
	}
	if !unicode.IsUpper([]rune(name)[0]) {
		return
	}
	seen[name] = true
	*dst = append(*dst, name)
}

func readReadme(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(b)
}

func countLinesWithPrefix(s, prefix string) int {
	n := 0
	for _, line := range strings.Split(s, "\n") {
		if strings.HasPrefix(line, prefix) {
			n++
		}
	}
	return n
}
