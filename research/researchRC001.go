package research

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Konstantin8105/Shell_generator/shellGenerator"
	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/plotutil"
	"github.com/gonum/plot/vg"
)

// RC001 - start research
func RC001() {

	researchName := "RC001"
	createResearchDir(researchName)

	Diameter := 5.
	Height := 15.
	n := 10
	calcTime := make(plotter.XYs, n)

	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	p.Title.Text = "Research : Time creating INP model file & Precision(Distance between points)"
	p.X.Label.Text = "Size of FE, meter"
	p.Y.Label.Text = "Calculation time,sec"

	for data := 0; data < 20; data++ {
		precisionMax := 1.0
		precisionMin := 0.010
		filename := fmt.Sprintf("%v%v.inp", researchFolder, researchName)
		for step := range calcTime {
			precision := precisionMin + (precisionMax-precisionMin)/float64(step+1.)

			calcTime[step].X = precision

			sh := shellGenerator.Shell{Height: Height, Diameter: Diameter, Precision: precision}
			start := time.Now()

			err := sh.GenerateINP(filename)
			if err != nil {
				panic(err)
			}

			calcTime[step].Y = time.Since(start).Seconds()
			_ = os.Remove(filename)
		}

		err = plotutil.AddLinePoints(p,
			fmt.Sprintf("Iteration %v", data), calcTime,
		)
		if err != nil {
			panic(err)
		}
	}
	// Save the plot to a PNG file.
	if err := p.Save(8*vg.Inch, 8*vg.Inch, string(researchFolder+string(filepath.Separator)+researchName+string(filepath.Separator)+researchName+".png")); err != nil {
		panic(err)
	}

}
