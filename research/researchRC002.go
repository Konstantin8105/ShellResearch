package research

import (
	"fmt"
	"math"
	"os"
	"path/filepath"

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
	precision := 0.5
	pointsOnLevel := 5
	pointsOnHeight := 5
	force := -1.0
	thk := 0.005
	modelFileName, err := createResearchFilename(researchName, "model.inp")
	if err != nil {
		fmt.Println("Wrong name of model : ", err)
		return
	}
	// get result of buckling factor
	dat, err := createResearchFilename(researchName, "model.dat")
	if err != nil {
		fmt.Println("Wrong file .dat: ", err)
		return
	}

	n := 10

	calcTime := make(plotter.XYs, n)

	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	p.Title.Text = "Research : Time creating INP model file & Precision(Distance between points)"
	p.X.Label.Text = "Size of FE, meter"
	p.Y.Label.Text = "Calculation time,sec"

	for iteration := 0; iteration < n; iteration++ {

		_ = os.Remove(modelFileName)
		_ = os.Remove(dat)

		// create cylinder model
		model, err := (shellGenerator.Shell{Height: height, Diameter: diameter, Precision: precision}).GenerateMesh(pointsOnLevel, pointsOnHeight)
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

		// save file
		err = model.Save(modelFileName)
		if err != nil {
			fmt.Println("Wrong saving : ", err)
			return
		}

		// calculate
		err = executeCalculix(modelFileName)
		if err != nil {
			fmt.Println("Wrong execution : ", err)
			return
		}

		bucklingFactors, err := getBucklingFactor(dat, model.Step.AmountBucklingShapes)
		if err != nil {
			fmt.Println("Wrong buckling factor analyze : ", err)
			return
		}

		// calculate the buckling force
		//fmt.Printf("%.4v\t%.4v\n", diameter, bucklingFactors[0])
		// calculate by Timoshenko formula
		timoshenkoStress := 0.605 * 2.e11 * thk / (diameter / 2.0)
		timoshenkoForce := timoshenkoStress * math.Pi * diameter * thk
		stress2 := 750. * math.Pow(math.Pi, 2.) * 2e11 / (12. * (1. - 0.3*0.3)) * math.Pow(thk/height, 2.)
		force2 := stress2 * math.Pi * diameter * thk
		// compare the result

		size := math.Pi * diameter * height / float64(pointsOnLevel) / float64(pointsOnHeight)
		fmt.Printf("%v\t%.4v\t%.4v\t%.4v\t%.4v\n", size, thk, diameter, force*bucklingFactors[0], timoshenkoForce)

		for i := 0; i < model.Step.AmountBucklingShapes; i++ {
			fmt.Printf("%.4f\t", bucklingFactors[i])
		}
		fmt.Println("")

		fmt.Println("onLevel = ", pointsOnLevel, "\tlevels = ", pointsOnHeight)

		fmt.Println("Buckling  2 : = ", force2)

		calcTime[iteration].X = float64(iteration) // float64(pointsOnHeight) //size
		calcTime[iteration].Y = math.Abs(force * bucklingFactors[0])

		pointsOnLevel += 5  //int(float64(pointsOnLevel) * 1.05)
		pointsOnHeight += 5 //int(float64(pointsOnHeight) * 1.1)
		//diameter = diameter * 1.1

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
