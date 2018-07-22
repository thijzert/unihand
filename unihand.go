package unihand

import (
	"archive/zip"
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"strings"
)

// The Character struct represents a single Unicode character.
type Character struct {
	// Whether this character is actually loaded in the database
	loaded bool

	// The Unicode code point
	Unicode uint32

	// The number of strokes in this character, for mainland and Taiwan
	Strokes   uint8
	StrokesTW uint8 `json:",omitempty"`

	// Pronunciation (pinyin)
	Pinyin string `json:",omitempty"`

	// Pronunciation for Taiwanese, if different
	PinyinTW string `json:",omitempty"`

	// English translation, where applicable
	English string `json:",omitempty"`

	// Source information from the Ideographic Rapporteur Group
	IRGSources struct {
		// PRC and Singapore
		G string `json:",omitempty"`
		// Hong Kong
		H string `json:",omitempty"`
		// Japan
		J string `json:",omitempty"`
		// Best Korea
		KP string `json:",omitempty"`
		// South Korea
		K string `json:",omitempty"`
		// Macao
		M string `json:",omitempty"`
		// Taiwan
		T string `json:",omitempty"`
		// Vietnam
		V string `json:",omitempty"`
	}

	// KangXi Radical number
	KangXiRadicalNumber uint8 `json:",omitempty"`

	// Radical-Stroke Counts
	RadicalStrokeCounts struct {
		KangXi  RScount
		Unicode RScount
	}
}

type RScount struct {
	Radical           uint8
	SimplifiedVersion bool
	Additional        int8
}

func (rs RScount) MarshalJSON() ([]byte, error) {
	if rs.SimplifiedVersion {
		return json.Marshal(fmt.Sprintf("%d'.%d", rs.Radical, rs.Additional))
	} else {
		return json.Marshal(fmt.Sprintf("%d'.%d", rs.Radical, rs.Additional))
	}
}

// The database contains the completish Unihan database in memory
var database = struct {
	Contiguous []Character
	Misc       map[uint32]Character
}{
	Contiguous: nil,
	Misc:       nil,
}

// A special marker that marks the end of the contiguous index
var endOfIndex uint32

type codeRange struct {
	From uint32
	To   uint32
}

// Code2index converts an unicode codepoint to an index in the contiguous range, or zero if it doesn't exist there
func code2index(code uint32) int {
	ranges := []codeRange{
		codeRange{0x4e00, 0x9fff},
		codeRange{0x2e80, 0x2eff},
		codeRange{0x3000, 0x303f},
		codeRange{0x31c0, 0x31ff},
		codeRange{0x3200, 0x32ff},
		codeRange{0x3300, 0x33ff},
		codeRange{0x3400, 0x4dff},
		codeRange{0xf900, 0xfaff},
		codeRange{0xfe30, 0xfe4f},
		codeRange{0x20000, 0x2a6df},
		codeRange{0x2a700, 0x2b73f},
		codeRange{0x2b740, 0x2b81f},
		codeRange{0x2b820, 0x2ceaf},
		codeRange{0x2ceb0, 0x2ebef},
		codeRange{0x2f800, 0x2fa1f},
		codeRange{endOfIndex, endOfIndex},
	}

	offset := 0
	for _, r := range ranges {
		if code >= r.From && code <= r.To {
			return offset + int(code) - int(r.From)
		}
		offset += 1 + int(r.To) - int(r.From)
	}
	return 0
}

func init() {
	// Assign a random code point from a Private Use Area
	endOfIndex = 0xffff
	for (endOfIndex & 0xfffe) == 0xfffe {
		endOfIndex = rand.Uint32() & 0xffff
	}
	endOfIndex = 0x100000 | endOfIndex

	database.Contiguous = make([]Character, code2index(endOfIndex))
	database.Misc = make(map[uint32]Character)
}

type fileFilter struct {
	Filename     string
	RecordFilter func(*Character, uint32, []string) error
}

// A list of file filters
var fileFilters = []fileFilter{}

// Initialise reads the Unihan.zip file into the in-memory database.
func Initialise(path string) error {
	r, err := zip.OpenReader(path)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, ff := range fileFilters {
		for _, file := range r.File {
			if file.Name == ff.Filename {
				r, err := file.Open()
				if err != nil {
					return err
				}
				br := bufio.NewReader(r)
				nline := 0
				for {
					nline++
					var line string
					line, err = br.ReadString('\n')
					if err != nil {
						break
					}
					if len(line) < 7 || line[0] == '#' {
						continue
					}

					var code uint32
					_, err = fmt.Sscanf(line, "U+%X\t", &code)
					if err != nil {
						err = fmt.Errorf("%s:%d: %s", ff.Filename, nline, err.Error())
						break
					}
					for n, c := range line {
						if c == '\t' {
							line = line[n+1:]
							break
						}
					}
					for len(line) > 0 && line[len(line)-1] == '\n' {
						line = line[:len(line)-1]
					}

					if code == endOfIndex {
						continue
					}

					idx := code2index(code)
					if idx > 0 {
						database.Contiguous[idx].loaded = true
						database.Contiguous[idx].Unicode = code
						err = ff.RecordFilter(&database.Contiguous[idx], code, strings.Split(line, "\t"))
					} else {
						c, _ := database.Misc[code]
						c.Unicode = code
						c.loaded = true
						err = ff.RecordFilter(&c, code, strings.Split(line, "\t"))
						database.Misc[code] = c
					}
					if err != nil {
						err = fmt.Errorf("%s:%d: %s", ff.Filename, nline, err.Error())
						break
					}
				}
				if err != io.EOF {
					return err
				}

				break
			}
		}
	}
	return nil
}

// LoadedCharacters returns the number of characters loaded into the database
func LoadedCharacters() int {
	rv := len(database.Misc)
	for _, c := range database.Contiguous {
		if c.loaded {
			rv++
		}
	}
	return rv
}

func Lookup(char uint32) (rv Character, err error) {
	if char == endOfIndex {
		err = fmt.Errorf("no way, Jose")
		return
	}

	idx := code2index(char)
	if idx > 0 {
		rv = database.Contiguous[idx]
		if !rv.loaded {
			err = fmt.Errorf("not found")
		}
		return
	} else {
		var ok bool
		rv, ok = database.Misc[char]
		if !ok {
			err = fmt.Errorf("not found")
		}
		return
	}
}
