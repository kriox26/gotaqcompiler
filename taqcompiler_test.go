package taqcompiler

import (
	"bufio"
	"fmt"
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
	err := testOutputProgramInit(comp)
	if err != nil {
		t.Error(err.Error())

	}
}

func TestNewCompilerFromFile(t *testing.T) {
	f, err := os.Open("testfiles/add.txt")
	if err != nil {
		t.Errorf("Failed when opening the test file. Error: %s", err.Error())
	}
	comp := NewCompilerFromFile(f)
	err = testOutputProgramInit(comp)
	if err != nil {
		t.Error(err.Error())

	}
}

func testOutputProgramInit(comp *Compiler) error {
	if comp != nil {
		for k, v := range comp.OutputProgram {
			if v != "0000000000000000" {
				return fmt.Errorf("compiler.OutputProgram should be inicialize with '0000000000000000' and instead has %s in the index %v", v, k)
			}
		}
	}
	return nil
}

func BenchmarkCompile(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		comp := NewCompilerFromString(correctProgram)
		comp.Compile()
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
	if !strings.HasPrefix(comp.CompilationErrors[0], errNoVarDefinition.Error()) {
		t.Errorf("The error should be errNoVarDefinition. But instead: %s\n", comp.CompilationErrors[0])
	}

	// When no program definition
	f, err = os.Open("testfiles/noprogramdef")
	if err != nil {
		t.Errorf("Failed when opening the test file. Error: %s", err.Error())
	}
	comp = NewCompilerFromFile(f)
	comp.Compile()
	if !strings.HasPrefix(comp.CompilationErrors[0], errNoProgramDefnition.Error()) {
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

func TestOutputProgram(t *testing.T) {
	comp := NewCompilerFromString(correctProgram)
	comp.Compile()
	f, err := os.Open("testfiles/outputProgramCorrect")
	defer f.Close()
	if err != nil {
		t.Errorf("Failed when opening the test file. Error: %s", err.Error())
	}
	f2, err := os.Create("testfiles/testingFile")
	defer f2.Close()
	if err != nil {
		t.Errorf("Failed when opening the test file. Error: %s", err.Error())
	}

	w := bufio.NewWriter(f2)
	for _, v := range comp.OutputProgram {
		w.WriteString(v + "\n")
	}

	var outputP []string
	scn := bufio.NewScanner(f)
	for scn.Scan() {
		line := scn.Text()
		// remove all tabs and new lines chars
		line = strings.Replace(strings.Replace(line, "\t", "", -1), "\n", "", -1)
		outputP = append(outputP, line)
	}

	for k, v := range comp.OutputProgram {
		if v != outputP[k] {
			t.Errorf("Wrong compilation. Line %v is %v and it should be %v", k, v, outputP[k])
		}
	}

}
