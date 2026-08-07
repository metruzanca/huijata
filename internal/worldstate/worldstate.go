package worldstate

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
)

// State is the WorldStateComponent of a run's world_state.xml.
type State struct {
	// Attrs are the WorldStateComponent's attributes.
	Attrs map[string]string
}

// Load reads and parses world_state.xml from the save slot folder.
func Load(savePath string) (*State, error) {
	data, err := os.ReadFile(filepath.Join(savePath, "world_state.xml"))
	if err != nil {
		return nil, err
	}
	var root element
	if err := xml.Unmarshal(data, &root); err != nil {
		return nil, err
	}
	var comp *element
	walk(&root, func(e *element) {
		if comp == nil && e.XMLName.Local == "WorldStateComponent" {
			comp = e
		}
	})
	if comp == nil {
		return nil, fmt.Errorf("world_state.xml has no WorldStateComponent")
	}
	return &State{Attrs: attrs(comp)}, nil
}

type element struct {
	XMLName xml.Name
	Attrs   []xml.Attr `xml:",any,attr"`
	Kids    []element  `xml:",any"`
}

func walk(e *element, fn func(*element)) {
	fn(e)
	for i := range e.Kids {
		walk(&e.Kids[i], fn)
	}
}

func attrs(e *element) map[string]string {
	m := make(map[string]string, len(e.Attrs))
	for _, a := range e.Attrs {
		m[a.Name.Local] = a.Value
	}
	return m
}
