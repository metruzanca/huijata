package stats

import (
	"encoding/xml"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/metruzanca/huijata/internal/salakieli"
)

// Stats is the decrypted content of a save's stats/_stats.salakieli file.
type Stats struct {
	// SessionDead reports whether the current run is over.
	SessionDead bool
	// Version is the file's STATS_VERSION.
	Version int
	// Entries are the tracked key/value stats, e.g. enemy names -> kills
	// and action ids -> casts.
	Entries []Entry
	// Session, Highest, Global and PrevBest hold the per-run aggregate
	// blocks written by the game.
	Session  Section
	Highest  Section
	Global   Section
	PrevBest Section
}

// Entry is one key/value pair of a KEY_VALUE_STATS block.
type Entry struct {
	Key   string
	Value string
}

// Section holds the attributes of one aggregate stats block (session,
// highest, global, prev_best). The game writes each stat as an XML attribute.
type Section map[string]string

// Load reads and decrypts stats/_stats.salakieli from the save slot folder.
func Load(savePath string) (*Stats, error) {
	data, err := salakieli.DecryptFile(filepath.Join(savePath, "stats", "_stats.salakieli"))
	if err != nil {
		return nil, err
	}
	return Parse(data)
}

// Parse parses the decrypted GameStats XML.
func Parse(data []byte) (*Stats, error) {
	var root element
	if err := xml.Unmarshal(data, &root); err != nil {
		return nil, err
	}
	if root.XMLName.Local != "GameStats" {
		return nil, fmt.Errorf("unexpected stats root element %q", root.XMLName.Local)
	}

	s := &Stats{}
	a := attrs(&root)
	s.Version = intVal(a["STATS_VERSION"])
	s.SessionDead = a["session_dead"] == "1"

	for i := range root.Kids {
		k := &root.Kids[i]
		switch k.XMLName.Local {
		case "KEY_VALUE_STATS":
			for j := range k.Kids {
				e := attrs(&k.Kids[j])
				if e["key"] != "" {
					s.Entries = append(s.Entries, Entry{Key: e["key"], Value: e["value"]})
				}
			}
		case "session":
			s.Session = section(k)
		case "highest":
			s.Highest = section(k)
		case "global":
			s.Global = section(k)
		case "prev_best":
			s.PrevBest = section(k)
		}
	}
	return s, nil
}

type element struct {
	XMLName xml.Name
	Attrs   []xml.Attr `xml:",any,attr"`
	Kids    []element  `xml:",any"`
}

func section(e *element) Section {
	return Section(attrs(e))
}

func attrs(e *element) map[string]string {
	m := make(map[string]string, len(e.Attrs))
	for _, a := range e.Attrs {
		m[a.Name.Local] = a.Value
	}
	return m
}

func intVal(s string) int {
	n, _ := strconv.Atoi(strings.TrimSpace(s))
	return n
}
