package main

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/godbus/dbus/v5"
)

// CREDIT: https://github.com/dawidd6/go-spotify-dbus
// These are internal methods so we can't import them :(

// Metadata contains Spotify player metadata
type Metadata struct {
	Artist      []string `spotify:"xesam:artist"`
	Title       string   `spotify:"xesam:title"`
	Album       string   `spotify:"xesam:album"`
	AlbumArtist []string `spotify:"xesam:albumArtist"`
	AutoRating  float64  `spotify:"xesam:autoRating"`
	DiskNumber  int32    `spotify:"xesam:discNumber"`
	TrackNumber uint32   `spotify:"xesam:trackNumber"`
	URL         string   `spotify:"xesam:url"`
	TrackID     string   `spotify:"mpris:trackid"`
	Length      int64    `spotify:"mpris:length"`
}

// parseMetadata returns a parsed Metadata struct
func parseMetadata(variant dbus.Variant) *Metadata {
	metadataMap := variant.Value().(map[string]dbus.Variant)
	metadataStruct := new(Metadata)

	valueOf := reflect.ValueOf(metadataStruct).Elem()
	typeOf := reflect.TypeOf(metadataStruct).Elem()

	for key, val := range metadataMap {
		for i := 0; i < typeOf.NumField(); i++ {
			field := typeOf.Field(i)
			if field.Tag.Get("spotify") == key {
				field := valueOf.Field(i)
				field.Set(reflect.ValueOf(val.Value()))
			}
		}
	}

	return metadataStruct
}

var (
	playingIcon = ""
	pausedIcon  = ""
)

func main() {
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	spotifyBus := conn.Object("org.mpris.MediaPlayer2.spotifyd", "/org/mpris/MediaPlayer2")
	playbackStatus, err := spotifyBus.GetProperty("org.mpris.MediaPlayer2.Player.PlaybackStatus")
	if err != nil {
		panic(err)
	}
	metadataProp, err := spotifyBus.GetProperty("org.mpris.MediaPlayer2.Player.Metadata")
	if err != nil {
		panic(err)
	}
	metadata := parseMetadata(metadataProp)
	if metadata.Title == "" {
		return
	}
	displayIcon := pausedIcon
	if playbackStatus.String() == "\"Playing\"" {
		displayIcon = playingIcon
	}
	fmt.Printf(
		"%s %s: %s",
		displayIcon,
		strings.Join(metadata.Artist[:1], ", "),
		metadata.Title,
	)
}
