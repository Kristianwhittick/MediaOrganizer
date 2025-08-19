package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/rwcarlsen/goexif/exif"
)

var metadataOnly bool

func main() {
	if len(os.Args) < 3 || len(os.Args) > 4 {
		usage()
		os.Exit(1)
	}

	workingDir := os.Args[1]
	outputDir := os.Args[2]
	
	if len(os.Args) == 4 && os.Args[3] == "-m" {
		metadataOnly = true
	}

	os.MkdirAll(outputDir, 0755)

	err := filepath.Walk(workingDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		// Skip hidden directories and files
		if strings.HasPrefix(info.Name(), ".") {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		
		if info.IsDir() {
			return nil
		}
		
		if isMediaFile(path) {
			return moveFile(path, outputDir)
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

func usage() {
	fmt.Println("Usage: main <source> <destination> [-m]")
	fmt.Println()
	fmt.Println("Where: <source> is the directory that contains many subfolder of files")
	fmt.Println("       <destination> is the output directory.")
	fmt.Println("       -m = metadata only mode (skip files without metadata date)")
}

func moveFile(sourceFile, outputDir string) error {
	name := filepath.Base(sourceFile)
	
	date := getDate(sourceFile)
	if date == nil {
		fmt.Printf("Skipping [%s] - no metadata date found\n", name)
		return nil
	}
	
	yearDir := filepath.Join(outputDir, date.Format("2006"))
	dateDir := filepath.Join(yearDir, date.Format("2006_01_02"))
	dest := filepath.Join(dateDir, name)

	if err := os.MkdirAll(dateDir, 0755); err != nil {
		return err
	}

	// Handle duplicate filenames
	if _, err := os.Stat(dest); err == nil {
		dest = getUniqueFilename(dateDir, name)
	}

	err := os.Rename(sourceFile, dest)
	if err != nil {
		return err
	}

	fmt.Printf("Moving [%s] to [%s]\n", name, dest)
	return nil
}

func isMediaFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".jpg" || ext == ".jpeg" || ext == ".mp4"
}

func getDate(filename string) *time.Time {
	ext := strings.ToLower(filepath.Ext(filename))
	
	// Try EXIF data for images
	if ext == ".jpg" || ext == ".jpeg" {
		if date := getExifDate(filename); date != nil {
			return date
		}
	}
	
	// Try MP4 metadata
	if ext == ".mp4" {
		if date := getMp4Date(filename); date != nil {
			return date
		}
	}
	
	// If metadata-only mode, return nil to skip
	if metadataOnly {
		return nil
	}
	
	// Fallback to file modification time
	if info, err := os.Stat(filename); err == nil {
		t := info.ModTime()
		return &t
	}
	
	t := time.Now()
	return &t
}

func getExifDate(filename string) *time.Time {
	file, err := os.Open(filename)
	if err != nil {
		return nil
	}
	defer file.Close()

	x, err := exif.Decode(file)
	if err != nil {
		return nil
	}

	tm, err := x.DateTime()
	if err != nil {
		return nil
	}

	return &tm
}

func getMp4Date(filename string) *time.Time {
	file, err := os.Open(filename)
	if err != nil {
		return nil
	}
	defer file.Close()

	// Look for mvhd atom which contains creation time
	buf := make([]byte, 8)
	for {
		n, err := file.Read(buf)
		if err != nil || n < 8 {
			break
		}
		
		size := binary.BigEndian.Uint32(buf[:4])
		atomType := string(buf[4:8])
		
		if atomType == "mvhd" {
			// Skip version and flags (4 bytes)
			file.Seek(4, 1)
			// Read creation time (4 bytes)
			timeBuf := make([]byte, 4)
			if n, err := file.Read(timeBuf); err == nil && n == 4 {
				// MP4 time is seconds since Jan 1, 1904
				mp4Time := binary.BigEndian.Uint32(timeBuf)
				// Convert to Unix time (subtract seconds between 1904 and 1970)
				unixTime := int64(mp4Time) - 2082844800
				if unixTime > 0 {
					t := time.Unix(unixTime, 0)
					return &t
				}
			}
			break
		}
		
		if size < 8 {
			break
		}
		file.Seek(int64(size-8), 1)
	}
	
	return nil
}

func getUniqueFilename(dir, name string) string {
	ext := filepath.Ext(name)
	base := strings.TrimSuffix(name, ext)
	
	for i := 0; ; i++ {
		var newName string
		if ext == "" {
			newName = name + strconv.Itoa(i)
		} else {
			newName = base + strconv.Itoa(i) + ext
		}
		
		newPath := filepath.Join(dir, newName)
		if _, err := os.Stat(newPath); os.IsNotExist(err) {
			return newPath
		}
	}
}