package utils

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"math"
	"math/big"
	"strings"
)

const (
	idLen               = 20
	txtNumbers          = "012346789"
	txtAlphabet         = "ABCDEFGHJKMNOQRSTUVXYZabcdefghjkmnoqrstuvxyz"
	txtSimpleCharacters = "13467ertyiadfhjkxcvbnERTYADFGHJKXCVBN"
	mimeTypeImage       = "image/png"
	glyphWidth          = 5
	glyphHeight         = 7
)

var (
	idChars = []byte(txtNumbers + txtAlphabet)
)

type generator interface {
	draw(content string) (imageItem, error)
	generate() (id, question, answer string)
}

type imageItem interface {
	WriteTo(w io.Writer) (int64, error)
	EncodeB64string() string
}

type Captcha struct{}

func NewCaptcha() *Captcha { return &Captcha{} }

func (c *Captcha) Generate() (id, base64Image, answer string, err error) {
	g := randomGenerator()
	id, content, answer := g.generate()
	captcha, err := g.draw(content)
	if err != nil {
		return "", "", "", err
	}
	return id, captcha.EncodeB64string(), answer, nil
}

func randomGenerator() generator {
	switch randomInt(3) {
	case 0:
		return &digitCaptcha{height: 80, width: 240, length: 5, dotCount: 80}
	case 1:
		return &textCaptcha{height: 80, width: 240, dotCount: 80, length: 6, source: txtSimpleCharacters}
	default:
		return &arithmeticCaptcha{height: 80, width: 240, dotCount: 80}
	}
}

type digitCaptcha struct {
	height, width, length, dotCount int
}

func (c *digitCaptcha) generate() (id, question, answer string) {
	id = randomID()
	answer = digitsToString(randomDigits(normalizeLength(c.length, 5)))
	return id, answer, answer
}

func (c *digitCaptcha) draw(content string) (imageItem, error) {
	w, h := normalizeSize(c.width, c.height)
	captcha := newCaptchaImage(w, h, randomLightColor())
	drawInterference(captcha, c.dotCount)
	if err := captcha.drawText(content); err != nil {
		return nil, err
	}
	return captcha, nil
}

type textCaptcha struct {
	background *color.RGBA
	source     string
	height     int
	width      int
	dotCount   int
	length     int
}

func (c *textCaptcha) generate() (id, content, answer string) {
	id = randomID()
	content = randomText(normalizeLength(c.length, 6), normalizeSource(c.source))
	return id, content, content
}

func (c *textCaptcha) draw(content string) (imageItem, error) {
	captcha := newCaptchaImageWithBackground(c.width, c.height, c.background)
	drawInterference(captcha, c.dotCount)
	if err := captcha.drawText(content); err != nil {
		return nil, err
	}
	return captcha, nil
}

type arithmeticCaptcha struct {
	background *color.RGBA
	height     int
	width      int
	dotCount   int
}

func (c *arithmeticCaptcha) generate() (id, question, answer string) {
	id = randomID()
	var result int32
	switch []string{"+", "-", "x"}[randomInt(3)] {
	case "+":
		a, b := int32(randomInt(20)), int32(randomInt(20))
		question, result = fmt.Sprintf("%d+%d=?", a, b), a+b
	case "x":
		a, b := int32(randomInt(10)), int32(randomInt(10))
		question, result = fmt.Sprintf("%dx%d=?", a, b), a*b
	default:
		a, b := int32(randomInt(80)+randomInt(20)), int32(randomInt(80))
		question, result = fmt.Sprintf("%d-%d=?", a, b), a-b
	}
	return id, question, fmt.Sprintf("%d", result)
}

func (c *arithmeticCaptcha) draw(question string) (imageItem, error) {
	captcha := newCaptchaImageWithBackground(c.width, c.height, c.background)
	drawInterference(captcha, c.dotCount)
	if err := captcha.drawText(question); err != nil {
		return nil, err
	}
	return captcha, nil
}

type captchaImage struct {
	nrgba  *image.NRGBA
	width  int
	height int
}

func newCaptchaImage(width, height int, background color.RGBA) *captchaImage {
	img := &captchaImage{width: width, height: height, nrgba: image.NewNRGBA(image.Rect(0, 0, width, height))}
	draw.Draw(img.nrgba, img.nrgba.Bounds(), &image.Uniform{C: background}, image.Point{}, draw.Src)
	return img
}

func newCaptchaImageWithBackground(width, height int, background *color.RGBA) *captchaImage {
	w, h := normalizeSize(width, height)
	if background != nil {
		return newCaptchaImage(w, h, *background)
	}
	return newCaptchaImage(w, h, randomLightColor())
}

func drawInterference(img *captchaImage, dotCount int) {
	img.drawSineLine()
	img.drawDots(maxInt(dotCount, 0))
}

func (img *captchaImage) drawSineLine() {
	a := randomInt(maxInt(img.height/2, 1))
	b := randomFloatRange(int64(-img.height/4), int64(img.height/4))
	f := randomFloatRange(int64(-img.height/4), int64(img.height/4))
	period := randomFloatRange(int64(maxInt(img.height, 1)), int64(maxInt(img.width/2, img.height+1)))
	w := (2 * math.Pi) / period
	c := randomDarkColor()
	for px, px2 := 0, int(randomFloatRange(int64(float64(img.width)*0.8), int64(img.width))); px < px2; px++ {
		py := float64(a)*math.Sin(w*float64(px)+f) + b + float64(img.width)/5
		for i := img.height / 5; i > 0; i-- {
			img.set(px+i, int(py), c)
		}
	}
}

func (img *captchaImage) drawDots(count int) {
	for i := 0; i < count; i++ {
		size := randomIntRange(1, maxInt(img.height/35, 2))
		img.drawBlock(randomInt(maxInt(img.width, 1)), randomInt(maxInt(img.height, 1)), size, size, randomDarkColor())
	}
}

func (img *captchaImage) drawText(text string) error {
	runes := []rune(text)
	if len(runes) == 0 {
		return errors.New("text must not be empty")
	}
	cellWidth := maxInt(img.width/len(runes), 1)
	for i, r := range runes {
		scale := minInt(maxInt(cellWidth/(glyphWidth+2), 1), maxInt(img.height/(glyphHeight+3), 1))
		if scale > 2 {
			scale -= randomInt(2)
		}
		x := i*cellWidth + maxInt((cellWidth-glyphWidth*scale)/2, 0) + randomIntRange(-scale, scale+1)
		y := maxInt((img.height-glyphHeight*scale)/2, 0) + randomIntRange(-scale, scale+1)
		img.drawRune(r, x, y, scale, randomDarkColor(), randomFloat64Range(-0.35, 0.35))
	}
	return nil
}

func (img *captchaImage) drawRune(r rune, x, y, scale int, c color.RGBA, skew float64) {
	for row, line := range glyphFor(r) {
		xOffset := int(float64(row-glyphHeight/2) * skew * float64(scale))
		for col, pixel := range line {
			if pixel == '1' {
				img.drawBlock(x+col*scale+xOffset, y+row*scale, scale, scale, c)
			}
		}
	}
}

func (img *captchaImage) drawBlock(x, y, w, h int, c color.RGBA) {
	for yy := 0; yy < h; yy++ {
		for xx := 0; xx < w; xx++ {
			img.set(x+xx, y+yy, c)
		}
	}
}

func (img *captchaImage) set(x, y int, c color.Color) {
	if x >= 0 && y >= 0 && x < img.width && y < img.height {
		img.nrgba.Set(x, y, c)
	}
}

func (img *captchaImage) binaryEncoding() []byte {
	var buf bytes.Buffer
	if err := png.Encode(&buf, img.nrgba); err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func (img *captchaImage) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(img.binaryEncoding())
	return int64(n), err
}

func (img *captchaImage) EncodeB64string() string {
	return fmt.Sprintf("data:%s;base64,%s", mimeTypeImage, base64.StdEncoding.EncodeToString(img.binaryEncoding()))
}

func normalizeLength(length, fallback int) int {
	if length > 0 {
		return length
	}
	return fallback
}

func normalizeSource(source string) string {
	if source != "" {
		return source
	}
	return txtSimpleCharacters
}

func normalizeSize(width, height int) (int, int) {
	if width <= 0 {
		width = 240
	}
	if height <= 0 {
		height = 80
	}
	return width, height
}

func digitsToString(b []byte) string {
	out := make([]byte, len(b))
	for i, n := range b {
		out[i] = n + '0'
	}
	return string(out)
}

func randomDigits(length int) []byte { return randomBytesMod(length, 10) }

func randomBytes(length int) []byte {
	b := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		panic("captcha: error reading random source: " + err.Error())
	}
	return b
}

func randomBytesMod(length int, mod byte) []byte {
	if length == 0 {
		return nil
	}
	if mod == 0 {
		panic("captcha: bad mod argument for randomBytesMod")
	}
	maxrb, out, i := 255-byte(256%int(mod)), make([]byte, length), 0
	for {
		for _, c := range randomBytes(length + length/4) {
			if c <= maxrb {
				out[i] = c % mod
				i++
				if i == length {
					return out
				}
			}
		}
	}
}

func randomText(size int, sourceChars string) string {
	if sourceChars == "" || size == 0 {
		return ""
	}
	if size >= len(sourceChars) {
		sourceChars = strings.Repeat(sourceChars, size)
	}
	source, text := []rune(sourceChars), make([]rune, size)
	for i := range text {
		text[i] = source[randomInt(len(source))]
	}
	return string(text)
}

func randomID() string {
	b := randomBytesMod(idLen, byte(len(idChars)))
	for i, c := range b {
		b[i] = idChars[c]
	}
	return string(b)
}

func randomFloatRange(min, max int64) float64 { return float64(min) + randomFloat64()*float64(max-min) }

func randomBaseColor() color.RGBA {
	red, green := randomInt(255), randomInt(255)
	blue := 400 - green - red
	if red+green > 400 {
		blue = 0
	} else if blue > 255 {
		blue = 255
	}
	return color.RGBA{R: uint8(red), G: uint8(green), B: uint8(blue), A: 255}
}

func randomDarkColor() color.RGBA {
	c, n := randomBaseColor(), float64(30+randomInt(255))
	return color.RGBA{R: uint8(math.Abs(math.Min(float64(c.R)-n, 255))), G: uint8(math.Abs(math.Min(float64(c.G)-n, 255))), B: uint8(math.Abs(math.Min(float64(c.B)-n, 255))), A: 255}
}

func randomLightColor() color.RGBA {
	return color.RGBA{R: uint8(randomInt(55) + 200), G: uint8(randomInt(55) + 200), B: uint8(randomInt(55) + 200), A: 255}
}

func randomIntRange(from, to int) int {
	if to-from <= 0 {
		return from
	}
	return randomInt(to-from) + from
}

func randomFloat64Range(from, to float64) float64 { return randomFloat64()*(to-from) + from }

func randomInt(max int) int {
	if max <= 0 {
		return 0
	}
	n, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		panic("captcha: error reading random source: " + err.Error())
	}
	return int(n.Int64())
}

func randomFloat64() float64 { return float64(randomInt(1<<53)) / (1 << 53) }

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func glyphFor(r rune) []string {
	if g, ok := glyphs[r]; ok {
		return g
	}
	if r >= 'a' && r <= 'z' {
		return glyphs[r-'a'+'A']
	}
	return glyphs['?']
}

var glyphs = map[rune][]string{
	'0': {"01110", "10001", "10011", "10101", "11001", "10001", "01110"}, '1': {"00100", "01100", "00100", "00100", "00100", "00100", "01110"}, '2': {"01110", "10001", "00001", "00010", "00100", "01000", "11111"},
	'3': {"11110", "00001", "00001", "01110", "00001", "00001", "11110"}, '4': {"00010", "00110", "01010", "10010", "11111", "00010", "00010"}, '5': {"11111", "10000", "10000", "11110", "00001", "00001", "11110"},
	'6': {"01110", "10000", "10000", "11110", "10001", "10001", "01110"}, '7': {"11111", "00001", "00010", "00100", "01000", "01000", "01000"}, '8': {"01110", "10001", "10001", "01110", "10001", "10001", "01110"},
	'9': {"01110", "10001", "10001", "01111", "00001", "00001", "01110"}, 'A': {"01110", "10001", "10001", "11111", "10001", "10001", "10001"}, 'B': {"11110", "10001", "10001", "11110", "10001", "10001", "11110"},
	'C': {"01111", "10000", "10000", "10000", "10000", "10000", "01111"}, 'D': {"11110", "10001", "10001", "10001", "10001", "10001", "11110"}, 'E': {"11111", "10000", "10000", "11110", "10000", "10000", "11111"},
	'F': {"11111", "10000", "10000", "11110", "10000", "10000", "10000"}, 'G': {"01111", "10000", "10000", "10011", "10001", "10001", "01111"}, 'H': {"10001", "10001", "10001", "11111", "10001", "10001", "10001"},
	'I': {"01110", "00100", "00100", "00100", "00100", "00100", "01110"}, 'J': {"00111", "00010", "00010", "00010", "00010", "10010", "01100"}, 'K': {"10001", "10010", "10100", "11000", "10100", "10010", "10001"},
	'L': {"10000", "10000", "10000", "10000", "10000", "10000", "11111"}, 'M': {"10001", "11011", "10101", "10101", "10001", "10001", "10001"}, 'N': {"10001", "11001", "10101", "10011", "10001", "10001", "10001"},
	'O': {"01110", "10001", "10001", "10001", "10001", "10001", "01110"}, 'P': {"11110", "10001", "10001", "11110", "10000", "10000", "10000"}, 'Q': {"01110", "10001", "10001", "10001", "10101", "10010", "01101"},
	'R': {"11110", "10001", "10001", "11110", "10100", "10010", "10001"}, 'S': {"01111", "10000", "10000", "01110", "00001", "00001", "11110"}, 'T': {"11111", "00100", "00100", "00100", "00100", "00100", "00100"},
	'U': {"10001", "10001", "10001", "10001", "10001", "10001", "01110"}, 'V': {"10001", "10001", "10001", "10001", "10001", "01010", "00100"}, 'W': {"10001", "10001", "10001", "10101", "10101", "10101", "01010"},
	'X': {"10001", "10001", "01010", "00100", "01010", "10001", "10001"}, 'Y': {"10001", "10001", "01010", "00100", "00100", "00100", "00100"}, 'Z': {"11111", "00001", "00010", "00100", "01000", "10000", "11111"},
	'a': {"00000", "00000", "01110", "00001", "01111", "10001", "01111"}, 'b': {"10000", "10000", "10110", "11001", "10001", "10001", "11110"}, 'c': {"00000", "00000", "01111", "10000", "10000", "10000", "01111"},
	'd': {"00001", "00001", "01101", "10011", "10001", "10001", "01111"}, 'e': {"00000", "00000", "01110", "10001", "11111", "10000", "01110"}, 'f': {"00110", "01001", "01000", "11100", "01000", "01000", "01000"},
	'g': {"00000", "01111", "10001", "10001", "01111", "00001", "01110"}, 'h': {"10000", "10000", "10110", "11001", "10001", "10001", "10001"}, 'i': {"00100", "00000", "01100", "00100", "00100", "00100", "01110"},
	'j': {"00010", "00000", "00110", "00010", "00010", "10010", "01100"}, 'k': {"10000", "10000", "10010", "10100", "11000", "10100", "10010"}, 'l': {"01100", "00100", "00100", "00100", "00100", "00100", "01110"},
	'm': {"00000", "00000", "11010", "10101", "10101", "10101", "10101"}, 'n': {"00000", "00000", "10110", "11001", "10001", "10001", "10001"}, 'o': {"00000", "00000", "01110", "10001", "10001", "10001", "01110"},
	'p': {"00000", "11110", "10001", "10001", "11110", "10000", "10000"}, 'q': {"00000", "01111", "10001", "10001", "01111", "00001", "00001"}, 'r': {"00000", "00000", "10110", "11001", "10000", "10000", "10000"},
	's': {"00000", "00000", "01111", "10000", "01110", "00001", "11110"}, 't': {"01000", "01000", "11100", "01000", "01000", "01001", "00110"}, 'u': {"00000", "00000", "10001", "10001", "10001", "10011", "01101"},
	'v': {"00000", "00000", "10001", "10001", "10001", "01010", "00100"}, 'w': {"00000", "00000", "10001", "10101", "10101", "10101", "01010"}, 'x': {"00000", "00000", "10001", "01010", "00100", "01010", "10001"},
	'y': {"00000", "10001", "10001", "10001", "01111", "00001", "01110"}, 'z': {"00000", "00000", "11111", "00010", "00100", "01000", "11111"}, '+': {"00000", "00100", "00100", "11111", "00100", "00100", "00000"},
	'-': {"00000", "00000", "00000", "11111", "00000", "00000", "00000"}, '=': {"00000", "00000", "11111", "00000", "11111", "00000", "00000"}, '?': {"01110", "10001", "00001", "00010", "00100", "00000", "00100"},
}
