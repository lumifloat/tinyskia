// Copyright 2012 The Chromium Authors
// Copyright 2026 Mozilla
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package text

var (
	FONT_FANTASY_FONT_FAMILY = []string{"Impact"}
	FONT_EMOJI_FONT_FAMILY   = []string{"Segoe UI Emoji", "Twemoji Mozilla"}
	FONT_MATH_FONT_FAMILY    = []string{"Cambria Math"}

	FONT_SERIF_FONT_FAMILY_ARAB      = []string{"Times New Roman"}
	FONT_SANS_SERIF_FONT_FAMILY_ARAB = []string{"Segoe UI", "Tahoma", "Arial"}
	FONT_MONOSPACE_FONT_FAMILY_ARAB  = []string{"Consolas"}
	FONT_CURSIVE_FONT_FAMILY_ARAB    = []string{"Comic Sans MS"}

	FONT_SERIF_FONT_FAMILY_ARMN      = []string{"Sylfaen"}
	FONT_SANS_SERIF_FONT_FAMILY_ARMN = []string{"Arial AMU"}
	FONT_MONOSPACE_FONT_FAMILY_ARMN  = []string{"Arial AMU"}
	FONT_CURSIVE_FONT_FAMILY_ARMN    = []string{}

	FONT_SERIF_FONT_FAMILY_BENG      = []string{"Vrinda", "Akaash", "Likhan", "Ekushey Punarbhaba"}
	FONT_SANS_SERIF_FONT_FAMILY_BENG = []string{"Nirmala Text", "Vrinda", "Akaash", "Likhan", "Ekushey Punarbhaba"}
	FONT_MONOSPACE_FONT_FAMILY_BENG  = []string{"Mitra Mono", "Likhan", "Mukti Narrow"}
	FONT_CURSIVE_FONT_FAMILY_BENG    = []string{}

	FONT_SERIF_FONT_FAMILY_CANS      = []string{"Aboriginal Serif", "BJCree Uni"}
	FONT_SANS_SERIF_FONT_FAMILY_CANS = []string{"Aboriginal Sans"}
	FONT_MONOSPACE_FONT_FAMILY_CANS  = []string{"Aboriginal Sans", "OskiDakelh", "Pigiarniq", "Uqammaq"}
	FONT_CURSIVE_FONT_FAMILY_CANS    = []string{}

	FONT_SERIF_FONT_FAMILY_CYRL      = []string{"Times New Roman"}
	FONT_SANS_SERIF_FONT_FAMILY_CYRL = []string{"Arial"}
	FONT_MONOSPACE_FONT_FAMILY_CYRL  = []string{"Consolas"}
	FONT_CURSIVE_FONT_FAMILY_CYRL    = []string{"Comic Sans MS"}

	FONT_SERIF_FONT_FAMILY_DEVA      = []string{"Kokila", "Raghindi"}
	FONT_SANS_SERIF_FONT_FAMILY_DEVA = []string{"Nirmala UI", "Mangal"}
	FONT_MONOSPACE_FONT_FAMILY_DEVA  = []string{"Mangal", "Nirmala UI"}
	FONT_CURSIVE_FONT_FAMILY_DEVA    = []string{}

	FONT_SERIF_FONT_FAMILY_ETHI      = []string{"Visual Geez Unicode", "Visual Geez Unicode Agazian"}
	FONT_SANS_SERIF_FONT_FAMILY_ETHI = []string{"GF Zemen Unicode"}
	FONT_MONOSPACE_FONT_FAMILY_ETHI  = []string{"Ethiopia Jiret"}
	FONT_CURSIVE_FONT_FAMILY_ETHI    = []string{"Visual Geez Unicode Title"}

	FONT_SERIF_FONT_FAMILY_GEOR      = []string{"Sylfaen", "BPG Paata Khutsuri U", "TITUS Cyberbit Basic"}
	FONT_SANS_SERIF_FONT_FAMILY_GEOR = []string{"BPG Classic 99U"}
	FONT_MONOSPACE_FONT_FAMILY_GEOR  = []string{"BPG Classic 99U"}
	FONT_CURSIVE_FONT_FAMILY_GEOR    = []string{}

	FONT_SERIF_FONT_FAMILY_GREK      = []string{"Times New Roman"}
	FONT_SANS_SERIF_FONT_FAMILY_GREK = []string{"Arial"}
	FONT_MONOSPACE_FONT_FAMILY_GREK  = []string{"Consolas"}
	FONT_CURSIVE_FONT_FAMILY_GREK    = []string{"Comic Sans MS"}

	FONT_SERIF_FONT_FAMILY_GUJR      = []string{"Shruti"}
	FONT_SANS_SERIF_FONT_FAMILY_GUJR = []string{"Shruti"}
	FONT_MONOSPACE_FONT_FAMILY_GUJR  = []string{"Shruti"}
	FONT_CURSIVE_FONT_FAMILY_GUJR    = []string{}

	FONT_SERIF_FONT_FAMILY_GURU      = []string{"Raavi", "Saab"}
	FONT_SANS_SERIF_FONT_FAMILY_GURU = []string{}
	FONT_MONOSPACE_FONT_FAMILY_GURU  = []string{"Raavi", "Saab"}
	FONT_CURSIVE_FONT_FAMILY_GURU    = []string{}

	FONT_SERIF_FONT_FAMILY_HANS      = []string{"SimSun", "MS Song", "SimSun-ExtB", "Noto Serif CJK SC", "Noto Serif SC"}
	FONT_SANS_SERIF_FONT_FAMILY_HANS = []string{"Microsoft YaHei", "SimHei", "Noto Sans CJK SC", "Noto Sans SC"}
	FONT_MONOSPACE_FONT_FAMILY_HANS  = []string{"NSimSun", "SimSun", "MS Song", "SimSun-ExtB"}
	FONT_CURSIVE_FONT_FAMILY_HANS    = []string{"KaiTi", "KaiTi_GB2312"}

	FONT_SERIF_FONT_FAMILY_HANT      = []string{"Times New Roman", "PMingLiu", "MingLiU", "MingLiU-ExtB", "Noto Serif CJK TC", "Noto Serif TC"}
	FONT_SANS_SERIF_FONT_FAMILY_HANT = []string{"Arial", "Microsoft JhengHei", "PMingLiu", "MingLiU", "MingLiU-ExtB", "Noto Sans CJK TC", "Noto Sans TC"}
	FONT_MONOSPACE_FONT_FAMILY_HANT  = []string{"MingLiU", "MingLiU-ExtB"}
	FONT_CURSIVE_FONT_FAMILY_HANT    = []string{"DFKai-SB"}

	FONT_SERIF_FONT_FAMILY_HEBR      = []string{"Narkisim", "David"}
	FONT_SANS_SERIF_FONT_FAMILY_HEBR = []string{"Arial"}
	FONT_MONOSPACE_FONT_FAMILY_HEBR  = []string{"Fixed Miriam Transparent", "Miriam Fixed", "Rod", "Consolas", "Courier New"}
	FONT_CURSIVE_FONT_FAMILY_HEBR    = []string{"Guttman Yad", "Ktav", "Arial"}

	FONT_SERIF_FONT_FAMILY_JPAN      = []string{"Yu Mincho", "MS PMincho", "MS Mincho", "Noto Serif CJK JP", "Noto Serif JP", "Meiryo", "Yu Gothic", "MS PGothic", "MS Gothic"}
	FONT_SANS_SERIF_FONT_FAMILY_JPAN = []string{"Meiryo", "Yu Gothic", "MS PGothic", "MS Gothic", "Noto Sans CJK JP", "Noto Sans JP"}
	FONT_MONOSPACE_FONT_FAMILY_JPAN  = []string{"BIZ UDGothic", "MS Gothic", "MS Mincho", "Meiryo", "Yu Gothic", "Yu Mincho"}
	FONT_CURSIVE_FONT_FAMILY_JPAN    = []string{}

	FONT_SERIF_FONT_FAMILY_KHMR      = []string{"PhnomPenh OT", ".Mondulkiri U GR 1.5", "Khmer OS"}
	FONT_SANS_SERIF_FONT_FAMILY_KHMR = []string{"Khmer OS"}
	FONT_MONOSPACE_FONT_FAMILY_KHMR  = []string{"Khmer OS", "Khmer OS System"}
	FONT_CURSIVE_FONT_FAMILY_KHMR    = []string{}

	FONT_SERIF_FONT_FAMILY_KNDA      = []string{"Tunga", "AksharUnicode"}
	FONT_SANS_SERIF_FONT_FAMILY_KNDA = []string{"Tunga", "AksharUnicode"}
	FONT_MONOSPACE_FONT_FAMILY_KNDA  = []string{"Tunga", "AksharUnicode"}
	FONT_CURSIVE_FONT_FAMILY_KNDA    = []string{}

	FONT_SERIF_FONT_FAMILY_KORE      = []string{"Batang", "Noto Serif CJK KR", "Noto Serif KR", "Gulim"}
	FONT_SANS_SERIF_FONT_FAMILY_KORE = []string{"Malgun Gothic", "Gulim", "Noto Sans CJK KR", "Noto Sans KR"}
	FONT_MONOSPACE_FONT_FAMILY_KORE  = []string{"GulimChe"}
	FONT_CURSIVE_FONT_FAMILY_KORE    = []string{"Gungsuh"}

	FONT_SERIF_FONT_FAMILY_LATN      = []string{"Times New Roman"}
	FONT_SANS_SERIF_FONT_FAMILY_LATN = []string{"Arial"}
	FONT_MONOSPACE_FONT_FAMILY_LATN  = []string{"Consolas"}
	FONT_CURSIVE_FONT_FAMILY_LATN    = []string{"Comic Sans MS"}

	FONT_SERIF_FONT_FAMILY_MLYM      = []string{"Rachana_w01", "AnjaliOldLipi", "Kartika", "ThoolikaUnicode"}
	FONT_SANS_SERIF_FONT_FAMILY_MLYM = []string{"Rachana_w01", "AnjaliOldLipi", "Kartika", "ThoolikaUnicode"}
	FONT_MONOSPACE_FONT_FAMILY_MLYM  = []string{"Rachana_w01", "AnjaliOldLipi", "Kartika", "ThoolikaUnicode"}
	FONT_CURSIVE_FONT_FAMILY_MLYM    = []string{}

	FONT_SERIF_FONT_FAMILY_ORYA      = []string{"ori1Uni", "Kalinga"}
	FONT_SANS_SERIF_FONT_FAMILY_ORYA = []string{"ori1Uni", "Kalinga"}
	FONT_MONOSPACE_FONT_FAMILY_ORYA  = []string{"ori1Uni", "Kalinga"}
	FONT_CURSIVE_FONT_FAMILY_ORYA    = []string{}

	FONT_SERIF_FONT_FAMILY_SINH      = []string{"Iskoola Pota", "AksharUnicode"}
	FONT_SANS_SERIF_FONT_FAMILY_SINH = []string{"Iskoola Pota", "AksharUnicode"}
	FONT_MONOSPACE_FONT_FAMILY_SINH  = []string{"Iskoola Pota", "AksharUnicode"}
	FONT_CURSIVE_FONT_FAMILY_SINH    = []string{}

	FONT_SERIF_FONT_FAMILY_TAML      = []string{"Latha"}
	FONT_SANS_SERIF_FONT_FAMILY_TAML = []string{}
	FONT_MONOSPACE_FONT_FAMILY_TAML  = []string{"Latha"}
	FONT_CURSIVE_FONT_FAMILY_TAML    = []string{}

	FONT_SERIF_FONT_FAMILY_TELU      = []string{"Gautami", "Akshar Unicode"}
	FONT_SANS_SERIF_FONT_FAMILY_TELU = []string{"Gautami", "Akshar Unicode"}
	FONT_MONOSPACE_FONT_FAMILY_TELU  = []string{"Gautami", "Akshar Unicode"}
	FONT_CURSIVE_FONT_FAMILY_TELU    = []string{}

	FONT_SERIF_FONT_FAMILY_THAI      = []string{"Tahoma"}
	FONT_SANS_SERIF_FONT_FAMILY_THAI = []string{"Tahoma"}
	FONT_MONOSPACE_FONT_FAMILY_THAI  = []string{"Tahoma"}
	FONT_CURSIVE_FONT_FAMILY_THAI    = []string{"Tahoma"}

	FONT_SERIF_FONT_FAMILY_TIBT      = []string{"Tibetan Machine Uni", "Jomolhari", "Microsoft Himalaya"}
	FONT_SANS_SERIF_FONT_FAMILY_TIBT = []string{"Tibetan Machine Uni", "Jomolhari", "Microsoft Himalaya"}
	FONT_MONOSPACE_FONT_FAMILY_TIBT  = []string{"Tibetan Machine Uni", "Jomolhari", "Microsoft Himalaya"}
	FONT_CURSIVE_FONT_FAMILY_TIBT    = []string{}
)
