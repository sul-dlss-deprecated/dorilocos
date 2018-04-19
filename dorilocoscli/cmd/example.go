package cmd

import (
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generator",
	Long:  "Generator",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		fmt.Printf("Using example template: %s\n", exampleFile)
		fmt.Printf("Number of examples to generate: %d\n", numToGenerate)
	},
}

var numToGenerate int
var exampleFile string

func init() {
	flag.IntVar(&numToGenerate, "n", 20, "Total number of examples to generate")
	flag.StringVar(&exampleFile, "s", "sdr3-object.json", "File to use as example template.")
}

// Execute the defined function above
func Execute() {
	if err := generateCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func shuffle(src []string) []string {
	final := make([]string, len(src))
	rand.Seed(time.Now().UTC().UnixNano())
	perm := rand.Perm(len(src))

	for i, v := range perm {
		final[v] = src[i]
	}
	return final
}

func copyFile(srcFil string, count int) {
	from, err := os.Open("./sourcefile.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer from.Close()

	to, err := os.OpenFile("./sourcefile.copy.txt", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer to.Close()

	_, err = io.Copy(to, from)
	if err != nil {
		log.Fatal(err)
	}
}
