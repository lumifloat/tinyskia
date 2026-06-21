package text

import (
	"math"
	"os"
	"runtime"
	"slices"
	"strings"
	"sync"

	"golang.org/x/image/font/sfnt"
	"golang.org/x/text/language"
)

type Registry struct {
	sync.Mutex

	assets map[string][]*Typeface
	locals map[string][]*Typeface
}

func NewRegistry() *Registry {
	return &Registry{
		assets: make(map[string][]*Typeface),
		locals: make(map[string][]*Typeface),
	}
}

func (reg *Registry) RegisterFont(font *sfnt.Font, family string, weight float64, style string) error {
	reg.Lock()
	defer reg.Unlock()

	typefaces := &Typeface{
		Font:   font,
		family: family,
		weight: weight,
		style:  style,
	}

	typefaces.Do(func() {})

	reg.assets[family] = append(reg.assets[family], typefaces)
	return nil
}

func (reg *Registry) ScanSystemFonts() error {
	var base string
	switch os := runtime.GOOS; os {
	case "windows":
		base = "C:/Windows/Fonts/"
	case "linux":
		base = "/usr/share/fonts/"
	case "darwin":
		base = "/Library/Fonts/"
	}
	entries, err := os.ReadDir(base)
	if err != nil {
		return err
	}
	reg.Lock()
	defer reg.Unlock()
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if !strings.HasSuffix(entry.Name(), ".ttf") && !strings.HasSuffix(entry.Name(), ".otf") &&
			!strings.HasSuffix(entry.Name(), ".ttc") {
			continue
		}
		fi, err := os.Open(base + entry.Name())
		if err != nil {
			continue
		}
		collection, err := sfnt.ParseCollectionReaderAt(fi)
		if err != nil {
			continue
		}
		for i := 0; i < collection.NumFonts(); i++ {
			ttf, err := collection.Font(i)
			if err != nil {
				continue
			}

			family, err := ttf.Name(nil, sfnt.NameIDTypographicFamily)
			if err != nil {
				family, err = ttf.Name(nil, sfnt.NameIDFamily)
				if err != nil {
					continue
				}
			}
			family = strings.ToLower(family)

			typefaces := &Typeface{
				Font:   nil,
				file:   base + entry.Name(),
				family: family,
				weight: 0,
				style:  "normal",
			}

			reg.locals[family] = append(reg.locals[family], typefaces)
		}
		fi.Close()
	}
	return nil
}

func (reg *Registry) GetFonts(family []string, weight float64, style string) []*Typeface {
	best := func(a, b *Typeface) int {
		if (a.style == style) != (b.style == style) {
			if a.style == style {
				return -1
			}
			return 1
		}
		if (a.weight == weight) != (b.weight == weight) {
			if a.weight == weight {
				return -1
			}
			return 1
		}
		return int(math.Abs(a.weight-weight) - math.Abs(b.weight-weight))
	}

	reg.Lock()
	defer reg.Unlock()

	// TODO: should more faster than this
	var chain = make([]*Typeface, 0, len(FONT_CURSIVE_FONT_FAMILY_TAML))
	for i := range family {
		family := strings.ToLower(family[i])
		if tf, ok := reg.assets[family]; ok {
			if len(tf) != 0 {
				slices.SortFunc(tf, best)
				chain = append(chain, tf[0])
			}
		}
		if tf, ok := reg.locals[family]; ok {
			if len(tf) != 0 {
				slices.SortFunc(tf, best)
				chain = append(chain, tf[0])
			}
		}
	}
	return chain
}

func (reg *Registry) GlyphIndex(b *sfnt.Buffer, r rune, script Script, lang language.Script, family []string, weight float64, style string) (*Typeface, sfnt.GlyphIndex, error) {
	locale := Scripts[lang.String()]
	var override = make([]string, 0, len(family))
	for i := range family {
		family := strings.ToLower(family[i])
		if IsGenericFamily(family) {
			override = append(override, GetPreferenceFonts(family, locale)...)
			override = append(override, GetPreferenceFonts(family, script)...)
			override = append(override, GetPreferenceFonts(family, ScriptLatn)...)
		} else {
			override = append(override, family)
		}
	}
	override = append(override, GetFallbackFonts(script)...)
	fonts := reg.GetFonts(override, weight, style)
	for i := range fonts {
		gid, err := fonts[i].GlyphIndex(b, r)
		if err != nil || gid == 0 {
			continue
		}
		return fonts[i], gid, nil
	}
	return nil, 0, nil
}
