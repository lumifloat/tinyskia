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
	case ScriptZzzz:
		// Ensure the switch covers all the Script enum values.
		break
	case ScriptZyyy, ScriptZinh, ScriptLatn, ScriptCyrl, ScriptGrek, ScriptArmn, ScriptHebr:
		// We always append Arial below, so no need to add it here.
		break
	case ScriptZmth, ScriptZsym, ScriptZsye:
		// Not currently returned by script run resolution (but some symbols may
		// be handled below).
		break
	case ScriptBopo, ScriptHanb, ScriptHans, ScriptHani:
		list = append(list, "SimSun", "SimSun-ExtB")
	case ScriptHant, ScriptHntl:
		list = append(list, "MingLiU", "MingLiU-ExtB")
	case ScriptHira, ScriptKana, ScriptHrkt, ScriptJpan:
		list = append(list, "Yu Gothic", "MS PGothic")
	case ScriptHang, ScriptJamo, ScriptKore:
		list = append(list, "Malgun Gothic")
	case ScriptYiii:
		list = append(list, "Microsoft Yi Baiti")
	case ScriptMong:
		list = append(list, "Mongolian Baiti")
	case ScriptTibt:
		list = append(list, "Microsoft Himalaya")
	case ScriptPhag:
		list = append(list, "Microsoft PhagsPa")
	case ScriptArab:
		// Default to Arial (added unconditionally below) for Arabic script.
		break
	case ScriptAran:
		list = append(list, "Urdu Typesetting")
	case ScriptSyrc, ScriptSyre:
		list = append(list, "Estrangelo Edessa")
	case ScriptThaa:
		list = append(list, "MV Boli")
	case ScriptBeng:
		list = append(list, "Vrinda", "Nirmala UI")
	case ScriptDeva:
		list = append(list, "Kokila", "Nirmala UI")
	case ScriptGujr:
		list = append(list, "Shruti", "Nirmala UI")
	case ScriptGuru:
		list = append(list, "Raavi", "Nirmala UI")
	case ScriptKnda:
		list = append(list, "Tunga", "Nirmala UI")
	case ScriptMlym:
		list = append(list, "Kartika", "Nirmala UI")
	case ScriptOrya:
		list = append(list, "Kalinga", "Nirmala UI")
	case ScriptTaml:
		list = append(list, "Latha", "Nirmala UI")
	case ScriptTelu:
		list = append(list, "Gautami", "Nirmala UI")
	case ScriptSinh:
		list = append(list, "Iskoola Pota", "Nirmala UI")
	case ScriptCakm, ScriptMtei, ScriptOlck, ScriptSora:
		list = append(list, "Nirmala UI")
	case ScriptMymr:
		list = append(list, "Myanmar Text")
	case ScriptKhmr:
		list = append(list, "Khmer UI")
	case ScriptLaoo:
		list = append(list, "Lao UI")
	case ScriptThai:
		list = append(list, "Tahoma", "Leelawadee UI")
	case ScriptTale:
		list = append(list, "Microsoft Tai Le")
	case ScriptBugi:
		list = append(list, "Leelawadee UI")
	case ScriptTalu:
		list = append(list, "Microsoft New Tai Lue")
	case ScriptJava:
		list = append(list, "Javanese Text")
	case ScriptGeor, ScriptGeok, ScriptLisu:
		list = append(list, "Segoe UI")
	case ScriptEthi:
		list = append(list, "Nyala", "Ebrima")
	case ScriptAdlm, ScriptNkoo, ScriptOsma, ScriptTfng, ScriptVaii:
		list = append(list, "Ebrima")
	case ScriptCans:
		list = append(list, "Euphemia")
	case ScriptCher, ScriptOsge:
		list = append(list, "Gadugi")
	case ScriptBrai, ScriptDsrt:
		list = append(list, "Segoe UI Symbol")
	case ScriptBrah, ScriptCari, ScriptXsux, ScriptCprt, ScriptEgyp, ScriptGlag,
		ScriptGoth, ScriptArmi, ScriptPhli, ScriptPrti, ScriptKhar, ScriptLyci,
		ScriptLydi, ScriptMerc, ScriptOgam, ScriptItal, ScriptXpeo, ScriptSarb,
		ScriptOrkh, ScriptPhnx, ScriptRunr, ScriptShaw, ScriptUgar:
		list = append(list, "Segoe UI Historic")
	// For some scripts where Windows doesn't supply a font by default,
	// there are Noto fonts that users might have installed:
	case ScriptAhom:
		list = append(list, "Noto Serif Ahom")
	case ScriptAvst:
		list = append(list, "Noto Sans Avestan")
	case ScriptBali:
		list = append(list, "Noto Sans Balinese")
	case ScriptBamu:
		list = append(list, "Noto Sans Bamum")
	case ScriptBass:
		list = append(list, "Noto Sans Bassa Vah")
	case ScriptBatk:
		list = append(list, "Noto Sans Batak")
	case ScriptBhks:
		list = append(list, "Noto Sans Bhaiksuki")
	case ScriptBuhd:
		list = append(list, "Noto Sans Buhid")
	case ScriptAghb:
		list = append(list, "Noto Sans Caucasian Albanian")
	case ScriptCham:
		list = append(list, "Noto Sans Cham")
	case ScriptCopt:
		list = append(list, "Noto Sans Coptic")
	case ScriptDupl:
		list = append(list, "Noto Sans Duployan")
	case ScriptElba:
		list = append(list, "Noto Sans Elbasan")
	case ScriptGran:
		list = append(list, "Noto Sans Grantha")
	case ScriptRohg:
		list = append(list, "Noto Sans Hanifi Rohingya")
	case ScriptHano:
		list = append(list, "Noto Sans Hanunoo")
	case ScriptHatr:
		list = append(list, "Noto Sans Hatran")
	case ScriptKthi:
		list = append(list, "Noto Sans Kaithi")
	case ScriptKali:
		list = append(list, "Noto Sans Kayah Li")
	case ScriptKhoj:
		list = append(list, "Noto Sans Khojki")
	case ScriptSind:
		list = append(list, "Noto Sans Khudawadi")
	case ScriptLepc:
		list = append(list, "Noto Sans Lepcha")
	case ScriptLimb:
		list = append(list, "Noto Sans Limbu")
	case ScriptLina:
		list = append(list, "Noto Sans Linear A")
	case ScriptLinb:
		list = append(list, "Noto Sans Linear B")
	case ScriptMahj:
		list = append(list, "Noto Sans Mahajani")
	case ScriptMand:
		list = append(list, "Noto Sans Mandaic")
	case ScriptMani:
		list = append(list, "Noto Sans Manichaean")
	case ScriptMarc:
		list = append(list, "Noto Sans Marchen")
	case ScriptMend:
		list = append(list, "Noto Sans Mende Kikakui")
	case ScriptMero:
		list = append(list, "Noto Sans Meroitic")
	case ScriptPlrd:
		list = append(list, "Noto Sans Miao")
	case ScriptModi:
		list = append(list, "Noto Sans Modi")
	case ScriptMroo:
		list = append(list, "Noto Sans Mro")
	case ScriptMult:
		list = append(list, "Noto Sans Multani")
	case ScriptNbat:
		list = append(list, "Noto Sans Nabataean")
	case ScriptNewa:
		list = append(list, "Noto Sans Newa")
	case ScriptHung:
		list = append(list, "Noto Sans Old Hungarian")
	case ScriptNarb:
		list = append(list, "Noto Sans Old North Arabian")
	case ScriptPerm:
		list = append(list, "Noto Sans Old Permic")
	case ScriptHmng:
		list = append(list, "Noto Sans Pahawh Hmong")
	case ScriptPalm:
		list = append(list, "Noto Sans Palmyrene")
	case ScriptPauc:
		list = append(list, "Noto Sans Pau Cin Hau")
	case ScriptPhlp:
		list = append(list, "Noto Sans Psalter Pahlavi")
	case ScriptRjng:
		list = append(list, "Noto Sans Rejang")
	case ScriptSamr:
		list = append(list, "Noto Sans Samaritan")
	case ScriptSaur:
		list = append(list, "Noto Sans Saurashtra")
	case ScriptShrd:
		list = append(list, "Noto Sans Sharada")
	case ScriptSidd:
		list = append(list, "Noto Sans Siddham")
	case ScriptSund:
		list = append(list, "Noto Sans Sundanese")
	case ScriptSylo:
		list = append(list, "Noto Sans Syloti Nagri")
	case ScriptTglg:
		list = append(list, "Noto Sans Tagalog")
	case ScriptTagb:
		list = append(list, "Noto Sans Tagbanwa")
	case ScriptLana:
		list = append(list, "Noto Sans Tai Tham")
	case ScriptTavt:
		list = append(list, "Noto Sans Tai Viet")
	case ScriptTakr:
		list = append(list, "Noto Sans Takri")
	case ScriptTirh:
		list = append(list, "Noto Sans Tirhuta")
	case ScriptWcho:
		list = append(list, "Noto Sans Wancho")
	case ScriptWara:
		list = append(list, "Noto Sans Warang Citi")
	default:
		break
	}

	// Arial is used as default fallback for system fallback, so always try that.
	list = append(list, "Arial")

	// Symbols/dingbats are generally Script=COMMON but may be resolved to any
	// surrounding script run. So we'll always append a couple of likely fonts
	// for such characters.
	// b := uint32(ch) >> 8
	// if script == ScriptZyyy || // Stray COMMON chars not resolved
	// 	(b >= 0x20 && b <= 0x2b) || b == 0x2e { // BMP symbols/punctuation/etc
	// 	// Segoe UI handles some punctuation/symbols that are missing from many text
	// 	// fonts.
	// 	list = append(list, "Segoe UI", "Segoe UI Symbol", "Cambria Math")
	// }

	// Arial Unicode MS also has lots of glyphs for obscure characters; try it as
	// a last resort, if available.
	list = append(list, "Arial Unicode MS")

	list = append(list, "Segoe UI Emoji", "Twemoji Mozilla")

	return list
}
