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

 // amazonq-ignore-next-line
	if err := os.MkdirAll(cleanOutput, 0777); err != nil {
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
	sourceAbs, err := filepath.Abs(filepath.Clean(source))
	if err != nil {
		return fmt.Errorf("invalid source path: %w", err)
	}

	_, err = filepath.Abs(filepath.Clean(output))
	if err != nil {
		return fmt.Errorf("invalid output path: %w", err)
	}

	if _, err := os.Stat(sourceAbs); os.IsNotExist(err) {
		return fmt.Errorf("source directory does not exist: %s", sourceAbs)
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

 // amazonq-ignore-next-line
	yearDir := filepath.Join(outputDir, date.Format("2006"))
 // amazonq-ignore-next-line
	dateDir := filepath.Join(yearDir, date.Format("2006_01_02"))
 // amazonq-ignore-next-line
	dest := filepath.Join(dateDir, name)

 // amazonq-ignore-next-line
	if err := os.MkdirAll(dateDir, 0777); err != nil {
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
