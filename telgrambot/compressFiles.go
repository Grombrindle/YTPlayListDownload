package telgrambot

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

func CompressWith7z(sourceFile, destArchive string) error {
	sevenZipPath := `C:\Program Files\7-Zip\7z.exe`
	cmd := exec.Command(sevenZipPath, "a", "-t7z", "-mx=9", destArchive, sourceFile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("7z compression failed: %v, output: %s", err, string(output))
	}
	return nil
}

func LowerBitRate(inputFile, outputFile string) error {
	ffmpegPath := `C:\Users\Damasco\Downloads\ffempg\ffmpeg-2025-06-04-git-a4c1a5b084-full_build\bin\ffmpeg.exe`
	cmd := exec.Command(ffmpegPath, "-i", inputFile, "-b:a", "192k", outputFile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg failed: %v, output: %s", err, string(output))
	}
	return nil
}

func GenerateOutputPaths(originalPath string) (lowerBitratePath, compressedPath string) {
	baseDir := filepath.Dir(originalPath)
	baseName := filepath.Base(originalPath)
	ext := filepath.Ext(baseName)
	fileNameOnly := strings.TrimSuffix(baseName, ext)

	lowerBitrateDir := filepath.Join(baseDir, "lower_bitrate")
	compressedDir := filepath.Join(baseDir, "compressed")

	lowerBitratePath = filepath.Join(lowerBitrateDir, fileNameOnly+"_128k"+ext)
	compressedPath = filepath.Join(compressedDir, fileNameOnly+"_128k.7z")

	return
}
