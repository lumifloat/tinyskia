// Copyright 2012 The Chromium Authors
// Copyright 2026 Mozilla
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package text

func GetFallbackFonts(script Script) []string {
	var list = make([]string, 0, 4)
	switch script {
	case ScriptBopo, ScriptHanb, ScriptHans, ScriptHani:
		list = append(list, "WenQuanYi Micro Hei", "Droid Sans Fallback")
	case ScriptHant, ScriptHntl:
		list = append(list, "WenQuanYi Micro Hei", "Droid Sans Fallback")
	case ScriptHira, ScriptKana, ScriptHrkt, ScriptJpan:
		list = append(list, "TakaoPGothic", "WenQuanYi Micro Hei", "Droid Sans Fallback")
	case ScriptHang, ScriptJamo, ScriptKore:
		list = append(list, "NanumGothic", "WenQuanYi Micro Hei", "Droid Sans Fallback")
	}

	list = append(list,
		"DejaVu Serif",
		"FreeSerif",
		"DejaVu Sans",
		"FreeSans",
		"Symbola",
		"Noto Sans Symbols",
		"Noto Sans Symbols2",
	)

	list = append(list, "Twemoji Mozilla")

	return list
}
