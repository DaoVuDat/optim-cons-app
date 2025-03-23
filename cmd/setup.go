/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	listui "golang-moaha-construction/cmd/ui/list"
	multiinput "golang-moaha-construction/cmd/ui/multiInput"
	"golang-moaha-construction/internal/algorithms"
	"golang-moaha-construction/internal/algorithms/aha"
	"golang-moaha-construction/internal/algorithms/ga"
	"golang-moaha-construction/internal/algorithms/moaha"
	"golang-moaha-construction/internal/data"
	"golang-moaha-construction/internal/objectives/multi"
	"golang-moaha-construction/internal/objectives/single"
	"os"
	"time"
)

const logo = `


 ________  ________                             ________  ________  _________  ___  _____ ______      
|\   ____\|\   __  \                           |\   __  \|\   __  \|\___   ___\\  \|\   _ \  _   \    
\ \  \___|\ \  \|\  \        ____________      \ \  \|\  \ \  \|\  \|___ \  \_\ \  \ \  \\\__\ \  \   
 \ \  \  __\ \  \\\  \      |\____________\     \ \  \\\  \ \   ____\   \ \  \ \ \  \ \  \\|__| \  \  
  \ \  \|\  \ \  \\\  \     \|____________|      \ \  \\\  \ \  \___|    \ \  \ \ \  \ \  \    \ \  \ 
   \ \_______\ \_______\                          \ \_______\ \__\        \ \__\ \ \__\ \__\    \ \__\
    \|_______|\|_______|                           \|_______|\|__|         \|__|  \|__|\|__|     \|__|

`

var (
	startStyle     = lipgloss.NewStyle().Padding(1).Foreground(lipgloss.Color("#06FCC6")).Bold(true)
	groupStyle     = lipgloss.NewStyle().Padding(1).Foreground(lipgloss.Color("#01FAC6")).Bold(true)
	logoStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#01FAC6")).Bold(true)
	selectedOption = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("#F9F9F9")).Bold(true)
	errorStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#eb3477")).Bold(true)
)

// slightly red color = #FF0000

type Options struct {
	algorithm         *listui.Output
	algoConfigs       *multiinput.Output
	objectiveFunction *listui.Output
	objConfigs        *multiinput.Output
}

// setup represents the run command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "To configure the objective and an algorithm.",
	Long:  `The command setup the algorithm for specific objective function to optimize.`,
	Run: func(cmd *cobra.Command, args []string) {

		objConfigs := []*data.Config{
			{
				Name:  single.SphereDimension,
				Value: "30",
			},
			{
				Name:  single.SphereUpperBound,
				Value: "100",
			},
			{
				Name:  single.SphereLowerBound,
				Value: "-100",
			},
		}

		objectiveFunction, _ := single.CreateSphere(objConfigs)

		algoConfigs := []*data.Config{
			{
				Name:  aha.NUM_AGENTS,
				Value: "30",
			},
			{
				Name:  aha.NUM_ITERS,
				Value: "500",
			},
		}

		var algorithm algorithms.Algorithm

		algorithm, _ = aha.Create(objectiveFunction, algoConfigs)
		//err := algorithm.Run()
		//if err != nil {
		//	fmt.Println(err)
		//	os.Exit(1)
		//}

		algoConfigs = []*data.Config{
			{
				Name:  ga.PopulationSizeParam,
				Value: "30",
			},
			{
				Name:  ga.MaxIterationsParam,
				Value: "500",
			},
			{
				Name:  ga.CrossoverRateParam,
				Value: "0.8",
			},
			{
				Name:  ga.MutationRateParam,
				Value: "0.2",
			},
			{
				Name:  ga.ElitismCountParam,
				Value: "5",
			},
		}

		algorithm, _ = ga.Create(objectiveFunction, algoConfigs)
		//err = algorithm.Run()
		//if err != nil {
		//	fmt.Println(err)
		//	os.Exit(1)
		//}

		// Multi-objective
		multiConfigs := []*data.Config{
			{
				Name:  multi.ZDT1Dimension,
				Value: "30",
			},
			{
				Name:  multi.ZDT1UpperBound,
				Value: "1",
			},
			{
				Name:  multi.ZDT1LowerBound,
				Value: "0",
			},
		}

		obj, _ := multi.CreateZDT1(multiConfigs)

		algoConfigs = []*data.Config{
			{
				Name:  moaha.NumAgents,
				Value: "1000",
			},
			{
				Name:  moaha.NumIters,
				Value: "1000",
			},
			{
				Name:  moaha.ArchiveSize,
				Value: "300",
			},
		}

		algorithm, _ = moaha.Create(obj, algoConfigs)
		now := time.Now()
		err := algorithm.Run()
		fmt.Println(time.Since(now))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// show results
		moahaAlgo := algorithm.(*moaha.MOAHAAlgorithm)

		fmt.Println("Number Of Archived Agents: ", len(moahaAlgo.Archive))
		fmt.Println("Number Of Agents: ", len(moahaAlgo.Agents))
		fmt.Println("Number Of Iterations: ", moahaAlgo.NumberOfIter)
		fmt.Println("Number Of Dimension: ", len(moahaAlgo.Agents[0].Position))

		for idx := 0; idx < obj.NumberOfObjectives(); idx++ {

			for i, agent := range moahaAlgo.Archive {

				fmt.Printf("  %f", agent.Value[idx])
				if i != len(moahaAlgo.Archive)-1 {
					fmt.Printf(",")
				} else {
					fmt.Printf(";")
				}

			}
			fmt.Printf("\n\n\n")
		}

		return
		//// show logo
		//fmt.Printf("%s\n", logoStyle.Render(logo))
		//
		//// initialize options
		//options := Options{
		//	objectiveFunction: &listui.Output{},
		//	algorithm:         &listui.Output{},
		//	objConfigs: &multiinput.Output{
		//		Output: make([]*data.Config, 0),
		//	},
		//	algoConfigs: &multiinput.Output{
		//		Output: make([]*data.Config, 0),
		//	},
		//}
		//
		//// select algorithm
		//m := listui.InitializeList(
		//	options.algorithm,
		//	"Select algorithm",
		//	uidata.Algorithms,
		//	list.NewDefaultDelegate(), 20, 14)
		//p := tea.NewProgram(m, tea.WithAltScreen())
		//
		//if _, err := p.Run(); err != nil {
		//	fmt.Println("Error running program:", err)
		//	os.Exit(1)
		//}
		//
		//fmt.Printf("Algorithm: %s\n", selectedOption.Render(options.algorithm.Output.Title()))
		//
		//// select supported objectives function
		//m = listui.InitializeList(
		//	options.objectiveFunction,
		//	"Select objective function to optimize",
		//	uidata.Objectives,
		//	list.NewDefaultDelegate(), 20, 14)
		//
		//p = tea.NewProgram(m, tea.WithAltScreen())
		//
		//if _, err := p.Run(); err != nil {
		//	fmt.Println("Error running program:", err)
		//	os.Exit(1)
		//}
		//fmt.Printf("Objective function: %s\n", selectedOption.Render(options.objectiveFunction.Output.Title()))
		//
		//// config for algorithm
		//algoConfigs := uidata.GetAlgorithmConfigs(options.algorithm.Output.Title())
		//for _, conf := range algoConfigs {
		//	// convert to input field
		//	inputField := &data.Config{
		//		Name:               conf.Name,
		//		ValidationFunction: conf.ValidationFunction,
		//	}
		//	options.algoConfigs.Output = append(options.algoConfigs.Output, inputField)
		//}
		//
		//p = tea.NewProgram(multiinput.InitialModel(options.algoConfigs, "Configure Algorithm"), tea.WithAltScreen())
		//if _, err := p.Run(); err != nil {
		//	fmt.Println("Error running program:", err)
		//	os.Exit(1)
		//}
		//
		//fmt.Printf("%s\n", groupStyle.Render("Configured Algorithm"))
		//for _, conf := range options.algoConfigs.Output {
		//	fmt.Printf("   - %s: %s\n", conf.Name, selectedOption.Render(conf.Value))
		//}
		//
		//// config for objective
		//objConfigs := uidata.GetObjectiveConfigs(options.objectiveFunction.Output.Title())
		//for _, conf := range objConfigs {
		//	// convert to input field
		//	inputField := &data.Config{
		//		Name:               conf.Name,
		//		ValidationFunction: conf.ValidationFunction,
		//	}
		//	options.objConfigs.Output = append(options.objConfigs.Output, inputField)
		//}
		//
		//p = tea.NewProgram(multiinput.InitialModel(options.objConfigs, "Configure Objective"), tea.WithAltScreen())
		//if _, err := p.Run(); err != nil {
		//	fmt.Println("Error running program:", err)
		//	os.Exit(1)
		//}
		//
		//fmt.Printf("%s\n", groupStyle.Render("Configured Objective"))
		//
		//for _, conf := range options.objConfigs.Output {
		//	fmt.Printf("  - %s: %s\n", conf.Name, selectedOption.Render(conf.Value))
		//}
		//
		//// Create Objective
		//var problem objectives.Problem
		//var err error
		//switch options.objectiveFunction.Output.Title() {
		//case "sphere":
		//	problem, err = single.CreateSphere(options.objConfigs.Output)
		//	if err != nil {
		//		fmt.Println(errorStyle.Render("Error: objective function"))
		//		fmt.Println(errorStyle.Render(err.Error()))
		//		os.Exit(1)
		//	}
		//case "cons-lay":
		//	//problem = objectives.NewConsLay()
		//}
		//
		//// Create Algorithm
		//var algo algorithms.Algorithm
		//switch options.algorithm.Output.Title() {
		//case "GWO":
		//	algo, err = gwo.Create(problem, options.algoConfigs.Output)
		//	if err != nil {
		//		fmt.Println(errorStyle.Render("Error: algorithm"))
		//		fmt.Println(errorStyle.Render(err.Error()))
		//		os.Exit(1)
		//	}
		//case "MOAHA":
		//	//algo = algorithms.NewMOAHA()
		//}
		//
		//// check if objective function and algorithm are of the same type
		//if algo.Type() != problem.Type() {
		//	fmt.Println(errorStyle.Render("Error: objective function and algorithm must be of the same type"))
		//}
		//
		//fmt.Printf("%s\n", groupStyle.Render("Optimizing..."))
		//
		//err = algo.Run()
		//if err != nil {
		//	fmt.Println(errorStyle.Render(err.Error()))
		//	os.Exit(1)
		//}

	},
}

func init() {
	rootCmd.AddCommand(setupCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
