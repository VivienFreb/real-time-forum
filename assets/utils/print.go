package utils

import (
	"fmt"
	"os"

	s "real/assets/struct"
)

// Print the Error with title, error message and code error for stop the solution (-1 to don't close)
func PrintError(title string, err error, status int) {
	fmt.Println(s.Color_Red, title, s.ResetAll)
	fmt.Println(s.ColorFontRGB(255, 100, 100), err, s.ResetAll)
	if status >= 0 {
		os.Exit(status)
	}
}
