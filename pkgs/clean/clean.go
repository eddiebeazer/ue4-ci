package clean

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
)

// Dist Deletes the given directory
func Dist(directory string) error {
	if directory == "" {
		return cli.Exit("No directory provided", 1)
	}
	err := os.RemoveAll(directory)
	if err != nil {
		return cli.Exit(fmt.Errorf("error deleting directory: %s", err), 1)
	}
	return nil
}
