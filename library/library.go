package library

import "strings"

const (
	AQUAMAN         = "magnet:?xt=urn:btih:b6e82665ef588bb6574db1f9780a0279274f407d&dn=Aquaman.2018.1080p.WEBRip.x264-MP4&tr=udp%3A%2F%2Ftracker.leechers-paradise.org%3A6969&tr=udp%3A%2F%2Ftracker.openbittorrent.com%3A80&tr=udp%3A%2F%2Fopen.demonii.com%3A1337&tr=udp%3A%2F%2Ftracker.coppersurfer.tk%3A6969&tr=udp%3A%2F%2Fexodus.desync.com%3A6969"
	FANTASTICBEASTS = "magnet:?xt=urn:btih:426ec6d01964bac82c0da451b8e67842608fcc61&dn=Fantastic.Beasts.The.Crimes.Of.Grindelwald.2018.1080p.WEBRip.x26&tr=udp%3A%2F%2Ftracker.leechers-paradise.org%3A6969&tr=udp%3A%2F%2Ftracker.openbittorrent.com%3A80&tr=udp%3A%2F%2Fopen.demonii.com%3A1337&tr=udp%3A%2F%2Ftracker.coppersurfer.tk%3A6969&tr=udp%3A%2F%2Fexodus.desync.com%3A6969"
)

func URL(name string) string {
	movies := make(map[string]string)
	movies["Aquaman"] = "magnet:?xt=urn:btih:b6e82665ef588bb6574db1f9780a0279274f407d&dn=Aquaman.2018.1080p.WEBRip.x264-MP4&tr=udp%3A%2F%2Ftracker.leechers-paradise.org%3A6969&tr=udp%3A%2F%2Ftracker.openbittorrent.com%3A80&tr=udp%3A%2F%2Fopen.demonii.com%3A1337&tr=udp%3A%2F%2Ftracker.coppersurfer.tk%3A6969&tr=udp%3A%2F%2Fexodus.desync.com%3A6969"
	movies["Fantastic Beasts"] = "magnet:?xt=urn:btih:426ec6d01964bac82c0da451b8e67842608fcc61&dn=Fantastic.Beasts.The.Crimes.Of.Grindelwald.2018.1080p.WEBRip.x26&tr=udp%3A%2F%2Ftracker.leechers-paradise.org%3A6969&tr=udp%3A%2F%2Ftracker.openbittorrent.com%3A80&tr=udp%3A%2F%2Fopen.demonii.com%3A1337&tr=udp%3A%2F%2Ftracker.coppersurfer.tk%3A6969&tr=udp%3A%2F%2Fexodus.desync.com%3A6969"
	for k, v := range movies {
		if name == k || name == strings.ToLower(name) {
			return v
		}
	}
	return ""
}
