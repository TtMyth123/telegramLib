package kit

import (
	"github.com/zelenin/go-tdlib/client"
	"path/filepath"
	"strings"
)

type FileType int

const (
	FT_gif FileType = 1
	FT_png FileType = 2
	FT_mp4 FileType = 3
	FT_mp3 FileType = 4
	FT_o   FileType = 0
)

func GetFileType(path string) (FileType, string) {
	ext := filepath.Ext(path)
	ext = strings.ToLower(ext)
	switch ext {
	case ".mp4":
		return FT_mp4, ext
	case ".mp3":
		return FT_mp3, ext
	case ".gif":
		return FT_gif, ext
	case ".png", ".jpg":
		return FT_png, ext
	}
	return FT_o, ext
}
func GetStr(str *string) string {
	if str == nil {
		return ""
	} else {
		return *str
	}
}

func GetStr2I(str string) *string {
	if str == "" {
		return nil
	} else {
		return &str
	}
}

func FormattedText(str string) *client.FormattedText {
	if str == "" {
		return nil
	}
	aFormattedText := client.FormattedText{
		Text: str,
	}
	return &aFormattedText
}

func GetFormattedText(FormattedText *client.FormattedText) string {
	if FormattedText == nil {
		return ""
	} else {
		return FormattedText.Text
	}
}
