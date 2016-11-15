package taqcompiler

import "mime/multipart"

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

// NewCompilerFromString creates a compiler for the program given as a string
func NewCompilerFromString(inProgram string) *Compiler {
	var inpProgram []string
	// Build the inputProgram array with all instructions and vars right here
	return &Compiler{
		inputProgram: inpProgram,
	}
}

// NewCompilerFromFile creates a compiler for the program given as a file
func NewCompilerFromFile(inProgramFile multipart.File) *Compiler {
	var inpProgram []string
	return &Compiler{
		inputProgram: inpProgram,
	}
}

// Compile runs a compilation on the program loaded in inputProgram. Stores the output of the pprogram in OutputProgram if there were no CompilationErrors, if OKCompilation is true.
func (c *Compiler) Compile() {
	c.OKCompilation = true
	c.OutputProgram = make([]string, 10)
	c.OutputProgram[0] = "HOLAAAA"
}
