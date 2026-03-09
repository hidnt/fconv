# fconv

**fconv** is a CLI utility for converting media files (images, videos, and audio) from one format to another.

## Dependencies

For the utility to work, the following must be installed: **[FFmpeg](https://ffmpeg.org/)**

FFmpeg install:
* **Linux (Arch based):** `sudo pacman -S ffmpeg`
* **MacOS:** `brew install ffmpeg`
* **Windows:** `winget install ffmpeg` or [official site](https://ffmpeg.org/download.html).

## Install

```bash
go install github.com/hidnt/fconv@latest
```
*(Make sure your `$GOPATH/bin` is added to your `$PATH` environment variable)*

## Usage

Command syntax:
```bash
fconv [FILES...] --to <EXTENSION> [OPTIONS] [flags]
```

### Flags

| Flag | Full name | Description |
| :--- | :--- | :--- |
| `-t` | `--to string` | Target extension |
| `-o` | `--output string` | Destination folder for saving (default ".") |
| `-r` | `--recursive` | Recursive directory traversal |
| `-L` | `--level int` | Level of recursion |
| `-d` | `--delete` | Delete files after convertion |
| `-f` | `--force` | Overwrite the target file if it already exists |
| `-h` | `--help` | Help for fconv|

## Examples

**1. Video convert:**
```bash
fconv ./video.mkv --to mp4
```

**2. Recursively convert all `.png` images to `.webp`, removing the originals:**
```bash
fconv ./images --to webp -r -d
```

**3. Recursive traversal of a folder, but not deeper than 2 levels of nesting:**
```bash
fconv ./data --to jpg -r -L 2
```

## Supported extensions

### Images
- **From:** `png`, `jpg`, `jpeg`, `webp`, `bmp`, `tiff`, `avif`, `ico`, `cur`, `heic`, `heif`
- **To:** `png`, `jpg`, `jpeg`, `webp`, `bmp`, `tiff`, `avif`, `ico`, `cur`

### Video
- **From:** `mp4`, `mkv`, `mov`, `avi`, `gif`, `wmv`, `ogg`
- **To:** `mp4`, `mkv`, `mov`, `avi`, `gif`, `wmv`

### Music
- **From:** `mp3`, `wav`, `flac`, `aac`
- **To:** `mp3`, `wav`, `flac`, `aac`

### Video to music
- **From** `mp4`, `mkv`, `mov`, `avi`, `wmv`, `ogg`
- **To** `mp3`, `wav`, `flac`, `aac`