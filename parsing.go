package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

type FailingTest struct {
	TestName     string
	FullTestName string
	TestError    string
}

// GetFailingTests takes the output of `go test ./... -run ^Test_divide$` and returns a list of files that contain failing tests
// Example input:
//--- FAIL: Test_divide (0.00s)
//--- FAIL: Test_divide/1_/_2_=_0 (0.00s)
//example_test.go:55: divide() = -2, want 0
//FAIL
//FAIL    github.com/vivalapanda/gofix/example    0.202s
//FAIL

func GetFailingTests(testOutputJSONL string) (map[string][]FailingTest, error) {
	filesWithFailingTests := make(map[string][]FailingTest)

	scanner := bufio.NewScanner(bytes.NewBufferString(testOutputJSONL))
	testOutputs := make([]map[string]string, 0)

	for scanner.Scan() {
		// Unmarshall, but if a field can't be unmarshalled to a string, just ignore it
		var output map[string]interface{}
		if err := json.Unmarshal(scanner.Bytes(), &output); err != nil {
			return nil, fmt.Errorf("error unmarshalling JSONL input: %v", err)
		}

		outputString := make(map[string]string)
		for key, value := range output {
			if valueString, ok := value.(string); ok {
				outputString[key] = valueString
			}
		}

		testOutputs = append(testOutputs, outputString)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading JSONL input: %v", err)
	}

	for _, output := range testOutputs {
		if output["Action"] == "fail" && output["Test"] != "" {
			testName := output["Test"]
			// packageName := output["Package"]
			// packagePath := strings.TrimPrefix(packageName, "github.com/vivalapanda/gofix/")
			filePath := ""

			// Find the related output with the file path
			for _, outputWithFile := range testOutputs {
				if outputWithFile["Test"] == testName && strings.Contains(outputWithFile["Output"], ".go:") {
					filePath = strings.Split(outputWithFile["Output"], ".go:")[0] + ".go"
					filePath = strings.Trim(filePath, " ")
					break
				}
			}

			if filePath == "" {
				continue
			}

			// filePath is not relative to the current directory, so we need to search our directory and all subdirectories
			// for filePath and replace filePath with the relative path
			err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				if info.Name() == filePath {
					filePath = path
					return nil
				}

				return nil
			})
			if err != nil {
				return nil, fmt.Errorf("error walking directory: %v", err)
			}

			// Find the test error
			testError := ""
			for _, outputWithError := range testOutputs {
				if outputWithError["Test"] == testName && strings.Contains(outputWithError["Output"], ": ") {
					testError = strings.Split(outputWithError["Output"], ": ")[1]
					break
				}
			}

			// TestName should only include the stuff before the first slash
			// however, we want to keep the full test name as a separate field
			fullTestName := testName
			if strings.Contains(testName, "/") {
				testName = strings.Split(testName, "/")[0]
			}

			failingTest := FailingTest{
				TestName:     testName,
				FullTestName: fullTestName,
				TestError:    testError,
			}

			if _, ok := filesWithFailingTests[filePath]; ok {
				filesWithFailingTests[filePath] = append(filesWithFailingTests[filePath], failingTest)
			} else {
				filesWithFailingTests[filePath] = []FailingTest{failingTest}
			}
		}
	}

	return filesWithFailingTests, nil
}

// GetFileFromPath gets *ast.File from the file paths
func GetFileFromPath(filePath string) (*ast.File, error) {
	fset := token.NewFileSet()

	file, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("error parsing file: %v", err)
	}

	return file, nil
}

// GetTestCode takes a list of files that contain failing tests and returns the code of the failing tests
func GetTestCode(file *ast.File, failedTestFuncName string) []*ast.FuncDecl {
	var testCode []*ast.FuncDecl
	for _, decl := range file.Decls {
		switch decl.(type) {
		case *ast.FuncDecl:
			fd := decl.(*ast.FuncDecl)
			if strings.HasPrefix(fd.Name.Name, "Test") && fd.Name.Name == failedTestFuncName {
				testCode = append(testCode, fd)
			}
		}
	}
	return testCode
}

// GetFunctionDependencies takes a list of function declarations and returns a list of function names that the function declarations depend on
func GetFunctionDependencies(funcDecls []*ast.FuncDecl) []string {
	var funcNames []string
	for _, decl := range funcDecls {
		ast.Inspect(decl, func(n ast.Node) bool {
			if callExpr, ok := n.(*ast.CallExpr); ok {
				if ident, ok := callExpr.Fun.(*ast.Ident); ok {
					funcNames = append(funcNames, ident.Name)
				}
			}
			return true
		})
	}
	return funcNames
}

// GetFunctionCode takes a list of files and a list of function names and returns the full source
// of the functions that match the function names
func GetFunctionCode(file *ast.File, funcNames []string) []*ast.FuncDecl {
	var funcCode []*ast.FuncDecl
	for _, decl := range file.Decls {
		switch decl.(type) {
		case *ast.FuncDecl:
			fd := decl.(*ast.FuncDecl)
			for _, funcName := range funcNames {
				if fd.Name.Name == funcName {
					funcCode = append(funcCode, fd)
				}
			}
		}
	}
	return funcCode
}

// FuncDeclToCode takes a function declaration and returns the code as a string
func FuncDeclToCode(funcDecl *ast.FuncDecl) string {
	var buf bytes.Buffer
	printer.Fprint(&buf, token.NewFileSet(), funcDecl)
	return buf.String()
}
