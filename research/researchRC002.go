package research

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/Konstantin8105/CalculixRPCclient/clientCalculix"
	"github.com/Konstantin8105/Convert-INP-to-STD-format/inp"
	"github.com/Konstantin8105/Shell_generator/shellGenerator"
	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/plotutil"
	"github.com/gonum/plot/vg"
)

// RC002 - compare with Timoshenko formula
// Goal :
// * if more FE, then more precision
// * if FE is many, then precision is lost
func RC002() {

	researchName := "RC002"
	createResearchDir(researchName)

	diameter := 5.8
	height := 12.0
	pointsOnLevel := 10
	pointsOnHeight := 10
	force := -1.0
	thk := 0.005

	n := 50

	calcTime := make(plotter.XYs, n)
	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	p.Title.Text = "Research : Buckling force depends of finite element size"
	p.X.Label.Text = "Iteration of size"
	p.Y.Label.Text = "Force, N"

	calcError := make(plotter.XYs, n)
	p2, err := plot.New()
	if err != nil {
		panic(err)
	}
	p2.Title.Text = "Research : Error depends of finite element size"
	p2.X.Label.Text = "Iteration of size"
	p2.Y.Label.Text = "Error, %"

	var inpModels []string

	for iteration := 0; iteration < n; iteration++ {
		fmt.Println("Prepare model : ", iteration)

		model, err := ShellModel(height, diameter, pointsOnLevel, pointsOnHeight, force, thk)
		if err != nil {
			fmt.Println("Cannot mesh")
			return
		}

		inpModels = append(inpModels, model)

		pointsOnLevel += 4
		pointsOnHeight += 8
	}

	var client clientCalculix.ClientCalculix
	client.Manager = *clientCalculix.NewServerManager()
	factor, err := client.CalculateForBuckle(inpModels)
	if err != nil {
		fmt.Printf("Error : %v.\n Factors = %v\n", err, factor)
		return
	}

	// create text file
	file := string(researchFolder + string(filepath.Separator) + researchName + string(filepath.Separator) + researchName + ".txt")
	// check file is exist
	if _, err := os.Stat(file); os.IsNotExist(err) {
		// create file
		newFile, err := os.Create(file)
		if err != nil {
			return
		}
		err = newFile.Close()
		if err != nil {
			return
		}
	}
	// open file
	f, err := os.OpenFile(file, os.O_WRONLY, 0777)
	if err != nil {
		return
	}
	for i := range inpModels {
		calcTime[i].X = float64(i)
		calcTime[i].Y = math.Abs(force * factor[i])
		calcError[i].X = float64(i)
		f0 := force * factor[i]
		ft := -0.6052275 * 2. * math.Pi * math.Pow(0.005, 2.) * 2.0e11
		e := (math.Abs(f0) - math.Abs(ft)) / math.Abs(ft) * 100.
		calcError[i].Y = e
		fmt.Fprintf(f, "f0 = %2.3E ft = %2.3E error = %+2.2f %v \n", f0, ft, e, "%")
	}

	err = f.Close()
	if err != nil {
		return
	}

	err = plotutil.AddLinePoints(p,
		fmt.Sprintf("Iteration"), calcTime,
	)
	if err != nil {
		panic(err)
	}

	err = plotutil.AddLinePoints(p2,
		fmt.Sprintf("Iteration"), calcError,
	)
	if err != nil {
		panic(err)
	}
	// Save the plot to a PNG file.
	if err := p.Save(8*vg.Inch, 8*vg.Inch, string(researchFolder+string(filepath.Separator)+researchName+string(filepath.Separator)+researchName+".png")); err != nil {
		panic(err)
	}
	if err := p2.Save(8*vg.Inch, 8*vg.Inch, string(researchFolder+string(filepath.Separator)+researchName+string(filepath.Separator)+researchName+"_error.png")); err != nil {
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

	fmt.Println("Add other model property")
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

	fmt.Println("End of adding property")

	lines := model.SaveINPtoLines()
	/*var buffer string
	for _, line := range lines {
		buffer += line + "\n"
	}*/

	fmt.Println("Return inp like string")
	return /* buffer*/ strings.Join(lines, "\n"), nil
}
