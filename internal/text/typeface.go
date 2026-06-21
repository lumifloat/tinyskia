package text

import (
	"sync"
	"sync/atomic"

	"golang.org/x/exp/mmap"
	"golang.org/x/image/font/sfnt"
)

type Typeface struct {
	sync.Once
	*sfnt.Font

	file  string
	index int

	family string
	weight float64
	style  string

	rejects [1024]uint64
}

func (tf *Typeface) GlyphIndex(b *sfnt.Buffer, r rune) (sfnt.GlyphIndex, error) {
	tf.Do(func() {
		reader, err := mmap.Open(tf.file)
		if err != nil {
			return
		}

		collection, err := sfnt.ParseCollectionReaderAt(reader)
		if err != nil {
			_ = reader.Close()
			return
		}

		ft, err := collection.Font(tf.index)
		if err != nil {
			_ = reader.Close()
			return
		}

		tf.Font = ft
	})

	if tf.Font == nil {
		return 0, sfnt.ErrNotFound
	}

	if r > 0xFFFF {
		return tf.Font.GlyphIndex(b, r)
	}

	// Fast reject
	idx, pos := r>>6, r&0x3F
	mask := atomic.LoadUint64(&tf.rejects[idx])
	if (mask & (uint64(1) << pos)) != 0 {
		return 0, sfnt.ErrNotFound
	}

	gid, err := tf.Font.GlyphIndex(b, r)
	if err == nil && gid == 0 {
		bitMask := uint64(1) << pos
		oldMask := atomic.LoadUint64(&tf.rejects[idx])
		newMask := oldMask | bitMask
		atomic.CompareAndSwapUint64(&tf.rejects[idx], oldMask, newMask)
		return 0, sfnt.ErrNotFound
	}
	return gid, err
}
