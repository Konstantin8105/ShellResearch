package research

import (
	"fmt"
	"math"
	"path/filepath"

	"github.com/Konstantin8105/CalculixRPCclient/clientCalculix"
	"github.com/Konstantin8105/Convert-INP-to-STD-format/inp"
	"github.com/Konstantin8105/Shell_generator/shellGenerator"
	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/plotutil"
	"github.com/gonum/plot/vg"
)

// RC002 - compare with Timoshenko formula
func RC002() {

	researchName := "RC002"
	createResearchDir(researchName)

	diameter := 3.
	height := 1.
	pointsOnLevel := 5
	pointsOnHeight := 5
	force := -1.0
	thk := 0.005

	n := 20

	calcTime := make(plotter.XYs, n)

	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	p.Title.Text = "Research : Buckling force depends of finite element size"
	p.X.Label.Text = "Iteration of size"
	p.Y.Label.Text = "Force, N"

	var inpModels []string

	for iteration := 0; iteration < n; iteration++ {
		fmt.Println("Prepare model : ", iteration)

		model, err := ShellModel(height, diameter, pointsOnLevel, pointsOnHeight, force, thk)
		if err != nil {
			fmt.Println("Cannot mesh")
			return
		}

		inpModels = append(inpModels, model)

		pointsOnLevel += 5
		pointsOnHeight += 5
	}

	var client clientCalculix.ClientCalculix
	client.Manager = *clientCalculix.NewServerManager()
	factor, err := client.CalculateForBuckle(inpModels)
	if err != nil {
		fmt.Println("Error : ", err)
		return
	}

	for i := range inpModels {
		calcTime[i].X = float64(i)
		calcTime[i].Y = math.Abs(force * factor[i])
	}

	err = plotutil.AddLinePoints(p,
		fmt.Sprintf("Iteration"), calcTime,
	)
	if err != nil {
		panic(err)
	}
	// Save the plot to a PNG file.
	if err := p.Save(8*vg.Inch, 8*vg.Inch, string(researchFolder+string(filepath.Separator)+researchName+string(filepath.Separator)+researchName+".png")); err != nil {
		panic(err)
	}
}

// ShellModel - shell model
func ShellModel(height float64, diameter float64, pointsOnLevel, pointsOnHeight int, force, thk float64) (resultInp string, err error) {
	var model inp.Format

	precision := 0.5

	// create cylinder model
	model, err = (shellGenerator.Shell{Height: height, Diameter: diameter, Precision: precision}).GenerateMesh(pointsOnLevel, pointsOnHeight)
	if err != nil {
		fmt.Println("Wrong mesh : ", err)
		return
	}

	// create fixed points
	fixName := "fix"
	model.AddNamedNodesOnLevel(0, fixName)
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
	model.AddNamedNodesOnLevel(height, loadName)
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
	forcePerPoint := force / float64(pointsOnLevel)
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
	var buffer string
	for _, line := range lines {
		buffer += line + "\n"
	}
	return buffer, nil
}
