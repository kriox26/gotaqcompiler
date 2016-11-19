package taqcompiler

import (
	"os"
	"strings"
	"testing"
)

var correctProgram = `var:
N1: 3
N2: 2
N3: 0
endvar
programa:
LOAD N1
ADD N2
STORE N3
HALT
end
`

func TestNewCompilerFromString(t *testing.T) {
	comp := NewCompilerFromString(correctProgram)
	if comp != nil {
		for k, v := range comp.OutputProgram {
			if v != "0000000000000000" {
				t.Errorf("compiler.OutputProgram should be inicialize with '0000000000000000' and instead has %s in the index %v", v, k)
			}
		}
	}
}

func TestNewCompilerFromFile(t *testing.T) {
	f, err := os.Open("testfiles/add.txt")
	if err != nil {
		t.Errorf("Failed when opening the test file. Error: %s", err.Error())
	}
	comp := NewCompilerFromFile(f)
	if comp != nil {
		for k, v := range comp.OutputProgram {
			if v != "0000000000000000" {
				t.Errorf("compiler.OutputProgram should be inicialize with '0000000000000000' and instead has %s in the index %v", v, k)
			}
		}
	}
}

func TestCompile(t *testing.T) {
	comp := NewCompilerFromString(correctProgram)
	comp.Compile()
	if comp.OKCompilation == false {
		t.Errorf("The compilation of the program should have been succesfull. The errors where: %v", comp.CompilationErrors)
	}

	// When no var definition
	f, err := os.Open("testfiles/novardef")
	if err != nil {
		t.Errorf("Failed when opening the test file. Error: %s", err.Error())
	}
	comp = NewCompilerFromFile(f)
	comp.Compile()
	if comp.CompilationErrors[0] != errNoVarDefinition.Error() {
		t.Errorf("The error should be errNoVarDefinition. But instead: %s\n", comp.CompilationErrors[0])
	}

	// When no program definition
	f, err = os.Open("testfiles/noprogramdef")
	if err != nil {
		t.Errorf("Failed when opening the test file. Error: %s", err.Error())
	}
	comp = NewCompilerFromFile(f)
	comp.Compile()
	if comp.CompilationErrors[0] != errNoProgramDefnition.Error() {
		t.Errorf("The error should be errNoProgramDefnition. But instead: %s\n", comp.CompilationErrors[0])
	}

	// When repetead variable names
	f, err = os.Open("testfiles/morethanonevar")
	if err != nil {
		t.Errorf("Failed when opening the test file. Error: %s", err.Error())
	}
	comp = NewCompilerFromFile(f)
	comp.Compile()
	if !strings.HasPrefix(comp.CompilationErrors[0], errVariableNameAlreadyDeclared.Error()) {
		t.Errorf("The error should be errVariableNameAlreadyDeclared. But instead: %s\n", comp.CompilationErrors[0])
	}

	// When repetead label definitions
	f, err = os.Open("testfiles/morethanonelabel")
	if err != nil {
		t.Errorf("Failed when opening the test file. Error: %s", err.Error())
	}
	comp = NewCompilerFromFile(f)
	comp.Compile()
	if !strings.HasPrefix(comp.CompilationErrors[0], errLabelAlreadyDefined.Error()) {
		t.Errorf("The error should be errLabelAlreadyDefined. But instead: %s\n", comp.CompilationErrors[0])
	}
}
