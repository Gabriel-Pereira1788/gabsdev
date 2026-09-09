// Package banner renders a per-article Open Graph banner (1200x630) that
// mirrors the site's own terminal-window look: same dark palette, same
// pixel-font title treatment, same monospace footer. It exists so every
// article gets a distinctive link-preview image (WhatsApp, Slack, Twitter,
// iMessage, ...) without anyone having to hand-design one per post.
package banner

import (
	"fmt"
	"image/color"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/fogleman/gg"

	"gabsdev-go/internal/domain"
	"gabsdev-go/internal/i18n"
)

const (
	width     = 1200
	height    = 630
	margin    = 40.0
	titlebarH = 56.0

	maxTitleLines = 4
	minTitlePt    = 24.0
	maxTitlePt    = 44.0

	pixelFontPath = "static/fonts/PressStart2P-Regular.ttf"
	monoFontPath  = "static/fonts/JetBrainsMono.ttf"
)

// Palette mirrors the --bg-1/--border/--text-* tokens in static/css/global.css
// (dark theme), since the generated banner is a screenshot-shaped stand-in
// for the site chrome rather than a themeable page.
var (
	colorBg       = hex("#0b0b0b")
	colorBorder   = hex("#2e2e2e")
	colorTitlebar = hex("#131313")
	colorText0    = hex("#ffffff")
	colorText1    = hex("#a6a6a6")
	colorText2    = hex("#6e6e6e")
)

func hex(h string) color.Color {
	var r, g, b int
	fmt.Sscanf(h, "#%02x%02x%02x", &r, &g, &b)
	return color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 255}
}

// Dir returns the on-disk directory generated banners for lang are written to.
func Dir(lang i18n.Lang) string {
	return filepath.Join("static", "images", "banners", string(lang))
}

// URLPath returns the public path (served by the static file handler) for a
// post's generated banner.
func URLPath(lang i18n.Lang, slug string) string {
	return "/static/images/banners/" + string(lang) + "/" + slug + ".png"
}

// GenerateAll renders a banner PNG for every post of every language listed in
// postsByLang, overwriting any file already on disk so an edited title is
// reflected on the next restart/deploy. Rendering is CPU-bound (font
// rasterization, no I/O beyond the final PNG write) and each post is
// independent, so work fans out across a bounded pool sized to the host's
// core count: wall time stays roughly constant as the number of posts grows
// instead of scaling linearly with a single-threaded loop.
func GenerateAll(postsByLang map[i18n.Lang][]domain.Post) error {
	type job struct {
		lang i18n.Lang
		post domain.Post
	}

	var jobs []job
	for lang, posts := range postsByLang {
		if err := os.MkdirAll(Dir(lang), 0o755); err != nil {
			return fmt.Errorf("banner: create dir %s: %w", Dir(lang), err)
		}
		for _, post := range posts {
			jobs = append(jobs, job{lang: lang, post: post})
		}
	}
	if len(jobs) == 0 {
		return nil
	}

	workers := runtime.NumCPU()
	if workers > len(jobs) {
		workers = len(jobs)
	}

	jobCh := make(chan job)
	errCh := make(chan error, len(jobs))
	var wg sync.WaitGroup
	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range jobCh {
				out := filepath.Join(Dir(j.lang), j.post.Slug+".png")
				if err := render(j.post, out); err != nil {
					errCh <- fmt.Errorf("banner: render %s/%s: %w", j.lang, j.post.Slug, err)
				}
			}
		}()
	}
	for _, j := range jobs {
		jobCh <- j
	}
	close(jobCh)
	wg.Wait()
	close(errCh)

	for err := range errCh {
		return err
	}
	return nil
}

func render(post domain.Post, outPath string) error {
	dc := gg.NewContext(width, height)

	dc.SetColor(colorBg)
	dc.Clear()

	// outer frame, echoes .term-window
	dc.SetColor(colorBorder)
	dc.SetLineWidth(2)
	dc.DrawRectangle(margin, margin, width-2*margin, height-2*margin)
	dc.Stroke()

	// titlebar strip, echoes .term-titlebar
	dc.SetColor(colorTitlebar)
	dc.DrawRectangle(margin, margin, width-2*margin, titlebarH)
	dc.Fill()
	dc.SetColor(colorBorder)
	dc.SetLineWidth(1)
	dc.DrawLine(margin, margin+titlebarH, width-margin, margin+titlebarH)
	dc.Stroke()

	// three outline dots + one filled, echoes .term-dot(.filled)
	dotY := margin + titlebarH/2
	dotX := margin + 28.0
	for i := range 3 {
		dc.SetColor(colorText1)
		dc.SetLineWidth(1.5)
		dc.DrawCircle(dotX+float64(i)*18, dotY, 5)
		if i == 2 {
			dc.Fill()
		} else {
			dc.Stroke()
		}
	}

	if err := dc.LoadFontFace(monoFontPath, 14); err != nil {
		return err
	}
	dc.SetColor(colorText1)
	dc.DrawStringAnchored("gabsdev / terminal blog", width-margin-20, dotY, 1, 0.35)

	contentX := margin + 48.0
	contentW := float64(width) - 2*margin - 96

	pt, lines, err := fitTitle(dc, post.Title, contentW)
	if err != nil {
		return err
	}
	lineH := pt * 1.5
	totalH := float64(len(lines)) * lineH
	availTop := margin + titlebarH
	availBottom := height - margin - 70.0
	startY := availTop + (availBottom-availTop-totalH)/2 + lineH/2

	dc.SetColor(colorText0)
	for i, line := range lines {
		dc.DrawStringAnchored(line, contentX, startY+float64(i)*lineH, 0, 0.35)
	}

	if err := dc.LoadFontFace(monoFontPath, 16); err != nil {
		return err
	}
	dc.SetColor(colorText2)
	dc.DrawString(strings.Join(post.Tags, " · "), contentX, height-margin-32)

	return dc.SavePNG(outPath)
}

// fitTitle picks the largest pixel-font size (in [minTitlePt, maxTitlePt])
// that wraps the title into at most maxTitleLines lines at contentW; long
// titles fall back to minTitlePt regardless of resulting line count.
func fitTitle(dc *gg.Context, title string, contentW float64) (float64, []string, error) {
	for pt := maxTitlePt; pt >= minTitlePt; pt -= 2 {
		if err := dc.LoadFontFace(pixelFontPath, pt); err != nil {
			return 0, nil, err
		}
		lines := dc.WordWrap(title, contentW)
		if len(lines) <= maxTitleLines {
			return pt, lines, nil
		}
	}
	if err := dc.LoadFontFace(pixelFontPath, minTitlePt); err != nil {
		return 0, nil, err
	}
	return minTitlePt, dc.WordWrap(title, contentW), nil
}
