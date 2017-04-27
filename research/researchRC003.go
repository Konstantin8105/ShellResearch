package research

import (
	"fmt"
	"math"
	"path/filepath"

	"github.com/Konstantin8105/CalculixRPCclient/clientCalculix"
	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/vg"
)

// RC003 - research
// research ratio (height/width) of FE and precision
func RC003() {

	researchName := "RC003"
	createResearchDir(researchName)

	diameter := 5.48
	height := 13.5
	force := -1.0
	thk := 0.005

	type pointCase struct {
		pointsOnLevel  int
		pointsOnHeight int
	}

	maxAmountPoints := []int{500, 2000, 5000, 10000, 20000, 30000, 40000}
	stepPoint := []int{5, 5, 20, 20, 50, 100, 100}

	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	p.Title.Text = "Research : precision at ratio(height/width or width/height) for rectange finite element S4"
	p.X.Label.Text = "Ratio (r<0 height/width-1.0);(r>0 1.0-width/height)"
	p.Y.Label.Text = "Error,%"

	for maxI := range maxAmountPoints {

		var cases []pointCase

		for l := 5; l < maxAmountPoints[maxI]; l += stepPoint[maxI] {
			h := int(float64(maxAmountPoints[maxI]) / float64(l))
			if l*h > maxAmountPoints[maxI] {
				continue
			}
			if h < 5 {
				continue
			}
			if len(cases) > 0 {
				c := cases[len(cases)-1]
				if c.pointsOnHeight == h && c.pointsOnLevel < l {
					cases[len(cases)-1] = pointCase{pointsOnLevel: l, pointsOnHeight: h}
					continue
				}
			}
			cases = append(cases, pointCase{pointsOnLevel: l, pointsOnHeight: h})
		}

		var inpModels []string
		for i, c := range cases {
			fmt.Println("Prepare model : ", i)

			model, err := ShellModel(height, diameter, c.pointsOnLevel, c.pointsOnHeight, force, thk, "S4")
			if err != nil {
				fmt.Printf("Cannot mesh : %v\n", err)
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

		graph := make(plotter.XYs, len(cases))

		for i := range factors {
			dl := math.Pi * diameter / float64(cases[i].pointsOnLevel)
			dh := height / float64(cases[i].pointsOnHeight)
			var ratio float64
			if dl > dh {
				ratio = dh/dl - 1
			} else {
				ratio = 1 - dl/dh
			}
			f0 := force * factors[i]
			ft := timoshenkoLoad(thk)
			var e float64
			if (math.Abs(f0)) > math.Abs(ft) {
				e = 1. - math.Abs(ft)/math.Abs(f0)
			} else {
				e = 1. - math.Abs(f0)/math.Abs(ft)
			}
			e = e * 100.
			graph[i].X = ratio
			graph[i].Y = e
			fmt.Printf("%5v %5v %.2e %.2e %.3v %.5e %.5e\t%3.1f\n",
				cases[i].pointsOnLevel,
				cases[i].pointsOnHeight,
				dl,
				dh,
				ratio,
				f0,
				ft,
				e)

		}

		{
			// Make a line plotter and set its style.
			l, err := plotter.NewLine(graph)
			if err != nil {
				panic(err)
			}
			l.LineStyle.Width = vg.Points(1)
			l.LineStyle.Color = GetColor(float64(maxI) / float64(len(maxAmountPoints)))

			// Add the plotters to the plot, with a legend
			// entry for each
			p.Add(l)
			p.Legend.Add(fmt.Sprintf("Graph with amount FE:%v", maxAmountPoints[maxI]), l)
		}

		// Save the plot to a PNG file.
		if err := p.Save(16*vg.Inch, 8*vg.Inch, string(researchFolder+string(filepath.Separator)+researchName+string(filepath.Separator)+researchName+fmt.Sprintf("-%v", maxAmountPoints[maxI])+".png")); err != nil {
			panic(err)
		}
	}
}
