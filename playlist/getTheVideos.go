package playlist

import (
	TB "YTPlayListDownload/telgrambot"
	"context"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	// "FreeGames/epicGames"
	// "github.com/gen2brain/beeep"
	//    h "github.com/go-rod/rod"
	// "github.com/go-rod/rod/lib/input"
	// "github.com/teocci/go-chrome-cookies/chrome"
	// un "github.com/Davincible/chromedp-undetected"
	// cu "github.com/lrakai/chromedp-undetected"
	// "github.com/chromedp/cdproto/cdp"
)

var (
	// spanText    string
	// videoNumber int
	// links       []string
	// customTitle []string
	// hrefs       []map[string]string
	// res         []byte
	outputPath       string = "C:/Users/Damasco/Music/goYTProgram"
	archivePath      string = "C:/Users/Damasco/Music/goYTProgram/archive.txt"
	lowerBitRatePath string
	compressedPath   string
)

func OpenLink() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	chromeProfilePath := filepath.Join(
		os.Getenv("LOCALAPPDATA"),
		"Google",
		"Chrome",
		"User Data",
		"Default",
	)

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.UserDataDir(chromeProfilePath),
		chromedp.WindowSize(1280, 800),
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-blink-features", "AutomationControlled"),
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	err := chromedp.Run(ctx,
		chromedp.Navigate("https://www.youtube.com/playlist?list=PLM08h2YldzVKcGVfO45R7IuC35SPY8L3Q"),
		chromedp.Sleep(3*time.Second),
	)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		err := chromedp.Run(ctx,
			chromedp.Evaluate(`window.scrollTo(0, document.documentElement.scrollHeight)`, nil),
			chromedp.Sleep(1*time.Second),
		)
		if err != nil {
			log.Fatal(err)
		}
	}
	time.Sleep(5 * time.Second)
	var links []string
	err = chromedp.Run(ctx,
		// chromedp.WaitVisible(`ytd-playlist-video-renderer`, chromedp.ByQuery),
		chromedp.Evaluate(`
			Array.from(document.querySelectorAll('ytd-playlist-video-renderer a#video-title')).map(a => a.href)`, &links),
	)
	if err != nil {
		log.Fatal(err)
	}

	var titles []string
	err = chromedp.Run(ctx,
		chromedp.Evaluate(`
			Array.from(document.querySelectorAll('ytd-playlist-video-renderer a#video-title')).map(a => a.title)`, &titles),
	)
	if err != nil {
		log.Fatal(err)
	}

	if len(links) != len(titles) {
		log.Fatalf("Mismatch between links (%d) and titles (%d)", len(links), len(titles))
	}

	TB.Init()
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	chatID := os.Getenv("CHAT_ID")
	if botToken == "" || chatID == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN and CHAT_ID must be set")
	}
	for i, link := range links {

		videoID := link[strings.Index(link, "v=")+2:]
		if idx := strings.Index(videoID, "&"); idx != -1 {
			videoID = videoID[:idx]
		}

		title := titles[i]
		if title == "" {
			title = videoID
		}

		filePath2, downloaded, err := DownloadAudio(link, outputPath, title, archivePath)
		if err != nil {
			log.Printf("Failed to download %s: %v\n", title, err)
			continue
		}
		if !downloaded {
			log.Printf("Skipped %s: already downloaded", title)
			continue
		}

		lowerBitRatePath, compressedPath = TB.GenerateOutputPaths(filePath2)

		if err := os.MkdirAll(filepath.Dir(lowerBitRatePath), 0755); err != nil {
			log.Fatalf("Failed to create directory for lower bitrate file: %v", err)
		}
		if err := os.MkdirAll(filepath.Dir(compressedPath), 0755); err != nil {
			log.Fatalf("Failed to create directory for compressed file: %v", err)
		}

		err = TB.LowerBitRate(filePath2, lowerBitRatePath)
		if err != nil {
			log.Fatal(err)
		}

		err = TB.CompressWith7z(lowerBitRatePath, compressedPath)
		if err != nil {
			log.Fatal(err)
		}

		err = TB.SendFileToTelegram(botToken, chatID, compressedPath)
		if err != nil {
			log.Printf("Failed to send file: %v\n", err)
		}

	}

	cmd := exec.Command("taskkill", "/IM", "chrome.exe", "/F")
	cmd.Run()

}
