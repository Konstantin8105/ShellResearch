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
	"os"
	"path/filepath"
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
	err := os.Mkdir("."+string(filepath.Separator)+researchFolder+string(filepath.Separator)+dir, 0777)
	if err != nil {
		panic(err)
	}
}
