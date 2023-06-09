package application

import (
	"regexp"
	"strings"
)

type Content struct {
	Inner string
}

func newContent(content string) *Content {
	return &Content{content}
}

func (c *Content) getMode() string {
	index := strings.LastIndex(c.Inner, "(")
	if index == -1 {
		return ""
	}

	return c.Inner[index+1 : len(c.Inner)-1]
}

func (c *Content) getPrompt() string {
	re := `\*\*(.*?)\*\*`
	regexp, _ := regexp.Compile(re)
	match := regexp.FindStringSubmatch(c.Inner)
	return match[1]
}

func (c *Content) getProcessRate() string {
	index := strings.LastIndex(c.Inner, "(")
	if index == -1 {
		return ""
	}

	index2 := strings.LastIndex(c.Inner[0:index], "(")
	if index2 == -1 {
		return ""
	}

	return c.Inner[index2+1 : index-2]
}
