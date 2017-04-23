// package research
//
// +- Research : --------+
// + Many tasks          +
// +---------------------+
//           |
//           V
// +- Task : ------------+
// +  Input  data: JSON  +
// +  Output data: JSON  +
// +---------------------+
//
//
//
//
// Server
// console interface
//       V
//       | Output: JSON
//       | * ID of worker
//       | * Input data
//       | * ID job for analyze - It is done
//       |
//       V
// Client
// no interface
// Calculate
//       |
//       |
//       |
//

package research

import (
	"bufio"
	"fmt"
	"image/color"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

var researchFolder string

func init() {
	researchFolder = "Research"
	realFolder := "." + string(filepath.Separator) + researchFolder
	if _, err := os.Stat(realFolder); os.IsNotExist(err) {
		err := os.Mkdir(realFolder, 0777)
		if err != nil {
			panic(err)
		}
	}
}

func createResearchDir(dir string) {
	realFolder := "." + string(filepath.Separator) + researchFolder + string(filepath.Separator) + dir
	if _, err := os.Stat(realFolder); os.IsNotExist(err) {
		err := os.Mkdir(realFolder, 0777)
		if err != nil {
			panic(err)
		}
	}
}

func createResearchFilename(researchName string, file string) (fileName string, err error) {
	fileName = "." + string(filepath.Separator) + researchFolder + string(filepath.Separator) + researchName + string(filepath.Separator) + file
	if _, err := os.Stat(fileName); os.IsExist(err) {
		return fileName, err
	}
	return fileName, nil
}

func executeCalculix(file string) (err error) {
	// check file is exist
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return fmt.Errorf("Cannot execute : %v", err)
	}
	// check caclulix is exist
	ccx := "ccx"
	// remove .INP
	file = file[:(len(file) - 4)]
	// execute
	out, err := exec.Command(ccx, "-i", file).Output()
	if err != nil {
		return fmt.Errorf("Try install from https://pkgs.org/download/calculix-ccx\nError in calculix execution: %v\n%v", err, out)
	}
	return nil
}

func getBucklingFactor(file string, amountBuckling int) (bucklingFactor []float64, err error) {
	if strings.ToUpper(file[len(file)-3:]) != "DAT" {
		return bucklingFactor, fmt.Errorf("Wrong filename : %v", file)
	}
	// check file is exist
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return bucklingFactor, fmt.Errorf("Cannot find file : %v", err)
	}
	// found buckling header
	inFile, err := os.Open(file)
	if err != nil {
		return bucklingFactor, err
	}
	defer func() {
		errFile := inFile.Close()
		if errFile != nil {
			if err != nil {
				err = fmt.Errorf("%v ; %v", err, errFile)
			} else {
				err = errFile
			}
		}
	}()
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)
	bucklingHeader := "B U C K L I N G   F A C T O R   O U T P U T"
	var found bool
	var numberLine int
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if !found {
			// empty line
			if len(line) == 0 {
				continue
			}

			if len(line) != len(bucklingHeader) {
				continue
			}

			if line == bucklingHeader {
				found = true
			}
		} else {
			numberLine++
			if numberLine >= 5+amountBuckling {
				break
			}
			if numberLine >= 5 {
				m, f, err := parseBucklingFactor(line)
				if err != nil {
					return bucklingFactor, err
				}
				if m != numberLine-4 {
					return bucklingFactor, fmt.Errorf("Wrong MODE NO: %v (%v) in line: %v", m, numberLine-4, line)
				}
				bucklingFactor = append(bucklingFactor, f)
			}
		}
	}
	if len(bucklingFactor) != amountBuckling {
		return bucklingFactor, fmt.Errorf("Wrong lenght of buckling lines in DAT file")
	}
	return bucklingFactor, nil
}

// Example:
//      4   0.4067088E+03
func parseBucklingFactor(line string) (mode int, factor float64, err error) {
	s := strings.Split(line, "   ")
	for i := range s {
		s[i] = strings.TrimSpace(s[i])
	}
	i, err := strconv.ParseInt(s[0], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("Error: string parts - %v, error - %v, in line - %v", s, err, line)
	}
	mode = int(i)

	factor, err = strconv.ParseFloat(s[1], 64)
	if err != nil {

		return 0, 0, fmt.Errorf("Error: string parts - %v, error - %v, in line - %v", s, err, line)
	}
	return mode, factor, nil
}

// GetColor - index from 0 to 1
func GetColor(index float64) color.RGBA {
	colorShema := []color.RGBA{
		color.RGBA{R: 255, G: 0, B: 255, A: 255},   // Magenta
		color.RGBA{R: 0, G: 255, B: 0, A: 255},     // Green
		color.RGBA{R: 0, G: 0, B: 255, A: 255},     // Blue
		color.RGBA{R: 255, G: 0, B: 0, A: 255},     // Red
		color.RGBA{R: 100, G: 100, B: 100, A: 255}, // Grey
		color.RGBA{R: 255, G: 255, B: 0, A: 255},   // Yellow
	}
	if index <= 0. {
		index = 1e-5
	}
	if index > 1. {
		index = 1 - 1e-5
	}
	// color = 0...6
	index = index*float64(len(colorShema)-1) - 1e-5
	i := int(index)
	index = index - float64(i)

	result := color.RGBA{
		R: uint8(float64(colorShema[i+1].R)*index + float64(colorShema[i].R)*(1.-index)),
		G: uint8(float64(colorShema[i+1].G)*index + float64(colorShema[i].G)*(1.-index)),
		B: uint8(float64(colorShema[i+1].B)*index + float64(colorShema[i].B)*(1.-index)),
		A: uint8(float64(colorShema[i+1].A)*index + float64(colorShema[i].A)*(1.-index)),
	}
	return result
}
