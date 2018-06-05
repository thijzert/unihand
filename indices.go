package unihand

import "fmt"

func init() {
	fileFilters = append(fileFilters, fileFilter{
		Filename:     "Unihan_IRGSources.txt",
		RecordFilter: irgSource,
	})
}

func irgSource(char *Character, code uint32, sources []string) error {
	if len(sources) < 2 {
		return fmt.Errorf("got invalid record [%s]", sources)
	}

	if sources[0] == "kIRG_GSource" {
		char.IRGSources.G = sources[1]
	} else if sources[0] == "kIRG_HSource" {
		char.IRGSources.H = sources[1]
	} else if sources[0] == "kIRG_JSource" {
		char.IRGSources.J = sources[1]
	} else if sources[0] == "kIRG_KPSource" {
		char.IRGSources.KP = sources[1]
	} else if sources[0] == "kIRG_KSource" {
		char.IRGSources.K = sources[1]
	} else if sources[0] == "kIRG_MSource" {
		char.IRGSources.M = sources[1]
	} else if sources[0] == "kIRG_TSource" {
		char.IRGSources.T = sources[1]
	} else if sources[0] == "kIRG_VSource" {
		char.IRGSources.V = sources[1]
	}

	return nil
}
