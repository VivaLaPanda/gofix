// The following is a Go script called "gofix" used to automatically fix failing Golang tests by using GPT-4.
// The user calls it by calling something like `go test ./... -run ^Test_divide$` and piping the results into gofix. Gofix then:
// - Figures out which files tests are failing in
// - Gets the code of the tests that are failing
// - Figures out what functions that test code depends on
// - Figures our what functions the test code calls
// - Get the text of those functions
// - Calls GPT-4 to determine how to fix the functions so that the test passes
package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func GetTestOutput() (string, error) {
	bytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return "", fmt.Errorf("failed to read from stdin: %v", err)
	}

	return string(bytes), nil
}

func main() {
	// Get the output of the test command from stdin (go test ./... -json)
	testOutput, err := GetTestOutput()

	// Print for debugging
	// fmt.Println(testOutput)

	// Get the files that contain failing tests
	failingTests, err := GetFailingTests(testOutput)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Files:")
	for filePath := range failingTests {
		fmt.Println(filePath)
	}

	for filePath, failingTestInfo := range failingTests {
		// Get the file from the file path
		file, err := GetFileFromPath(filePath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		for _, failedTest := range failingTestInfo {
			// Print the name of the failing test
			fmt.Println("Failing test:", failedTest)

			// Get the code of the failing tests
			testCode := GetTestCode(file, failedTest.TestName)

			// Get the functions that the test code depends on
			funcs := GetFunctionDependencies(testCode)

			// Get the name of the non-test file by stripping the _test.go suffix
			nonTestFilePath := strings.Replace(filePath, "_test.go", ".go", 1)

			// Get the file from the non-test file path
			nonTestFile, err := GetFileFromPath(nonTestFilePath)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			// Get the code of the functions that the test code depends on
			funcAst := GetFunctionCode(nonTestFile, funcs)

			// Make the testCode string
			var testCodeStrings []string
			for _, decl := range testCode {
				testCodeStrings = append(testCodeStrings, FuncDeclToCode(decl))
			}

			// Make the function code strings
			var funcCodeStrings []string
			for _, decl := range funcAst {
				funcCodeStrings = append(funcCodeStrings, FuncDeclToCode(decl))
			}

			// Print out the state here to make sure it's working
			// fmt.Println("Test Code:")
			//for _, testCodeString := range testCodeStrings {
			//	fmt.Println(testCodeString)
			//}
			//// fmt.Println("Function Code:")
			//for _, funcCodeString := range funcCodeStrings {
			//	fmt.Println(funcCodeString)
			//}

			// Construct the prompt for GPT-4
			prompt := ConstructPrompt(testCodeStrings, funcCodeStrings, failedTest)

			// Print the prompt for debugging
			// fmt.Println(prompt)

			// Call GPT-4
			resp, err := callGpt4(prompt)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			// Print the response for debugging
			fmt.Println(resp)
		}
	}
}

// ConstructPrompt takes in the test code, the function code, and failingTestInfo and constructs the prompt for GPT-4
func ConstructPrompt(testCodes, funcCodes []string, failingTestInfo FailingTest) string {
	// Join test code and funcCode to look like code blocks .i.e.
	// ```go
	// func add(a, b int) int {
	// 	return a + b
	// }
	// ```
	testCodes = append([]string{"```go"}, testCodes...)
	testCodes = append(testCodes, "```")
	funcCodes = append([]string{"```go"}, funcCodes...)
	funcCodes = append(funcCodes, "```")

	testCode := strings.Join(testCodes, "\n")
	funcCode := strings.Join(funcCodes, "\n")

	prompt := "You are part of an elite automated software fixing team. You will be given the code of a failing test, and the failing test output.\n\n" +
		"Please provide an explanation of why the test is failing (either the test code being bad or the function code being bad) and how to fix it. Generally bias\n" +
		"towards assuming the test code is correct and the function code is what needs to be fixed.\n" +
		"Assume that comments correctly describe the desired functionality of the function.\n\n" +
		"You should recommend specific code changes if possible.':" +
		"\n\nTest Code:\n" + testCode +
		"\n\nFunction Code:\n" + funcCode +
		"\nTest Error: " + failingTestInfo.TestError +
		"\n\nYour explanation and recommended changes:"

	return prompt
}

//// RunGoTest runs "go test ./... -json and returns the output
//func RunGoTest() (string, error) {
//	// Run "go test ./... -json and get the output using the os/exec package
//	cmd := exec.Command("go", "test", "./...", "-json")
//	out, err := cmd.Output()
//	if err != nil {
//		return "", err
//	}
//
//	return string(out), nil
//}
