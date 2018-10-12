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

// --------------------------------
// Author(URL) -----------thumbnail
// Title(URL)  -----------
// Description -----------
// Fileds      -----------
// ----Image -------------
// Fotter(image,text)*timestamp

func embed(title, url, desc, timestamp string, color int, author *discordgo.MessageEmbedAuthor, thumbnail *discordgo.MessageEmbedThumbnail, footer *discordgo.MessageEmbedFooter, image *discordgo.MessageEmbedImage) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		URL:         url,
		Title:       title,
		Description: desc,
		Timestamp:   timestamp,
		Color:       color,
		Footer:      footer,
		//image bottom
		Image: image,
		//image right up
		Thumbnail: thumbnail,
		Author:    author,
		Fields:    nil,
	}
}

func msgSend(content string, embed *discordgo.MessageEmbed) *discordgo.MessageSend {
	return &discordgo.MessageSend{
		Content: content,
		Embed:   embed,
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
