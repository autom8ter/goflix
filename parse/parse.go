package parse

import (
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var patterns = []struct {
	name string
	// Use the last matching pattern. E.g. Year.
	last bool
	kind reflect.Kind
	// REs need to have 2 sub expressions (groups), the first one is "raw", and
	// the second one for the "clean" value.
	// E.g. Epiode matching on "S01E18" will result in: raw = "E18", clean = "18".
	re *regexp.Regexp
}{
	{"season", false, reflect.Int, regexp.MustCompile(`(?i)(s?([0-9]{1,2}))[ex]`)},
	{"episode", false, reflect.Int, regexp.MustCompile(`(?i)([ex]([0-9]{2})(?:[^0-9]|$))`)},
	{"episode", false, reflect.Int, regexp.MustCompile(`(-\s+([0-9]{1,})(?:[^0-9]|$))`)},
	{"year", true, reflect.Int, regexp.MustCompile(`\b(((?:19[0-9]|20[0-9])[0-9]))\b`)},

	{"resolution", false, reflect.String, regexp.MustCompile(`\b(([0-9]{3,4}p))\b`)},
	{"quality", false, reflect.String, regexp.MustCompile(`(?i)\b(((?:PPV\.)?[HP]DTV|(?:HD)?CAM|B[DR]Rip|(?:HD-?)?TS|(?:PPV )?WEB-?DL(?: DVDRip)?|HDRip|DVDRip|DVDRIP|CamRip|W[EB]BRip|BluRay|DvDScr|telesync))\b`)},
	{"codec", false, reflect.String, regexp.MustCompile(`(?i)\b((xvid|[hx]\.?26[45]))\b`)},
	{"audio", false, reflect.String, regexp.MustCompile(`(?i)\b((MP3|DD5\.?1|Dual[\- ]Audio|LiNE|DTS|AAC[.-]LC|AAC(?:\.?2\.0)?|AC3(?:\.5\.1)?))\b`)},
	{"region", false, reflect.String, regexp.MustCompile(`(?i)\b(R([0-9]))\b`)},
	{"size", false, reflect.String, regexp.MustCompile(`(?i)\b((\d+(?:\.\d+)?(?:GB|MB)))\b`)},
	{"website", false, reflect.String, regexp.MustCompile(`^(\[ ?([^\]]+?) ?\])`)},
	{"language", false, reflect.String, regexp.MustCompile(`(?i)\b((rus\.eng|ita\.eng))\b`)},
	{"sbs", false, reflect.String, regexp.MustCompile(`(?i)\b(((?:Half-)?SBS))\b`)},
	{"container", false, reflect.String, regexp.MustCompile(`(?i)\b((MKV|AVI|MP4))\b`)},

	{"group", false, reflect.String, regexp.MustCompile(`\b(- ?([^-]+(?:-={[^-]+-?$)?))$`)},

	{"extended", false, reflect.Bool, regexp.MustCompile(`(?i)\b(EXTENDED(:?.CUT)?)\b`)},
	{"hardcoded", false, reflect.Bool, regexp.MustCompile(`(?i)\b((HC))\b`)},
	{"proper", false, reflect.Bool, regexp.MustCompile(`(?i)\b((PROPER))\b`)},
	{"repack", false, reflect.Bool, regexp.MustCompile(`(?i)\b((REPACK))\b`)},
	{"widescreen", false, reflect.Bool, regexp.MustCompile(`(?i)\b((WS))\b`)},
	{"unrated", false, reflect.Bool, regexp.MustCompile(`(?i)\b((UNRATED))\b`)},
	{"threeD", false, reflect.Bool, regexp.MustCompile(`(?i)\b((3D))\b`)},
}

func init() {
	for _, pat := range patterns {
		if pat.re.NumSubexp() != 2 {
			fmt.Printf("Pattern %q does not have enough capture groups. want 2, got %d\n", pat.name, pat.re.NumSubexp())
			os.Exit(1)
		}
	}
}

// TorrentInfo is the resulting structure returned by Parse
type TorrentInfo struct {
	Title      string
	Season     int    `json:"season,omitempty"`
	Episode    int    `json:"episode,omitempty"`
	Year       int    `json:"year,omitempty"`
	Resolution string `json:"resolution,omitempty"`
	Quality    string `json:"quality,omitempty"`
	Codec      string `json:"codec,omitempty"`
	Audio      string `json:"audio,omitempty"`
	Group      string `json:"group,omitempty"`
	Region     string `json:"region,omitempty"`
	Extended   bool   `json:"extended,omitempty"`
	Hardcoded  bool   `json:"hardcoded,omitempty"`
	Proper     bool   `json:"proper,omitempty"`
	Repack     bool   `json:"repack,omitempty"`
	Container  string `json:"container,omitempty"`
	Widescreen bool   `json:"widescreen,omitempty"`
	Website    string `json:"website,omitempty"`
	Language   string `json:"language,omitempty"`
	Sbs        string `json:"sbs,omitempty"`
	Unrated    bool   `json:"unrated,omitempty"`
	Size       string `json:"size,omitempty"`
	ThreeD     bool   `json:"3d,omitempty"`
}

func setField(tor *TorrentInfo, field, raw, val string) {
	ttor := reflect.TypeOf(tor)
	torV := reflect.ValueOf(tor)
	field = strings.Title(field)
	v, _ := ttor.Elem().FieldByName(field)
	//fmt.Printf("    field=%v, type=%+v, value=%v\n", field, v.Type, val)
	switch v.Type.Kind() {
	case reflect.Bool:
		torV.Elem().FieldByName(field).SetBool(true)
	case reflect.Int:
		clean, _ := strconv.ParseInt(val, 10, 64)
		torV.Elem().FieldByName(field).SetInt(clean)
	case reflect.Uint:
		clean, _ := strconv.ParseUint(val, 10, 64)
		torV.Elem().FieldByName(field).SetUint(clean)
	case reflect.String:
		torV.Elem().FieldByName(field).SetString(val)
	}
}

// Parse breaks up the given filename in TorrentInfo
func Parse(filename string) (*TorrentInfo, error) {
	tor := &TorrentInfo{}
	//fmt.Printf("filename %q\n", filename)

	var startIndex, endIndex = 0, len(filename)
	cleanName := strings.Replace(filename, "_", " ", -1)
	for _, pattern := range patterns {
		matches := pattern.re.FindAllStringSubmatch(cleanName, -1)
		if len(matches) == 0 {
			continue
		}
		matchIdx := 0
		if pattern.last {
			// Take last occurence of element.
			matchIdx = len(matches) - 1
		}
		//fmt.Printf("  %s: pattern:%q match:%#v\n", pattern.name, pattern.re, matches[matchIdx])

		index := strings.Index(cleanName, matches[matchIdx][1])
		if index == 0 {
			startIndex = len(matches[matchIdx][1])
			//fmt.Printf("    startIndex moved to %d [%q]\n", startIndex, filename[startIndex:endIndex])
		} else if index < endIndex {
			endIndex = index
			//fmt.Printf("    endIndex moved to %d [%q]\n", endIndex, filename[startIndex:endIndex])
		}
		setField(tor, pattern.name, matches[matchIdx][1], matches[matchIdx][2])
	}

	// Start process for title
	//fmt.Println("  title: <internal>")
	raw := strings.Split(filename[startIndex:endIndex], "(")[0]
	cleanName = raw
	if strings.HasPrefix(cleanName, "- ") {
		cleanName = raw[2:]
	}
	if strings.ContainsRune(cleanName, '.') && !strings.ContainsRune(cleanName, ' ') {
		cleanName = strings.Replace(cleanName, ".", " ", -1)
	}
	cleanName = strings.Replace(cleanName, "_", " ", -1)
	//cleanName = re.sub('([\[\(_]|- )$', '', cleanName).strip()
	setField(tor, "title", raw, strings.TrimSpace(cleanName))

	return tor, nil
}
