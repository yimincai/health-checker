package utils

import (
	"bytes"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"time"

	"github.com/golang/freetype"
	"github.com/yimincai/health-checker/fonts"
	"golang.org/x/image/font"
)

func TimeFormat(t time.Time) string {
	months := []string{"01", "02", "03", "04", "05", "06", "07", "08", "09", "10", "11", "12"}
	weekdays := []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}

	chineseFormat := fmt.Sprintf("%d/%s/%d %s %02d:%02d:%02d",
		t.Year(),
		months[t.Month()-1],
		t.Day(),
		weekdays[t.Weekday()],
		t.Hour(),
		t.Minute(),
		t.Second())

	return chineseFormat
}

func MonthToInt(m time.Month) int {
	// Map month to integer representation
	monthMap := map[time.Month]int{
		time.January:   1,
		time.February:  2,
		time.March:     3,
		time.April:     4,
		time.May:       5,
		time.June:      6,
		time.July:      7,
		time.August:    8,
		time.September: 9,
		time.October:   10,
		time.November:  11,
		time.December:  12,
	}

	// Retrieve integer representation from the map
	if val, ok := monthMap[m]; ok {
		return val
	}

	panic("Invalid month")
}

func IsVaildateDate(y, m, d int) bool {
	//check year
	if y < 1 || y > 9999 {
		return false
	}

	//check month
	if m < 1 || m > 12 {
		return false
	}

	//check day
	if d < 1 || d > 31 {
		return false
	}

	//check day in month
	if d > 30 && (m == 4 || m == 6 || m == 9 || m == 11) {
		return false
	}

	if m == 2 {
		if y%400 == 0 || (y%4 == 0 && y%100 != 0) {
			if d > 29 {
				return false
			}
		} else {
			if d > 28 {
				return false
			}
		}
	}

	return true
}

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
