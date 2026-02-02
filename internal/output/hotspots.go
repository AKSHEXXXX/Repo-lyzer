package output

import (
	"fmt"
	"os"

	"github.com/agnivo988/Repo-lyzer/internal/analyzer"
	"github.com/olekukonko/tablewriter"
)

// PrintHotspots prints the identify hotspots in a table format
func PrintHotspots(hotspots []analyzer.Hotspot) {
	if len(hotspots) == 0 {
		return
	}

	fmt.Println("\n🔥 Hotspot files (Complex & Frequently Changed)")

	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"File", "Score", "Churn", "Size", "Complexity", "Issues"})

	for _, h := range hotspots {
		table.Append([]string{
			h.FilePath,
			fmt.Sprintf("%d", h.Score),
			fmt.Sprintf("%d", h.ChurnScore),
			fmt.Sprintf("%d", h.SizeScore),
			fmt.Sprintf("%d", h.Complexity),
			h.Reason,
		})
	}

	table.Render()
}
