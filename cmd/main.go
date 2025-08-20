package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"mediaorganizer/organizer"
)

const version = "1.0.0"

func main() {
	var (
		metadataOnly = flag.Bool("m", false, "metadata only mode (skip files without metadata date)")
		showVersion  = flag.Bool("version", false, "show version information")
		showHelp     = flag.Bool("help", false, "show help information")
	)

	flag.Usage = usage
	flag.Parse()

	// amazonq-ignore-next-line
	if *showVersion {
		fmt.Printf("MediaOrganizer v%s\n", version)
		return
	}

	// amazonq-ignore-next-line
	if *showHelp {
		usage()
		return
	}

	args := flag.Args()
	if len(args) != 2 {
		usage()
		os.Exit(1)
	}

	sourceDir := args[0]
	outputDir := args[1]

	if strings.TrimSpace(sourceDir) == "" {
		fmt.Fprintf(os.Stderr, "Error: source directory cannot be empty\n")
		os.Exit(1)
	}
	if strings.TrimSpace(outputDir) == "" {
		fmt.Fprintf(os.Stderr, "Error: destination directory cannot be empty\n")
		os.Exit(1)
	}

	if info, err := os.Stat(sourceDir); err != nil {
		if os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Error: source directory does not exist: %s\n", sourceDir)
		} else {
			fmt.Fprintf(os.Stderr, "Error: cannot access source directory: %v\n", err)
		}
		os.Exit(1)
	} else if !info.IsDir() {
		fmt.Fprintf(os.Stderr, "Error: source path is not a directory: %s\n", sourceDir)
		os.Exit(1)
	}

	if parentDir := filepath.Dir(outputDir); parentDir != "." {
		if _, err := os.Stat(parentDir); err != nil && os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Error: destination parent directory does not exist: %s\n", parentDir)
			os.Exit(1)
		}
	}

	// amazonq-ignore-next-line
	org := organizer.New(*metadataOnly)
	if err := org.OrganizeFiles(sourceDir, outputDir); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "MediaOrganizer v%s - Organize media files by date\n\n", version)
	fmt.Fprintf(os.Stderr, "Usage: %s [options] <source> <destination>\n\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Arguments:\n")
	fmt.Fprintf(os.Stderr, "  <source>       Source directory containing media files\n")
	fmt.Fprintf(os.Stderr, "  <destination>  Output directory for organized files\n\n")
	fmt.Fprintf(os.Stderr, "Options:\n")
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "\nSupported formats: JPEG, JPG, MP4\n")
}
