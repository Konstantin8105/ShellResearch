package main

import (
	"github.com/Konstantin8105/Convert-INP-to-STD-format/utils"
	"github.com/Konstantin8105/ShellResearch/research"
	"github.com/Konstantin8105/Shell_generator/imperfection"
)

func main() {
	//research.RC001()
	//research.RC002()
	//research.RC003()
	//research.RC004()
	//research.RC005()

	model, _ := research.RegularPerfectModel(10.0, 5.0, 100, 100, -1.0, 0.010, "S4")
	imperfection.Ovalization(&model, 1, 0.010, 0.0)

	_ = utils.CreateNewFile("Imper1thk.inp", model.SaveINPtoLines())
}
