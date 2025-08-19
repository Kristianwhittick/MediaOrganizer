package organizer

import (
	"bytes"
	"encoding/binary"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/rwcarlsen/goexif/exif"
)

const mp4ToUnixEpochOffset = 2082844800 // Seconds between MP4 epoch (1904-01-01) and Unix epoch (1970-01-01)
const maxMP4ReadSize = 1024 * 1024      // Limit MP4 parsing to first 1MB

func IsMediaFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg", ".mp4":
		return true
	default:
		return false
	}
}

func GetDate(filename string, metadataOnly bool) *time.Time {
	ext := strings.ToLower(filepath.Ext(filename))
	
	if ext == ".jpg" || ext == ".jpeg" {
		if date := getExifDate(filename); date != nil {
			return date
		}
	}
	
	if ext == ".mp4" {
		if date := getMp4Date(filename); date != nil {
			return date
		}
	}
	
	if metadataOnly {
		return nil
	}
	
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

	buf := make([]byte, 4096)
	mvhdAtom := []byte("mvhd")
	totalRead := 0
	
	for totalRead < maxMP4ReadSize {
		n, err := file.Read(buf)
		if err != nil || n < 8 {
			break
		}
		totalRead += n
		
		for i := 0; i <= n-8; i++ {
			if i+20 <= n && bytes.Equal(buf[i+4:i+8], mvhdAtom) {
				size := binary.BigEndian.Uint32(buf[i:i+4])
				if size < 32 {
					continue
				}
				
				mp4Time := binary.BigEndian.Uint32(buf[i+16:i+20])
				unixTime := int64(mp4Time) - mp4ToUnixEpochOffset
				if unixTime > 0 {
					t := time.Unix(unixTime, 0)
					return &t
				}
				return nil
			}
		}
		
		if n == 4096 {
			if _, err := file.Seek(-7, 1); err != nil {
				break
			}
		}
	}
	
	return nil
}

func GetUniqueFilename(dir, name string) string {
	ext := filepath.Ext(name)
	base := strings.TrimSuffix(name, ext)
	
	for i := 1; ; i++ {
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