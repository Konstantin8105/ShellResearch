package research

import (
	"fmt"
	"image/color"
	"math"
	"path/filepath"

	"github.com/Konstantin8105/CalculixRPCclient/clientCalculix"
	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/vg"
)

// RC003 - research
// research amount FE and precision
func RC003() {

	researchName := "RC003"
	createResearchDir(researchName)

	diameter := 5.48
	height := 13.5
	startPointsOnLevel := 100
	stepOnLevel := 50
	startPointsOnHeight := 100
	stepOnHeight := 50
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

		client := clientCalculix.NewClient()
		factors, err := client.CalculateForBuckle(inpModels)
		if err != nil {
			fmt.Printf("Error : %v.\n Factors = %v\n", err, factors)
			return
		}

		for j := 0; j < n; j++ {
			pointsOnHeight := startPointsOnHeight + j*stepOnHeight
			graph[j].X = float64(pointsOnHeight)
			graph[j].Y = float64(math.Abs(force * factors[j]))
		}

		// Make a line plotter and set its style.
		l, err := plotter.NewLine(graph)
		if err != nil {
			panic(err)
		}
		l.LineStyle.Width = vg.Points(1)
		l.LineStyle.Color = getColor(float64(i) / float64(n))

		// Add the plotters to the plot, with a legend
		// entry for each
		p.Add(l)
		p.Legend.Add(fmt.Sprintf("PointOnLevel:%v", pointsOnLevel), l)

		// Save the plot to a PNG file.
		if err := p.Save(8*vg.Inch, 8*vg.Inch, string(researchFolder+string(filepath.Separator)+researchName+string(filepath.Separator)+researchName+".Graph"+fmt.Sprintf("%v", i)+".png")); err != nil {
			panic(err)
		}
	}
}

// index from 0 to 1
func getColor(index float64) color.RGBA {
	colorShema := []color.RGBA{
		color.RGBA{R: 255, G: 0, B: 0, A: 255},
		color.RGBA{R: 255, G: 255, B: 0, A: 255},
		color.RGBA{R: 255, G: 0, B: 255, A: 255},
		color.RGBA{R: 0, G: 255, B: 255, A: 255},
	}
	i := int(index*float64(len(colorShema)-1) - 1e-5)
	if i >= len(colorShema) {
		i = len(colorShema) - 1
	}
	if i < 0 {
		i = 0
	}
	index = index*float64(len(colorShema)) - float64(int(index*float64(len(colorShema))))
	if index < 0. {
		index = 1e5
	}
	if index > 1. {
		index = 1. - 1e-5
	}
	return color.RGBA{
		R: uint8(float64(colorShema[i].R) + float64(float64(colorShema[i+1].R-colorShema[i].R)*index)),
		G: uint8(float64(colorShema[i].G) + float64(float64(colorShema[i+1].G-colorShema[i].G)*index)),
		B: uint8(float64(colorShema[i].B) + float64(float64(colorShema[i+1].B-colorShema[i].B)*index)),
		A: uint8(float64(colorShema[i].A) + float64(float64(colorShema[i+1].A-colorShema[i].A)*index)),
	}
}
