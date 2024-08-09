package model

import (
	"image/color"
	"johtotimes.com/src/assert"
	"time"

	"github.com/afocus/captcha"
	"johtotimes.com/src/constants"
)

type Captcha struct {
	UUID      string
	Value     string
	Image     *captcha.Image
	CreatedAt time.Time
}

func NewCaptcha(captchaID string) Captcha {
	c := captcha.New()
	c.SetSize(256, 64)
	c.SetDisturbance(captcha.HIGH)
	// White font color
	c.SetFrontColor(color.White)
	// Transparent background with a different accent color
	c.SetBkgColor(
		color.RGBA{R: 255}, // transparent
		color.RGBA{B: 255}, // blue
		color.RGBA{G: 153}, // green
	)
	err := c.SetFont(constants.AssetPath + "/fonts/Annon.ttf")
	assert.NoError(err, "Captcha: Error setting font")
	img, str := c.Create(6, captcha.UPPER)
	return Captcha{
		UUID:  captchaID,
		Value: str,
		Image: img,
	}
}
