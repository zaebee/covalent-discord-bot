// Package model keeps structures for meme domain and elasticsearch helper enities.
package model

import "time"

// Meme implements meme image.
type Meme struct {
	ID        string
	Url       string
	ImageHash string
	Author    Author
	Timestamp time.Time
	Reactions []string
}

// Author implements meme author.
type Author struct {
	ID       string
	Username string
	Avatar   string
}

// ImageHash implements unique hash for image content.
type ImageHash struct {
	Url  string `json:"url"`
	Hash string `json:"hash"`
}
