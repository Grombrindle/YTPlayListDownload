package playlist

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	// "FreeGames/epicGames"
	"github.com/chromedp/chromedp"
	// "github.com/gen2brain/beeep"
	//    h "github.com/go-rod/rod"
	// "github.com/go-rod/rod/lib/input"
	// "github.com/teocci/go-chrome-cookies/chrome"
	// un "github.com/Davincible/chromedp-undetected"
	// cu "github.com/lrakai/chromedp-undetected"
	// "github.com/chromedp/cdproto/cdp"
)

func DownloadAudio(url, outputPath, customTitle, archivePath string) (string, bool, error) {
	videoID := url[strings.Index(url, "v=")+2:]
	if idx := strings.Index(videoID, "&"); idx != -1 {
		videoID = videoID[:idx]
	}

	safeTitle := strings.Map(func(r rune) rune {
		switch r {
		case '/', '\\', ':', '*', '?', '"', '<', '>', '|':
			return '-'
		default:
			return r
		}
	}, customTitle)

	filename := fmt.Sprintf("%s [%s]", safeTitle, videoID)

	cmd := exec.Command("yt-dlp",
		"--ffmpeg-location", `C:\Users\Damasco\Downloads\ffempg\ffmpeg-2025-06-04-git-a4c1a5b084-full_build\bin`,
		"-x",
		"--audio-format", "mp3",
		"-o", fmt.Sprintf("%s/%s.%%(ext)s", outputPath, filename),
		"--no-playlist",
		"--download-archive", archivePath,
		url,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", false, fmt.Errorf("yt-dlp failed: %v, output: %s", err, string(output))
	}

	outStr := string(output)
	if strings.Contains(outStr, "Skipping video") {
		return "", false, nil
	}
	if strings.Contains(outStr, "Destination:") && strings.Contains(outStr, "100%") {
		filePath := filepath.Join(outputPath, filename+".mp3")
		return filePath, true, nil
	}
	return "", false, nil
}

func ExtractNumber(xpath string, ctx context.Context) (int, error) {
	var spanText string
	err := chromedp.Run(ctx,
		// chromedp.WaitVisible(xpath, chromedp.BySearch),
		chromedp.TextContent(xpath, &spanText, chromedp.BySearch),
	)
	if err != nil {
		return 0, fmt.Errorf("chromedp run error: %w", err)
	}

	fmt.Println("Raw extracted text:", spanText)

	re := regexp.MustCompile(`(\d+)\s+video[s]?`)
	matches := re.FindStringSubmatch(spanText)
	if len(matches) < 2 {
		return 0, fmt.Errorf("number not found in text: %q", spanText)
	}

	numberStr := matches[1]
	number, err := strconv.Atoi(numberStr)
	if err != nil {
		return 0, fmt.Errorf("conversion error: %w", err)
	}

	return number, nil
}
