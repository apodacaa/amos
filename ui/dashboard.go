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
		fmt.Sprintf("%s %s", keyStyle.Render("[n]"), menuItemStyle.Render("New Entry")),
	)
	menuItems = append(menuItems,
		fmt.Sprintf("%s %s", keyStyle.Render("[a]"), menuItemStyle.Render("Add Todo")),
	)
	menuItems = append(menuItems,
		fmt.Sprintf("%s %s", keyStyle.Render("[t]"), menuItemStyle.Render("Todos")),
	)
	menuItems = append(menuItems,
		fmt.Sprintf("%s %s", keyStyle.Render("[e]"), menuItemStyle.Render("Entries")),
	)

	menu := strings.Join(menuItems, "\n")
	help := helpStyle.Render("Press q to quit")

	content := titleStyle.Render("AMOS") + "\n\n" + menu + "\n" + help

	return containerStyle.Render(content)
}
