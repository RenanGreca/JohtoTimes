package comment

import "time"

type Captcha struct {
	ID        int64
	UUID      string
	Value     string
	Image     string
	CreatedAt time.Time
}
