package unihand

import (
	"fmt"
)

func init() {
	// The field 'kRSUnicode' is included with the IRGSources document for some reason
	fileFilters = append(fileFilters, fileFilter{
		Filename:     "Unihan_RadicalStrokeCounts.txt",
		RecordFilter: radicalStrokeCount,
	}, fileFilter{
		Filename:     "Unihan_IRGSources.txt",
		RecordFilter: radicalStrokeCount,
	})

	fileFilters = append(fileFilters, fileFilter{
		Filename:     "Unihan_DictionaryLikeData.txt",
		RecordFilter: totalStrokes,
	})
}

func radicalStrokeCount(char *Character, code uint32, fields []string) error {
	if len(fields) < 2 {
		return fmt.Errorf("got invalid record [%s]", fields)
	}

	if fields[0] == "kRSKangXi" || fields[0] == "kRSUnicode" {
		var rsc RScount
		_, err := fmt.Sscanf(fields[1], "%d.%d", &rsc.Radical, &rsc.Additional)
		if err != nil {
			_, err := fmt.Sscanf(fields[1], "%d'.%d", &rsc.Radical, &rsc.Additional)
			if err != nil {
				return err
			} else {
				rsc.SimplifiedVersion = true
			}
		}
		if fields[0] == "kRSKangXi" {
			char.RadicalStrokeCounts.KangXi = rsc
		} else if fields[0] == "kRSUnicode" {
			char.RadicalStrokeCounts.Unicode = rsc
		}
	}

	return nil
}

func totalStrokes(char *Character, code uint32, fields []string) error {
	if len(fields) < 2 {
		return fmt.Errorf("got invalid record [%s]", fields)
	}

	if fields[0] == "kTotalStrokes" {
		var hans, hant uint8
		n, _ := fmt.Sscan(fields[1], &hans, &hant)
		if n >= 1 {
			char.Strokes = hans
			if n == 1 {
				char.StrokesTW = hans
			} else {
				char.StrokesTW = hant
			}
		}
	}
	return nil
}
