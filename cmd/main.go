package main

import (
	"flag"
	"fmt"
	"os"

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
