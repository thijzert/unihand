package unihand

import (
	"fmt"
	"log"
)

func init() {
	// Convert Pinyin from impossible-to-type to easy-to-type format
	fileFilters = append(fileFilters, fileFilter{
		Filename:     "Unihan_Readings.txt",
		RecordFilter: pinyin,
	})

	// English translation
	fileFilters = append(fileFilters, fileFilter{
		Filename:     "Unihan_Readings.txt",
		RecordFilter: englishDefinition,
	})
}

func englishDefinition(char *Character, code uint32, fields []string) error {
	if len(fields) < 2 {
		return fmt.Errorf("got invalid record [%s]", fields)
	}

	if fields[0] == "kDefinition" {
		char.English = fields[1]
	}

	return nil
}

func pinyin(char *Character, code uint32, fields []string) error {
	if len(fields) < 2 {
		return fmt.Errorf("got invalid record [%s]", fields)
	}

	if fields[0] == "kMandarin" {
		var hans, hant string
		n, _ := fmt.Sscan(fields[1], &hans, &hant)
		if n >= 1 {
			var t error
			char.Pinyin, t = toneNumber(hans)
			if t != nil {
				log.Printf("odd tone for U+%04X (%c): %s", code, rune(code), t.Error())
				return nil
			}
			if n == 1 {
				// char.PinyinTW = hans
			} else {
				char.PinyinTW, _ = toneNumber(hant)
			}
		}
	}
	return nil
}

var toneChars = []struct {
	Char rune
	Tone int
	Word string
}{
	// Combining characters
	{0x0304, 1, ""},
	{0x0301, 2, ""},
	{0x030c, 3, ""},
	{0x0300, 4, ""},

	// Non-combining characters
	{0x0101, 1, "a"},
	{0x00E1, 2, "a"},
	{0x01CE, 3, "a"},
	{0x00E0, 4, "a"},

	{0x0113, 1, "e"},
	{0x00E9, 2, "e"},
	{0x011B, 3, "e"},
	{0x00E8, 4, "e"},

	{0x012B, 1, "i"},
	{0x00ED, 2, "i"},
	{0x01D0, 3, "i"},
	{0x00EC, 4, "i"},

	{0x014D, 1, "o"},
	{0x00F3, 2, "o"},
	{0x01D2, 3, "o"},
	{0x00F2, 4, "o"},

	{0x016B, 1, "u"},
	{0x00FA, 2, "u"},
	{0x01D4, 3, "u"},
	{0x00F9, 4, "u"},

	// U-umlaut -> v
	{0x00FC, 0, "v"},
	{0x01D6, 1, "v"},
	{0x01D8, 2, "v"},
	{0x01DA, 3, "v"},
	{0x01DC, 4, "v"},

	// Whatever, man
	{0x0144, 2, "n"},
	{0x0148, 3, "n"},
	{0x01f9, 4, "n"},
	{0x1e3f, 2, "m"},
}

func toneNumber(pinyin string) (string, error) {
	tone := 5
	word := ""
	for _, p := range pinyin {
		found := false
		for _, py := range toneChars {
			if py.Char == p {
				found = true
				if py.Tone != 0 {
					tone = py.Tone
				}
				word += py.Word
				break
			}
		}

		if found {
			// Yay!
		} else if p == 0x0308 {
			// Combining Umlaut - change 'u' to 'v'
			if len(word) > 0 && word[len(word)-1] == 'u' {
				word = word[:len(word)-1] + "v"
			} else {
				log.Printf("umlaut? %02x", []byte(pinyin))
			}
		} else if p >= 'a' && p <= 'z' {
			word += string(p)
		} else {
			return "", fmt.Errorf("invalid character U+%04x in pinyin '%s'", p, pinyin)
		}
	}
	return fmt.Sprintf("%s%d", word, tone), nil
}
