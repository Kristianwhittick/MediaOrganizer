# MediaOrganizer

A fast, cross-platform media file organizer that automatically sorts photos and videos by date into a hierarchical folder structure.

## Features

- **Smart Date Detection**: Extracts creation dates from EXIF data (JPEG/JPG) and MP4 metadata
- **Hierarchical Organization**: Creates `YYYY/YYYY_MM_DD/` folder structure for easy browsing
- **Fast File Operations**: Uses filesystem move operations (no copying) for speed
- **Duplicate Handling**: Automatically renames duplicate files with numeric suffixes
- **Metadata-Only Mode**: Option to skip files without extractable metadata dates
- **Hidden File Filtering**: Ignores hidden directories and files (starting with `.`)
- **Cross-Platform**: Supports Windows, Linux, and ARM architectures

## Supported File Types

- **Images**: JPEG, JPG (with EXIF date extraction)
- **Videos**: MP4 (with metadata date extraction)

## Installation

### Download Pre-built Binaries

Download the appropriate binary for your platform from the releases page:

- `mediaorganizer-windows-x64.exe` - Windows 64-bit
- `mediaorganizer-linux-x64` - Linux 64-bit
- `mediaorganizer-linux-arm` - Linux ARM (ARMv5+)

### Build from Source

```bash
# Clone the repository
git clone <repository-url>
cd MediaOrganizer

# Install dependencies
go mod tidy

# Build for current platform
go build -o mediaorganizer ./cmd

# Or build for all platforms
make
```

## Usage

### Basic Usage

```bash
# Organize media files
./mediaorganizer /path/to/source /path/to/destination

# Windows
mediaorganizer.exe C:\Photos C:\Organized

# Show help
./mediaorganizer --help

# Show version
./mediaorganizer --version
```

### Metadata-Only Mode

```bash
# Only process files with extractable metadata dates
./mediaorganizer -m /path/to/source /path/to/destination
```

### Command Line Options

- `<source>` - Source directory containing media files
- `<destination>` - Output directory for organized files
- `-m` - Metadata-only mode (skip files without metadata dates)
- `--help` - Show help information
- `--version` - Show version information

## Output Structure

Files are organized into a hierarchical structure:

```
destination/
├── 2023/
│   ├── 2023_01_15/
│   │   ├── IMG_001.jpg
│   │   └── VIDEO_001.mp4
│   └── 2023_12_25/
│       └── holiday_photo.jpg
└── 2024/
    └── 2024_03_10/
        └── spring_video.mp4
```

## Date Detection Priority

1. **EXIF Data** (JPEG/JPG files) - Original photo timestamp
2. **MP4 Metadata** (MP4 files) - Video creation timestamp
3. **File Modification Time** (fallback, unless `-m` flag is used)

## Examples

### Organize All Media Files

```bash
./mediaorganizer ./unsorted ./organized
```

### Process Only Files with Metadata

```bash
./mediaorganizer -m ./photos ./sorted
```

### Cross-Platform Paths

```bash
# Linux/macOS
./mediaorganizer /home/user/photos /home/user/organized

# Windows
mediaorganizer.exe "C:\Users\Name\Pictures" "D:\Organized Photos"

# Relative paths
./mediaorganizer ../photos ./sorted
```

## Building

### Prerequisites

- Go 1.21 or later
- Make (optional, for cross-compilation)

### Build Commands

```bash
# Build all platforms
make

# Build specific platforms
make build-windows  # Windows x64
make build-linux    # Linux x64
make build-arm      # Linux ARM

# Clean build artifacts
make clean

# Install dependencies
make deps

# Run tests
go test ./organizer
```

### Manual Cross-Compilation

```bash
# Windows x64
GOOS=windows GOARCH=amd64 go build -o mediaorganizer-windows.exe ./cmd

# Linux x64
GOOS=linux GOARCH=amd64 go build -o mediaorganizer-linux ./cmd

# Linux ARM
GOOS=linux GOARCH=arm GOARM=5 go build -o mediaorganizer-arm ./cmd
```

## Dependencies

- [github.com/rwcarlsen/goexif](https://github.com/rwcarlsen/goexif) - EXIF data extraction

## Performance

- **Fast**: Uses `os.Rename()` for same-filesystem moves (no data copying)
- **Memory Efficient**: Processes files one at a time with optimized MP4 parsing
- **Scalable**: Handles thousands of files efficiently
- **Secure**: Path validation prevents directory traversal attacks

## Limitations

- Only processes JPEG/JPG and MP4 files
- MP4 metadata extraction is basic (mvhd atom only)
- Cross-filesystem moves will copy data (slower)
- Hidden files and directories are skipped

## License

See LICENSE file for details.

## Contributing

Contributions welcome! Please feel free to submit issues and pull requests.