// Package main provides ...
package main

import "github.com/bwmarrin/discordgo"

func embedFooter(text, iconURL string) *discordgo.MessageEmbedFooter {
	return &discordgo.MessageEmbedFooter{
		Text:         text,
		IconURL:      iconURL,
		ProxyIconURL: "",
	}
}

func embedImage(url string) *discordgo.MessageEmbedImage {
	return &discordgo.MessageEmbedImage{
		URL:      url,
		ProxyURL: "",
		Width:    0,
		Height:   0,
	}
}

func embedThumbnail(url string) *discordgo.MessageEmbedThumbnail {
	return &discordgo.MessageEmbedThumbnail{
		URL:      url,
		ProxyURL: "",
		Width:    0,
		Height:   0,
	}
}

func embedAuthor(name, url, iconURL string) *discordgo.MessageEmbedAuthor {
	return &discordgo.MessageEmbedAuthor{
		URL:          url,
		Name:         name,
		IconURL:      iconURL,
		ProxyIconURL: "",
	}
}

func embedFields(size int) []*discordgo.MessageEmbedField {
	return make([]*discordgo.MessageEmbedField, 0, size)
}

type embedInfo struct {
	title     string
	url       string
	desc      string
	timestamp string
	color     int
}

// --------------------------------
// Author(URL) -----------thumbnail
// Title(URL)  -----------
// Description -----------
// Fileds      -----------
// ----Image -------------
// Fotter(image,text)*timestamp

func embed(info *embedInfo, author *discordgo.MessageEmbedAuthor, thumbnail *discordgo.MessageEmbedThumbnail, footer *discordgo.MessageEmbedFooter, image *discordgo.MessageEmbedImage, fields []*discordgo.MessageEmbedField) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		URL:         info.url,
		Title:       info.title,
		Description: info.desc,
		Timestamp:   info.timestamp,
		Color:       info.color,
		Footer:      footer,
		//image bottom
		Image: image,
		//image right up
		Thumbnail: thumbnail,
		Author:    author,
		Fields:    fields,
	}
}

func mef(n string, v string, i bool) *discordgo.MessageEmbedField {
	f := &discordgo.MessageEmbedField{
		Name:   n,
		Value:  v,
		Inline: i,
	}
	return f
}
