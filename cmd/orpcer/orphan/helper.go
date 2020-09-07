package orphan

import (
	"log"
	"os"
	"sort"
	"strconv"

	"github.com/olekukonko/tablewriter"
)

func printResults(results [][]string) {
	println()
	table := tablewriter.NewWriter(os.Stdout)
	table.NumLines()
	table.SetHeader([]string{"Domain", "AGE(days)", "CREATED AT"})
	table.SetBorder(true)     // Set Border to false
	table.AppendBulk(results) // Add Bulk Data
	table.Render()
}

func sortSlice(slice [][]string) [][]string {
	sort.SliceStable(slice, func(i, j int) bool {
		var err error
		var inta int
		var intb int
		inta, err = strconv.Atoi(slice[i][1])
		if err != nil {
			log.Fatal(err)
		}
		intb, err = strconv.Atoi(slice[j][1])
		if err != nil {
			log.Fatal(err)
		}
		if inta > intb {
			return true
		}
		return false
	})
	return slice
}
