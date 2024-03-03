package types

type Link struct {
  Name string
  URL string
}

type CSS struct {
  Path string
}

type HeaderData struct {
  SocialLinks []Link
  Categories []Link
  CSS []string
  Colors Colors
}

type Colors struct {
  Background string
  White string
  Text string
  NeonRed string
  NeonBlue string
}

