package main

import (
	"fmt"

	"github.com/Konstantin8105/CalculixRPCclient/clientCalculix"
	"github.com/Konstantin8105/ShellResearch/research"
)

func main() {
	inp, err := research.ShellModel(3., 5., 10, 10, 1., 0.005)
	if err != nil {
		fmt.Println("Cannot mesh")
		return
	}
	var client clientCalculix.ClientCalculix
	client.Manager = *clientCalculix.NewServerManager()
	factors, err := client.CalculateForBuckle([]string{inp})
	for i, factor := range factors {
		fmt.Println("Factor ", i, " -- ", factor)
	}
	research.RC002()
}
