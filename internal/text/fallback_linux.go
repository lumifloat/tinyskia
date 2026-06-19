// Copyright 2012 The Chromium Authors
// Copyright 2026 Mozilla
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package text

func GetCommonFallbackFonts(list []string, ch rune, script Script) {
	list = append(list,
		"DejaVu Serif",
		"FreeSerif",
		"DejaVu Sans",
		"FreeSans",
		"Symbola",
		"Noto Sans Symbols",
		"Noto Sans Symbols2",
	)

	if ch >= 0x3000 && ((ch < 0xe000) ||
		(ch >= 0xf900 && ch < 0xfff0) ||
		((uint32(ch) >> 16) == 2)) {
		list = append(list,
			"TakaoPGothic",
			"Droid Sans Fallback",
			"WenQuanYi Micro Hei",
			"NanumGothic",
		)
	}

	list = append(list, "Twemoji Mozilla")
}
