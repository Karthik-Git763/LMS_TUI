package models

import (
	"fmt"
	"os/exec"
	"runtime"
)

func OpenBrowser(url string) error {
	var cmd *exec.Cmd
	
	switch runtime.GOOS {
		case "linux":
			cmd = exec.Command("xdg-open", url)
		case "windows":
			cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
		case "darwin":
			cmd = exec.Command("open", url)
		default:
			fmt.Println("Unsupported operating system")
			return nil
	}
	return cmd.Run()
}