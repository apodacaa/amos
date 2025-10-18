package ui

import (
	"fmt"
	"strings"
)

// RenderDashboard renders the main dashboard view
func RenderDashboard(width, height int) string {
	containerStyle := GetContainerStyle(width, height)
	titleStyle := GetTitleStyle(width)
	var menuItems []string

	menuItems = append(menuItems,
		menuItemStyle.Render(fmt.Sprintf("%s New Entry", keyStyle.Render("[n]"))),
	)
	menuItems = append(menuItems,
		menuItemStyle.Render(fmt.Sprintf("%s Todos", keyStyle.Render("[t]"))),
	)
	menuItems = append(menuItems,
		menuItemStyle.Render(fmt.Sprintf("%s Entries", keyStyle.Render("[e]"))),
	)
	menuItems = append(menuItems,
		menuItemStyle.Render(fmt.Sprintf("%s Search", keyStyle.Render("[s]"))),
	)

	menu := strings.Join(menuItems, "\n")
	help := helpStyle.Render("Press q to quit")

	content := titleStyle.Render("AMOS") + "\n\n" + menu + "\n" + help

	return containerStyle.Render(content)
}
