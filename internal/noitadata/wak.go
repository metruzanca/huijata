package noitadata

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

// wakEntry is one file described in a Noita .wak archive directory.
type wakEntry struct {
	name string
	// end is the absolute file offset where this entry's data ends,
	// which is also where the next entry's data begins.
	end uint32
}

// Wak is a Noita data.wak archive.
type Wak struct {
	path       string
	dataOffset uint32
	entries    []wakEntry
}

// OpenWak reads a data.wak header and directory. The data section is only
// read when an entry is requested.
func OpenWak(path string) (*Wak, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var hdr [24]byte
	if _, err := io.ReadFull(f, hdr[:]); err != nil {
		return nil, fmt.Errorf("read %s header: %w", path, err)
	}
	dataOffset := binary.LittleEndian.Uint32(hdr[8:12])

	dir := make([]byte, int(dataOffset)-len(hdr))
	if _, err := io.ReadFull(f, dir); err != nil {
		return nil, fmt.Errorf("read %s directory: %w", path, err)
	}

	w := &Wak{path: path, dataOffset: dataOffset}
	pos := 0
	for pos < len(dir) {
		if pos+8 > len(dir) {
			break
		}
		nameLen := int(binary.LittleEndian.Uint32(dir[pos:]))
		if nameLen == 0 || nameLen > 512 || pos+4+nameLen+8 > len(dir) {
			break
		}
		name := string(dir[pos+4 : pos+4+nameLen])
		end := binary.LittleEndian.Uint32(dir[pos+4+nameLen:])
		w.entries = append(w.entries, wakEntry{name: name, end: end})
		pos += 4 + nameLen + 8
	}
	return w, nil
}

// Find returns the contents of the entry with the given name (decompressed
// if it was stored compressed).
func (w *Wak) Find(name string) ([]byte, error) {
	for i := range w.entries {
		if w.entries[i].name != name {
			continue
		}
		start := w.dataOffset
		if i > 0 {
			start = w.entries[i-1].end
		}
		end := w.entries[i].end
		if end < start || end-start > 200<<20 {
			return nil, fmt.Errorf("%s: bad wak bounds %d..%d", name, start, end)
		}
		blob := make([]byte, end-start)
		f, err := os.Open(w.path)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		if _, err := f.ReadAt(blob, int64(start)); err != nil {
			return nil, err
		}
		return maybeInflate(blob), nil
	}
	return nil, fmt.Errorf("data.wak has no entry %q", name)
}

// maybeInflate returns the zlib-decompressed contents of data when it looks
// like a zlib stream, and data itself otherwise.
func maybeInflate(data []byte) []byte {
	if len(data) < 2 || data[0] != 0x78 {
		return data
	}
	r, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		return data
	}
	out, err := io.ReadAll(r)
	if err != nil {
		return data
	}
	return out
}
