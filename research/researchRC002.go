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
	"github.com/gonum/plot/vg"
)

// RC002 - compare with Timoshenko formula
// Goal :
// * if more FE, then more precision
// * if FE is many, then precision is lost
func RC002() {

	researchName := "RC002"
	createResearchDir(researchName)

	diameter := 5.48
	height := 13.5
	force := -1.0
	thk := 0.005

	maxAmountPoints := 40000

	pointsOnLevel := []int{5}
	pointsOnHeight := []int{5}

	for {
		lStep := 5
		hStep := 5
		l := pointsOnLevel[len(pointsOnLevel)-1] + lStep
		h := pointsOnHeight[len(pointsOnHeight)-1] + hStep
		pointsOnLevel = append(pointsOnLevel, l)
		pointsOnHeight = append(pointsOnHeight, h)

		if (l+lStep)*(h+hStep) >= maxAmountPoints {
			break
		}
	}

	n := len(pointsOnLevel)

	typeOfFEs := []string{
		"S4",
		"S4R",
		"S8",
		"S8R",
		"S3",
		"S6",
	}

	calcTime := make(plotter.XYs, n)
	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	p.Title.Text = "Research : Buckling force depends of finite element size"
	p.X.Label.Text = "Amount of points"
	p.Y.Label.Text = "Force, N"

	calcError := make(plotter.XYs, n)
	p2, err := plot.New()
	if err != nil {
		panic(err)
	}
	p2.Title.Text = "Research : Error depends of finite element size"
	p2.X.Label.Text = "Amount of points"
	p2.Y.Label.Text = "Error, %"

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

	for indexFE, fe := range typeOfFEs {
		var inpModels []string
		var amountPoints []int

		for i := 0; i < n; i++ {
			fmt.Println("Prepare model : ", len(amountPoints))

			model, err := ShellModel(height, diameter, pointsOnLevel[i], pointsOnHeight[i], force, thk, fe)
			if err != nil {
				fmt.Printf("Cannot mesh : %v\n", err)
				return
			}

			inpModels = append(inpModels, model)
			amountPoints = append(amountPoints, pointsOnLevel[i]*pointsOnHeight[i])
		}

		client := clientCalculix.NewClient()
		factor, err := client.CalculateForBuckle(inpModels)
		if err != nil {
			fmt.Printf("Error : %v.\n Factors = %v\n", err, factor)
			return
		}

		fmt.Fprintf(f, "Type of finite element : %v\n", fe)

		for i := range inpModels {
			calcTime[i].X = float64(amountPoints[i])
			calcTime[i].Y = math.Abs(force * factor[i])
			calcError[i].X = float64(amountPoints[i])
			f0 := force * factor[i]
			ft := -0.6052275 * 2. * math.Pi * math.Pow(0.005, 2.) * 2.0e11
			// amount =    24025	f0 = -3.317E+06 ft = -1.901E+07 error = +0.83 %
			var e float64
			if (math.Abs(f0)) > math.Abs(ft) {
				e = 1. - math.Abs(ft)/math.Abs(f0)
			} else {
				e = 1. - math.Abs(f0)/math.Abs(ft)
			}
			e = e * 100.
			calcError[i].Y = e
			fmt.Fprintf(f, "amount = %8v\tf0 = %2.3E ft = %2.3E error = %+2.2f %v \n", amountPoints[i], f0, ft, e, "%")
		}

		{
			// Make a line plotter and set its style.
			l, err := plotter.NewLine(calcTime)
			if err != nil {
				panic(err)
			}
			l.LineStyle.Width = vg.Points(1)
			l.LineStyle.Color = GetColor(float64(indexFE) / float64(len(typeOfFEs)))

			// Add the plotters to the plot, with a legend
			// entry for each
			p.Add(l)
			p.Legend.Add(fmt.Sprintf("Finite Element :%v", fe), l)
		}
		{
			// Make a line plotter and set its style.
			l, err := plotter.NewLine(calcError)
			if err != nil {
				panic(err)
			}
			l.LineStyle.Width = vg.Points(1)
			l.LineStyle.Color = GetColor(float64(indexFE) / float64(len(typeOfFEs)))

			// Add the plotters to the plot, with a legend
			// entry for each
			p2.Add(l)
			p2.Legend.Add(fmt.Sprintf("Finite Element :%v", fe), l)
		}

		// Save the plot to a PNG file.
		if err := p.Save(16*vg.Inch, 8*vg.Inch, string(researchFolder+string(filepath.Separator)+researchName+string(filepath.Separator)+researchName+".png")); err != nil {
			panic(err)
		}
		if err := p2.Save(16*vg.Inch, 8*vg.Inch, string(researchFolder+string(filepath.Separator)+researchName+string(filepath.Separator)+researchName+"_error.png")); err != nil {
			panic(err)
		}
	}
	err = f.Close()
	if err != nil {
		fmt.Println("Cannot close file")
		return
	}
}

// ShellModel - shell model
func ShellModel(height float64, diameter float64, pointsOnLevel, pointsOnHeight int, force, thk float64, typeOfFE string) (resultInp string, err error) {
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
