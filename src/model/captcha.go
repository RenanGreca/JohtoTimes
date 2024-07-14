package model

import (
	"image/color"
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
	cap := captcha.New()
	cap.SetSize(256, 64)
	cap.SetDisturbance(captcha.HIGH)
	// White font color
	cap.SetFrontColor(color.White)
	// Transparent background with a different accent color
	cap.SetBkgColor(
		color.RGBA{255, 0, 0, 0}, // transparent
		color.RGBA{0, 0, 255, 0}, // blue
		color.RGBA{0, 153, 0, 0}, // green
	)
	cap.SetFont(constants.AssetPath + "/fonts/Annon.ttf")
	img, str := cap.Create(6, captcha.UPPER)
	captcha := Captcha{
		UUID:  captchaID,
		Value: str,
		Image: img,
	}

	return captcha
}
