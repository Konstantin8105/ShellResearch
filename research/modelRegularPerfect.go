package research

import (
	"fmt"
	"strings"

	"github.com/Konstantin8105/Convert-INP-to-STD-format/inp"
	"github.com/Konstantin8105/Shell_generator/shellGenerator"
)

// ShellModel - shell model
func regularPerfectModel(height float64, diameter float64, pointsOnLevel, pointsOnHeight int, force, thk float64, typeOfFE string) (resultInp string, err error) {
	var model inp.Format

	precision := 0.5

	// create cylinder model
	model, err = (shellGenerator.Shell{Height: height, Diameter: diameter, Precision: precision}).GenerateMesh(pointsOnLevel, pointsOnHeight)
	if err != nil {
		fmt.Println("Wrong mesh : ", err)
		return
	}

	// modify finite element
	s4, err := inp.GetFiniteElementByName("S4")
	if err != nil {
		return "", fmt.Errorf("Error : %v", err)
	}
	s, err := inp.GetFiniteElementByName(typeOfFE)
	if err != nil {
		return "", fmt.Errorf("Error : %v", err)
	}
	err = model.ChangeTypeFiniteElement(s4, s)
	if err != nil {
		return "", fmt.Errorf("Error in change FE: %v", err)
	}

	// create fixed points
	fixName := "fix"
	_ = model.AddNamedNodesOnLevel(0, fixName)
	model.Boundary = append(model.Boundary, inp.BoundaryProperty{
		NodesByName:   fixName,
		StartFreedom:  1,
		FinishFreedom: 1,
		Value:         0,
	})
	model.Boundary = append(model.Boundary, inp.BoundaryProperty{
		NodesByName:   fixName,
		StartFreedom:  2,
		FinishFreedom: 2,
		Value:         0,
	})
	model.Boundary = append(model.Boundary, inp.BoundaryProperty{
		NodesByName:   fixName,
		StartFreedom:  3,
		FinishFreedom: 3,
		Value:         0,
	})

	// create load points
	loadName := "load"
	size := model.AddNamedNodesOnLevel(height, loadName)
	model.Boundary = append(model.Boundary, inp.BoundaryProperty{
		NodesByName:   loadName,
		StartFreedom:  1,
		FinishFreedom: 1,
		Value:         0,
	})
	model.Boundary = append(model.Boundary, inp.BoundaryProperty{
		NodesByName:   loadName,
		StartFreedom:  3,
		FinishFreedom: 3,
		Value:         0,
	})
	forcePerPoint := force / float64(size)
	model.Step.Loads = append(model.Step.Loads, inp.Load{
		NodesByName: loadName,
		Direction:   2,
		LoadValue:   forcePerPoint,
	})

	// addshell property
	model.ShellSections = append(model.ShellSections, inp.ShellSection{
		ElementName: shellGenerator.ShellName,
		Thickness:   thk,
	})

	// create linear buckling
	model.Step.AmountBucklingShapes = 1

	lines := model.SaveINPtoLines()

	return strings.Join(lines, "\n"), nil
}
