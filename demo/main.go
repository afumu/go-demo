package main

import "fmt"

type Menu struct {
	ID       int
	ParentID int
	Name     string
	SubMenus []*Menu
}

func BuildMenu(menus []Menu) []*Menu {
	// Create a map to store the menu items by their ID
	menuMap := make(map[int]*Menu)
	for _, menu := range menus {
		menuMap[menu.ID] = &menu
	}

	// Iterate through the menu items and link the submenus
	for _, menu := range menus {
		if menu.ParentID != 0 {
			parentMenu := menuMap[menu.ParentID]
			parentMenu.SubMenus = append(parentMenu.SubMenus, &menu)
		}
	}

	// Find the top-level menu items and return them
	var topLevelMenus []*Menu
	for _, menu := range menus {
		if menu.ParentID == 0 {
			topLevelMenus = append(topLevelMenus, &menu)
		}
	}
	return topLevelMenus
}

func main() {
	menus := []Menu{
		{ID: 1, ParentID: 0, Name: "Menu 1"},
		{ID: 2, ParentID: 0, Name: "Menu 2"},
		{ID: 3, ParentID: 1, Name: "Menu 1.1"},
		{ID: 4, ParentID: 1, Name: "Menu 1.2"},
		{ID: 5, ParentID: 2, Name: "Menu 2.1"},
	}
	topLevelMenus := BuildMenu(menus)
	fmt.Println(topLevelMenus)

}
