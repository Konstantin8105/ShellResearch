package research

import (
	"fmt"
	"path/filepath"

	"github.com/Konstantin8105/CalculixRPCclient/clientCalculix"
	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/plotutil"
	"github.com/gonum/plot/vg"
)

// RC003 - research
// research amount FE and precision
func RC003() {

	researchName := "RC003"
	createResearchDir(researchName)

	diameter := 5.8
	height := 12.0
	startPointsOnLevel := 10
	stepOnLevel := 4
	startPointsOnHeight := 10
	stepOnHeight := 4
	force := -1.0
	thk := 0.005

	n := 20

	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	p.Title.Text = "Research : Buckling force depends of finite element size"
	p.X.Label.Text = "amount of points on level"
	p.Y.Label.Text = "Force, N"

	for i := 0; i < n; i++ {

		var inpModels []string
		graph := make(plotter.XYs, n)

		pointsOnLevel := startPointsOnLevel + i*stepOnLevel

		for j := 0; j < n; j++ {
			pointsOnHeight := startPointsOnHeight + j*stepOnHeight
			model, err := ShellModel(height, diameter, pointsOnLevel, pointsOnHeight, force, thk)
			if err != nil {
				fmt.Println("Cannot mesh")
				return
			}

			inpModels = append(inpModels, model)
		}

		var client clientCalculix.ClientCalculix
		client.Manager = *clientCalculix.NewServerManager()
		factors, err := client.CalculateForBuckle(inpModels)
		if err != nil {
			fmt.Printf("Error : %v.\n Factors = %v\n", err, factors)
			return
		}

		err = plotutil.AddLinePoints(p,
			fmt.Sprintf("PointOnLevel:%v", pointsOnLevel), graph,
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
