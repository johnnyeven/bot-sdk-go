package template

import (
	"github.com/johnnyeven/bot-sdk-go/bot/directive"
)

const (
	PLAIN_TEXT = "PlainText"
	RICH_TEXT  = "RichText"
)

var textTypeMap = map[string]bool{
	PLAIN_TEXT: true,
	RICH_TEXT:  true,
}

type BaseTemplate struct {
	directive.BaseDirective
	Token           string `json:"token"`
	BackgroundImage *Image `json:"backgroundImage,omitempty"`
	Title           string `json:"title"`
}

func (t *BaseTemplate) SetTitle(title string) {
	t.Title = title
}

func (t *BaseTemplate) SetBackgroundImageUrl(url string) {
	t.SetBackgroundImage(NewImage(url))
}

func (t *BaseTemplate) SetBackgroundImage(background *Image) {
	t.BackgroundImage = background
}

type Image struct {
	Url          string `json:"url"`
	WidthPixels  int    `json:"widthPixels,omitempty"`
	HeightPixels int    `json:"heightPixels,omitempty"`
}

func NewImage(url string) *Image {
	image := &Image{}
	image.Url = url
	return image
}

func (i *Image) SetWidth(width int) *Image {
	i.WidthPixels = width
	return i
}

func (i *Image) SetHeight(height int) *Image {
	i.HeightPixels = height
	return i
}

type Text struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func NewText(textType, text string) *Text {
	ok := textTypeMap[textType]
	if !ok {
		textType = PLAIN_TEXT
	}

	t := &Text{}
	t.Type = textType
	t.Text = text
	return t
}
