package utils

import (
	"bytes"
	"flag"
	"image"
	"image/color"
	"image/draw"
	"image/png"

	"github.com/golang/freetype"
	"github.com/yimincai/health-checker/fonts"
	"golang.org/x/image/font"
)

var (
	dpi = flag.Float64("dpi", 72, "screen resolution in Dots Per Inch")
	// fontfile = flag.String("fontfile", "font.ttf", "filename of the ttf font")
	fontBytes = fonts.DefaultFont
	hinting   = flag.String("hinting", "none", "none | full")
	size      = flag.Float64("size", 18, "font size in points")
	spacing   = flag.Float64("spacing", 2, "line spacing (e.g. 2 means double spaced)")
	wonb      = flag.Bool("whiteonblack", true, "white text on a black background")
)

// var text = []string{
// 	"┌─────────────────────┬────────────┬────────────────────────────────────┬──────────┐",
// 	"│ ID                  │ NAME       │ URL                                │ INTERVAL │",
// 	"├─────────────────────┼────────────┼────────────────────────────────────┼──────────┤",
// 	"│ 1789957222001586176 │ service_1  │ https://api.example.com/api/health │ 60 s     │",
// 	"│ 1789957391636017152 │ service_2  │ https://example.com/               │ 60 s     │",
// 	"└─────────────────────┴────────────┴────────────────────────────────────┴──────────┘",
// }

func GenerateImage(text []string) (*bytes.Reader, error) {
	flag.Parse()

	// // Read the font data.
	// fontBytes, err := os.ReadFile(*fontfile)
	// if err != nil {
	// 	return nil, err
	// }
	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return nil, err
	}

	// Initialize the context.
	fg, bg := image.Black, image.White
	ruler := color.RGBA{0xdd, 0xdd, 0xdd, 0xff}
	if *wonb {
		fg, bg = image.White, image.Black
		ruler = color.RGBA{0x22, 0x22, 0x22, 0xff}
	}
	rgba := image.NewRGBA(image.Rect(0, 0, 640*2, 480*2))
	draw.Draw(rgba, rgba.Bounds(), bg, image.Point{}, draw.Src)
	c := freetype.NewContext()
	c.SetDPI(*dpi)
	c.SetFont(f)
	c.SetFontSize(*size)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(fg)
	switch *hinting {
	default:
		c.SetHinting(font.HintingNone)
	case "full":
		c.SetHinting(font.HintingFull)
	}

	// Draw the guidelines.
	for i := 0; i < 200; i++ {
		rgba.Set(10, 10+i, ruler)
		rgba.Set(10+i, 10, ruler)
	}

	// Draw the text.
	pt := freetype.Pt(10, 10+int(c.PointToFixed(*size)>>6))
	for _, s := range text {
		_, err = c.DrawString(s, pt)
		if err != nil {
			return nil, err
		}
		pt.Y += c.PointToFixed(*size * *spacing)
	}

	var buf bytes.Buffer
	err = png.Encode(&buf, rgba)
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(buf.Bytes())

	return reader, nil
}
