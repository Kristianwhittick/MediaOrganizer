package organizer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Organizer struct {
	MetadataOnly bool
}

func New(metadataOnly bool) *Organizer {
	return &Organizer{MetadataOnly: metadataOnly}
}

func (o *Organizer) OrganizeFiles(sourceDir, outputDir string) error {
	cleanSource := filepath.Clean(sourceDir)
	cleanOutput := filepath.Clean(outputDir)

	if err := o.validatePaths(cleanSource, cleanOutput); err != nil {
		return err
	}

	if err := os.MkdirAll(cleanOutput, 0700); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	return filepath.Walk(cleanSource, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasPrefix(info.Name(), ".") {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if info.IsDir() {
			return nil
		}

		if IsMediaFile(path) {
			return o.moveFile(path, cleanOutput)
		}
		return nil
	})
}

func (o *Organizer) validatePaths(source, output string) error {
	if strings.Contains(source, "..") || strings.Contains(output, "..") {
		return fmt.Errorf("path traversal not allowed")
	}

	if !filepath.IsAbs(source) {
		abs, err := filepath.Abs(source)
		if err != nil {
			return fmt.Errorf("invalid source path: %w", err)
		}
		source = abs
	}

	if !filepath.IsAbs(output) {
		abs, err := filepath.Abs(output)
		if err != nil {
			return fmt.Errorf("invalid output path: %w", err)
		}
		output = abs
	}

	if _, err := os.Stat(source); os.IsNotExist(err) {
		return fmt.Errorf("source directory does not exist: %s", source)
	}

	return nil
}

func (o *Organizer) moveFile(sourceFile, outputDir string) error {
	name := filepath.Base(sourceFile)

	date := GetDate(sourceFile, o.MetadataOnly)
	if date == nil {
		fmt.Printf("Skipping [%s] - no metadata date found\n", name)
		return nil
	}

	yearDir := filepath.Join(outputDir, date.Format("2006"))
	dateDir := filepath.Join(yearDir, date.Format("2006_01_02"))
	dest := filepath.Join(dateDir, name)

	if err := os.MkdirAll(dateDir, 0700); err != nil {
		return fmt.Errorf("failed to create date directory: %w", err)
	}

	if _, err := os.Stat(dest); err == nil {
		dest = GetUniqueFilename(dateDir, name)
	}

	fmt.Printf("Moving [%s] to [%s]\n", name, dest)
	if err := os.Rename(sourceFile, dest); err != nil {
		return fmt.Errorf("failed to move file: %w", err)
	}

	return nil
}
