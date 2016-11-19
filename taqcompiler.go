package taqcompiler

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
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
	errLabelAlreadyDefined         = errors.New("This label has already been used")
	errExtraArguments              = errors.New("Instruction has more arguments than it should")
	errNoProgramDefnition          = errors.New("The program definition wasn't found")
	errNoVarDefinition             = errors.New("The var definition wasn't found")
)

var codeOps = map[string]string{
	"LOAD":  "00000010",
	"ADD":   "00000011",
	"STORE": "00000001",
	"JUMP":  "00000000",
	"SUB":   "00000110",
	"AND":   "00000100",
	"JZ":    "00000101",
	"NOP":   "00000111",
	"HALT":  "00001000",
}

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
		inputProgram:  inpProgram,
		OutputProgram: initOutputProgram(),
	}
}

// NewCompilerFromFile creates a compiler for the program given as a file
func NewCompilerFromFile(inProgramFile io.Reader) *Compiler {
	var inpProgram []string
	scn := bufio.NewScanner(inProgramFile)
	for scn.Scan() {
		line := scn.Text()
		// remove all tabs and new lines chars
		line = strings.Replace(strings.Replace(line, "\t", "", -1), "\n", "", -1)
		inpProgram = append(inpProgram, line)
	}
	return &Compiler{
		inputProgram:  inpProgram,
		OutputProgram: initOutputProgram(),
	}
}

// Compile runs a compilation on the program loaded in inputProgram.
// It stores the compilation errors in Compiler.CompilationErrors and
// sets OKCompilation to false if there were any errors, otherwise
// store the program in Compiler.OutputProgram and set OKCompilation to true.
func (c *Compiler) Compile() {
	prog := program{
		variables:    make(map[string]string),
		instructions: make([]string, 20, 50),
		c:            c,
	}
	// start with the variables
	prog.loadVariables()
	prog.loadInstructions()
	if len(prog.c.CompilationErrors) == 0 {
		prog.c.OKCompilation = true
	} else {
		prog.c.OKCompilation = false
	}
}

func (prog *program) loadVariables() {
	// Load all variables
	if strings.TrimSpace(prog.c.inputProgram[0]) != "var:" {
		prog.c.CompilationErrors = append(prog.c.CompilationErrors, errNoVarDefinition.Error())
		return
	}
	for k, v := range prog.c.inputProgram {
		if strings.TrimSpace(v) == "endvar" {
			break
		}
		v = strings.Replace(strings.TrimSpace(v), " ", "", -1)
		line := strings.Split(v, ":") // line[0] holds id, and line[1] holds the value
		if line[0] != "var" {
			if _, ok := prog.variables[line[0]]; ok {
				prog.c.CompilationErrors = append(prog.c.CompilationErrors, fmt.Sprintf("%s, on line: %v", errVariableNameAlreadyDeclared.Error(), k+1))
			} else {
				prog.variables[line[0]] = line[1]
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

func (prog *program) loadInstructions() {
	var jumps []int
	code := prog.c.inputProgram
	p, err := indexOfInst(code, "programa:")
	start := p + 1
	if err != nil {
		prog.c.CompilationErrors = append(prog.c.CompilationErrors, err.Error())
		return
	}
	for i := start; i < len(code); i++ {
		if code[i] == "end" {
			break
		}
		inst := strings.Split(code[i], " ")
		if len(inst) > 2 {
			// it has a label, so we need to store the memory address
			if _, ok := codeOps[strings.Replace(inst[0], ":", "", -1)]; ok {
				prog.c.CompilationErrors = append(prog.c.CompilationErrors, fmt.Sprintf("%s, on line: %v", errLabelAlreadyDefined.Error(), i+1))
			} else {
				codeOps[strings.Replace(inst[0], ":", "", -1)] = strings.Replace(fmt.Sprintf("%08v", strconv.FormatInt(int64(i-start), 2)), " ", "", -1)
			}
		} else if len(inst) == 1 {
			// HALT
			prog.c.OutputProgram[i-start] = fmt.Sprintf("%016s", codeOps[inst[0]])
			continue
		}
		codeOp := inst[len(inst)-2]
		if strings.ToLower(codeOp) == "jump" || strings.ToLower(codeOp) == "jz" {
			// add to jumps
			jumps = append(jumps, i)
		} else {
			op := inst[len(inst)-1]
			ind := prog.indexOfOp(op)
			prog.c.OutputProgram[i-start] = codeOps[codeOp] + fmt.Sprintf("%08s", strings.Replace(strconv.FormatInt(int64(255-ind), 2), " ", "0", -1))
		}
	}

	for _, k := range jumps {
		inst := strings.Split(prog.c.inputProgram[k], " ")
		prog.c.OutputProgram[k-start] = fmt.Sprintf("%s%s", codeOps[inst[len(inst)-2]], codeOps[strings.TrimSpace(inst[len(inst)-1])])
	}

}

func (prog *program) indexOfOp(op string) int {
	i := 0
	for k := range prog.variables {
		if strings.TrimSpace(op) == strings.TrimSpace(k) {
			return i
		}
		i++
	}
	return -1
}

func indexOfInst(sl []string, inst string) (int, error) {
	for k, v := range sl {
		if strings.TrimSpace(v) == inst {
			return k, nil
		}
	}
	return -1, errNoProgramDefnition
}

func initOutputProgram() []string {
	// Load all "0000000000000000" in compiler.OutputProgram()
	outputProgram := make([]string, 256, 256)
	for k := range outputProgram {
		outputProgram[k] = "0000000000000000"
	}
	return outputProgram
}
