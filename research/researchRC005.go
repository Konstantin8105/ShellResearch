package research

import (
	"fmt"
	"math"
	"strings"

	"github.com/Konstantin8105/CalculixRPCclient/clientCalculix"
	"github.com/Konstantin8105/Shell_generator/imperfection"
)

// RC005 - critical buckling load and ovalization
func RC005() {

	researchName := "RC005"
	createResearchDir(researchName)

	diameter := 5.0 //5.48
	height := 10.0  //13.5
	force := -1.0
	thk := 0.010 //0.005

	maxAmountPoints := 20000

	dh := math.Sqrt(math.Pi * diameter * height / (2.0 * float64(maxAmountPoints)))
	pointsOnLevel := int(math.Pi * diameter / (2.0 * dh))
	pointsOnHeight := int(height / dh)

	imperfectionAmplitude := []float64{0.0}
	imperfectionStep := 0.1

	for i := 0; i < 10; i++ {
		imperfectionAmplitude = append(imperfectionAmplitude, imperfectionAmplitude[len(imperfectionAmplitude)-1]+imperfectionStep)
	}

	var inpModels []string
	for _, amplitude := range imperfectionAmplitude {
		model, err := RegularPerfectModel(height, diameter, pointsOnLevel, pointsOnHeight, force, thk, "S4")
		if err != nil {
			fmt.Printf("Cannot mesh : %v\n", err)
			return
		}
		imperfection.Ovalization(&model, 1, amplitude*thk, 0.0)

		inpModels = append(inpModels, strings.Join(model.SaveINPtoLines(), "\n"))
	}

	client := clientCalculix.NewClient()
	factor, err := client.CalculateForBuckle(inpModels)
	if err != nil {
		fmt.Printf("Error : %v.\n Factors = %v\n", err, factor)
		return
	}

	for i := range inpModels {
		if i == 0 {
			fmt.Printf("%.5v\t%.5E\tTimoshenko = %.5E\n", imperfectionAmplitude[i]*thk, factor[i], timoshenkoLoad(thk))
		} else {
			fmt.Printf("%.5v\t%.5E\t%.2v\n", imperfectionAmplitude[i]*thk, factor[i], (factor[0]-factor[i])/factor[0]*100)
		}
	}
	for i := range inpModels {
		fmt.Printf("Imper = %.4v\tCritical stress = %.5v\n", imperfectionAmplitude[i], factor[i]/(math.Pi*diameter*thk)*1.0e-6)
	}
}
