package taqcompiler

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"strconv"
	"strings"
)

// Compiler represents the compiler that compiles a TAQ program
type Compiler struct {
	// CompilationErrors has all encountered errors when compiling
	CompilationErrors []string
	inputProgram      []string

	// OutputProgram is the 256 lines of the compiled program
	OutputProgram []string

	// OKCompilation is true if no errors were found
	OKCompilation bool
}

const progLen = 256

var (
	errVariableNameAlreadyDeclared = errors.New("This variable name has already been declared before")
	errExtraArguments              = errors.New("Instruction has more arguments than it should")
)

type program struct {
	variables    map[string]string
	instructions []string
	c            *Compiler
}

// NewCompilerFromString creates a compiler for the program given as a string
func NewCompilerFromString(inProgram string) *Compiler {
	var inpProgram []string
	// remove all tabs chars
	inProgram = strings.Replace(inProgram, "\t", "", -1)
	inpProgram = strings.Split(inProgram, "\n")
	// build the inputProgram array with all instructions and vars right here
	return &Compiler{
		inputProgram: inpProgram,
	}
}

// NewCompilerFromFile creates a compiler for the program given as a file
func NewCompilerFromFile(inProgramFile multipart.File) *Compiler {
	var inpProgram []string
	scn := bufio.NewScanner(inProgramFile)
	for scn.Scan() {
		line := scn.Text()
		// remove all tabs and new lines chars
		line = strings.Replace(strings.Replace(line, "\t", "", -1), "\n", "", -1)
		inpProgram = append(inpProgram, line)
	}
	return &Compiler{
		inputProgram: inpProgram,
	}
}

// Compile runs a compilation on the program loaded in inputProgram.
// It stores the compilation errors in Compiler.CompilationErrors and
// sets OKCompilation to false if there were any errors, otherwise
// store the program in Compiler.OutputProgram and set OKCompilation to true.
func (c *Compiler) Compile() {
	c.initCompiler()
	prog := program{
		variables:    make(map[string]string),
		instructions: make([]string, 20, 50),
		c:            c,
	}
	// start with the variables
	prog.loadVariables()
}

func (prog *program) loadVariables() {
	// Load all variables
	for k, v := range prog.c.inputProgram {
		if v == "endvar" {
			break
		}
		v = strings.Replace(strings.TrimSpace(v), " ", "", -1)
		line := strings.Split(v, ":") // line[0] holds id, and line[1] holds the value
		if line[0] != "var" {
			if _, ok := prog.variables[line[0]]; ok {
				prog.c.CompilationErrors = append(prog.c.CompilationErrors, fmt.Sprintf("%s, on line: %v", errVariableNameAlreadyDeclared.Error(), k+1))
			} else {
				prog.variables[line[0]] = line[1]
				fmt.Printf("Line %v has value %s\n", line[0], line[1])
				aux, err := strconv.Atoi(line[1])
				if err != nil {
					log.Fatalf(err.Error())
				}
				bin := strconv.FormatInt(int64(aux), 2)
				prog.c.OutputProgram[progLen-len(prog.variables)] = fmt.Sprintf("%016s", bin)
			}
		}
	}
}

func (c *Compiler) initCompiler() {
	// Load all "0000000000000000" in compiler.OutputProgram()
	c.OutputProgram = make([]string, 256, 256)
}
