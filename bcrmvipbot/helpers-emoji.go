// Package main provides ...
package main

import "regexp"

func emojiFromContent(content string) (id, name string) {
	exp := regexp.MustCompile(`<:(.+?):(\d{18})>`)
	result := exp.FindAllStringSubmatch(content, -1)
	if len(result) == 0 {
		return
	}
	id = result[0][2]
	name = result[0][1]
	return
}
