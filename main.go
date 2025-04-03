package main

import (
	"fmt"
	"log"
	"os"
	"os/exec" // For executing commands
	"path/filepath"
	"strings"

	"github.com/gotk3/gotk3/gtk"
)

type Option struct {
	Name string
	Path string
}

func main() {
	// Initialize GTK
	gtk.Init(nil)

	// Create a new top-level window
	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}
	win.SetTitle("Poyo")
	win.SetDefaultSize(400, 200)
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	// Create a vertical box to hold UI elements
	vbox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 10)
	vbox.SetMarginTop(20)
	vbox.SetMarginBottom(20)
	vbox.SetMarginStart(20)
	vbox.SetMarginEnd(20)

	// Create an entry field for the username
	entry, _ := gtk.EntryNew()
	entry.SetPlaceholderText("Enter your username")
	selectionLabel, err := gtk.LabelNew("Current Selection: Nothing") // Create a label for the selection menu
	if err != nil {
		log.Fatal("Unable to create label:", err) // Handle error in label creation
	}
	// Create a password field
	password_field, err := gtk.EntryNew()
	if err != nil {
		log.Fatal("Unable to create entry:", err)
	}
	password_field.SetVisibility(false) // Hide the text to make it secret

	vbox.PackStart(entry, false, false, 0)          // Add the username field
	vbox.PackStart(password_field, false, false, 0) // Add the password field
	vbox.PackStart(selectionLabel, false, false, 0) // Add the label to show the current selection

	// Create a login button
	btn, _ := gtk.ButtonNewWithLabel("Login")
	btn.Connect("clicked", func() {
		// Get the text from the entry field when the button is clicked
		username, err := entry.GetText()
		if err != nil {
			log.Fatal("Error getting text from entry:", err)
		}
		fmt.Printf("Username: %s\n", username)
		// Here you can add your authentication logic
	})

	// Create a Menu for the pop-out selection
	menu, _ := gtk.MenuNew()

	// Create menu items
	f := findDesktops()
	for _, desktop := range f {
		// Create a local variable to store the value of 'desktop' to avoid closure issues
		desktopCopy := desktop

		entry, err := gtk.MenuItemNewWithLabel(desktopCopy.Name) // Create a new menu item with the session name
		if err != nil {
			log.Fatal("Unable to create menu item:", err)
		}

		// Connect the activate signal
		entry.Connect("activate", func() {
			selectionLabel.SetText(fmt.Sprintf("Current Selection: %s", desktopCopy.Name)) // Update the label with the selected session
			launchDesktop(desktopCopy)                                                     // Launch the selected app when the menu item is activated
		})

		menu.Add(entry) // Add the menu item to the menu
	} // Show the menu items
	menu.ShowAll()

	// Create an icon button in the top-right corner
	iconButton, _ := gtk.ButtonNew()
	iconImage, _ := gtk.ImageNewFromFile("icon.png") // Replace with the path to your icon image
	iconButton.SetImage(iconImage)

	// Position the icon button in the top-right corner using a fixed container
	fixed, _ := gtk.FixedNew()
	fixed.Put(iconButton, 350, 10) // Adjust position as needed
	vbox.PackStart(fixed, true, true, 0)

	// When the icon button is clicked, show the pop-up menu
	iconButton.Connect("clicked", func() {
		menu.PopupAtPointer(nil) // Show the menu at the current pointer position
	})

	// Add button to the vertical box
	vbox.PackStart(btn, false, false, 0)

	// Add the vertical box to the window
	win.Add(vbox)

	// Show all elements
	win.ShowAll()

	// Run the GTK main loop
	gtk.Main()
}

/*
// Function to launch the selected app

	func launchApp(app string) {
		var cmd *exec.Cmd

		// Select the appropriate command based on the app
		switch app {
		case "Plasma":
			cmd = exec.Command("startplasma-x11")
		case "Sway":
			cmd = exec.Command("sway")
		case "i3":
			cmd = exec.Command("i3")
		case "Openbox":
			cmd = exec.Command("openbox-session")
		default:
			fmt.Println("Unknown app:", app)
			return
		}

		// Run the selected command
		err := cmd.Start()
		if err != nil {
			fmt.Printf("Failed to launch %s: %s\n", app, err)
		} else {
			fmt.Printf("%s launched successfully!\n", app)
		}
	}
*/

func authCheck() {
}

func launchDesktop(desktop *Option) {
	cmd := parseDesktop(desktop)

	err := exec.Command(cmd).Start() // Start the command to launch the desktop environment
}

func parseDesktop(desktop *Option) string {
}

func findDesktops() []*Option {
	var desktops []*Option

	folder := "/usr/share/wayland-sessions/"

	// Read the files in the directory
	files, err := os.ReadDir(folder)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return desktops // Return empty if the directory can't be read
	}

	// Iterate through the files and process them
	for _, file := range files {
		if file.IsDir() {
			// Skip directories
			continue
		}

		// Get the file name
		fileName := file.Name()

		// Only process files that end with ".desktop"
		if strings.HasSuffix(fileName, ".desktop") {
			// Create a new Option struct for the session
			session := new(Option)

			// Remove the ".desktop" extension
			session.Name = fileName[:len(fileName)-8]

			// Capitalize the first letter of the session name
			session.Name = capitalizeFirstLetter(session.Name)

			// Get the absolute file path
			absPath, err := filepath.Abs(filepath.Join(folder, fileName))
			if err != nil {
				fmt.Println("Error getting absolute path:", err)
				continue
			}
			session.Path = absPath

			// Append the session to the desktops slice
			desktops = append(desktops, session)
		}
	}

	// Return the slice of Option structs
	return desktops
}

func capitalizeFirstLetter(word string) string {
	if len(word) == 0 {
		return word // Return the word unchanged if it's empty
	}
	return strings.ToUpper(string(word[0])) + word[1:]
}
