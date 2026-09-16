package utils

import (
	"bytes"
	"encoding/base64"
	"image/png"
	"strings"
	"testing"
)

func assertCaptchaImage(t *testing.T, dataURL string) {
	t.Helper()
	encoded, ok := strings.CutPrefix(dataURL, "data:image/png;base64,")
	if !ok {
		t.Fatal("captcha must be a PNG data URL")
	}
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		t.Fatalf("decode captcha base64: %v", err)
	}
	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("decode captcha PNG: %v", err)
	}
	if bounds := img.Bounds(); bounds.Dx() != 240 || bounds.Dy() != 80 {
		t.Fatalf("unexpected captcha size: %v", bounds)
	}
}

func TestGenerateCaptcha(t *testing.T) {
	id, image, answer, err := GenerateCaptcha()
	if err != nil {
		t.Fatalf("generate captcha: %v", err)
	}
	if len(id) != idLen || answer == "" {
		t.Fatalf("invalid captcha ID or empty answer: id=%q answer=%q", id, answer)
	}
	assertCaptchaImage(t, image)
}

func TestCaptchaGenerators(t *testing.T) {
	generators := map[string]generator{
		"digit":      &digitCaptcha{height: 80, width: 240, length: 5, dotCount: 80},
		"text":       &textCaptcha{height: 80, width: 240, length: 6, dotCount: 80, source: txtSimpleCharacters},
		"arithmetic": &arithmeticCaptcha{height: 80, width: 240, dotCount: 80},
	}
	for name, g := range generators {
		t.Run(name, func(t *testing.T) {
			id, content, answer := g.generate()
			if len(id) != idLen || content == "" || answer == "" {
				t.Fatalf("incomplete captcha: id=%q content=%q answer=%q", id, content, answer)
			}
			image, err := g.draw(content)
			if err != nil {
				t.Fatalf("draw captcha: %v", err)
			}
			assertCaptchaImage(t, image.EncodeB64string())
		})
	}
}
