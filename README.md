gofix is a tool used to automatically fix failing Golang tests by using GPT-4.
The user calls it by calling something like `go test ./... -run ^Test_divide$` and piping the results into gofix. Gofix then:
- Figures out which files tests are failing in
- Gets the code of the tests that are failing
- Figures out what functions that test code depends on
- Figures our what functions the test code calls
- Get the text of those functions
- Calls GPT-4 to determine how to fix the functions so that the test passes

## NOTE

Reading in from stdin is broken right now, instead call `go test ./... -json -run ^Test_divide$ > test_output.jsonl`
and then run the executable with `gofix`. I'll fix this soon.