package ui

import (
	"fmt"
	"strings"

	"github.com/apodacaa/amos/internal/helpers"
	"github.com/charmbracelet/lipgloss"
)

// RenderLineGraph renders a full-screen line graph of weekly activity
func RenderLineGraph(weekStats []helpers.WeekStats, width, height int) string {
	if len(weekStats) == 0 {
		return ""
	}

	// Styles
	axisStyle := lipgloss.NewStyle().Foreground(mutedColor)
	textStyle := lipgloss.NewStyle().Foreground(subtleColor)

	// Find max value for scaling
	maxValue := 0
	for _, ws := range weekStats {
		if ws.EntryCount > maxValue {
			maxValue = ws.EntryCount
		}
		if ws.TodoCount > maxValue {
			maxValue = ws.TodoCount
		}
	}

	if maxValue == 0 {
		maxValue = 1 // Avoid division by zero
	}

	// Graph dimensions (leave room for title, legend, axes)
	graphHeight := height - 6 // Title + legend + spacing
	if graphHeight < 10 {
		graphHeight = 10
	}
	graphWidth := width - 8 // Y-axis labels

	// Create grid
	grid := make([][]rune, graphHeight)
	for i := range grid {
		grid[i] = make([]rune, graphWidth)
		for j := range grid[i] {
			grid[i][j] = ' '
		}
	}

	// Plot entries and todos as scatter plot (dots only, no connecting lines)
	dataPoints := max(len(weekStats)-1, 1)

	// Plot entries first (white dots)
	for i, ws := range weekStats {
		x := int(float64(i) * float64(graphWidth-1) / float64(dataPoints))
		y := graphHeight - 1 - int(float64(ws.EntryCount)*float64(graphHeight-1)/float64(maxValue))
		if x >= 0 && x < graphWidth && y >= 0 && y < graphHeight {
			grid[y][x] = '●'
		}
	}

	// Plot todos (cyan dots, or diamond if overlapping)
	for i, ws := range weekStats {
		x := int(float64(i) * float64(graphWidth-1) / float64(dataPoints))
		y := graphHeight - 1 - int(float64(ws.TodoCount)*float64(graphHeight-1)/float64(maxValue))
		if x >= 0 && x < graphWidth && y >= 0 && y < graphHeight {
			if grid[y][x] == '●' {
				grid[y][x] = '◆' // Both overlap
			} else {
				grid[y][x] = '○'
			}
		}
	}

	// Build output
	var lines []string

	// Render grid with Y-axis
	yStep := maxValue / 5
	if yStep == 0 {
		yStep = 1
	}

	for i := 0; i < graphHeight; i++ {
		// Y-axis label
		yValue := maxValue - (i * maxValue / max(graphHeight-1, 1))
		yLabel := fmt.Sprintf("%4d", yValue)
		yLabelInterval := max(graphHeight/5, 1)
		if i%yLabelInterval != 0 && i != 0 {
			yLabel = "    " // Only show labels every 1/5th
		}

		// Render grid line with monochrome styling
		gridLine := string(grid[i])

		lines = append(lines, axisStyle.Render(yLabel)+" "+axisStyle.Render("┤")+" "+textStyle.Render(gridLine))
	}

	// X-axis (extend to cover all labels)
	xAxis := "     " + axisStyle.Render("└"+strings.Repeat("─", graphWidth+4))
	lines = append(lines, xAxis)

	// X-axis labels positioned under their data points
	labelLine := make([]rune, graphWidth+10) // +6 for Y-axis space "     ┤" + extra for rightmost label
	for i := range labelLine {
		labelLine[i] = ' '
	}

	// Position each week label under its data point
	for i, ws := range weekStats {
		x := int(float64(i) * float64(graphWidth-1) / float64(dataPoints))
		labelPos := x + 6 // Offset for Y-axis labels

		// Write week label at calculated position
		for j, char := range ws.WeekLabel {
			if labelPos+j < len(labelLine) {
				labelLine[labelPos+j] = char
			}
		}
	}

	lines = append(lines, axisStyle.Render(string(labelLine)))

	// Legend with monochrome styling
	lines = append(lines, "")
	lines = append(lines, textStyle.Render("Entries ●  Todos ○"))

	return strings.Join(lines, "\n")
}
