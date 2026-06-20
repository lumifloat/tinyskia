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
	case ScriptZyyy, ScriptZinh, ScriptLatn, ScriptCyrl, ScriptGrek:
		// In most cases, COMMON and INHERITED characters will be merged into
		// their context, but if they occur without any specific script context
		// we'll just try common default fonts here.
		list = append(list, "Lucida Grande")
	case ScriptZmth, ScriptZsym, ScriptZsye:
		// Not currently returned by script run resolution (but see below, after
		// the switch).
		break
	case ScriptBopo, ScriptHanb, ScriptHans, ScriptHani:
		// Han-derived scripts are unified; we can't be sure which language
		// font to try first, but standard standard defaults here.
		list = append(list, "Songti SC", "SimSun-ExtB")
	case ScriptHant, ScriptHntl:
		list = append(list, "Songti TC", "MingLiU-ExtB")
	case ScriptHira, ScriptKana, ScriptHrkt, ScriptJpan:
		list = append(list, "Hiragino Sans", "Hiragino Kaku Gothic ProN")
	case ScriptJamo, ScriptKore, ScriptHang:
		list = append(list, "Nanum Gothic", "Apple SD Gothic Neo")
	// For most other scripts, macOS comes with a default font we can use.
	case ScriptArab, ScriptAran:
		list = append(list, "Geeza Pro")
	case ScriptArmn:
		list = append(list, "Mshtakan")
	case ScriptBeng:
		list = append(list, "Bangla Sangam MN")
	case ScriptCher:
		list = append(list, "Plantagenet Cherokee")
	case ScriptCopt:
		list = append(list, "Noto Sans Coptic")
	case ScriptDsrt:
		list = append(list, "Baskerville")
	case ScriptDeva:
		list = append(list, "Devanagari Sangam MN")
	case ScriptEthi:
		list = append(list, "Kefa")
	case ScriptGeor:
		list = append(list, "Helvetica")
	case ScriptGoth:
		list = append(list, "Noto Sans Gothic")
	case ScriptGujr:
		list = append(list, "Gujarati Sangam MN")
	case ScriptGuru:
		list = append(list, "Gurmukhi MN")
	case ScriptHebr:
		list = append(list, "Lucida Grande")
	case ScriptKnda:
		list = append(list, "Kannada MN")
	case ScriptKhmr:
		list = append(list, "Khmer MN")
	case ScriptLaoo:
		list = append(list, "Lao MN")
	case ScriptMlym:
		list = append(list, "Malayalam Sangam MN")
	case ScriptMong:
		list = append(list, "Noto Sans Mongolian")
	case ScriptMymr:
		list = append(list, "Myanmar MN")
	case ScriptOgam:
		list = append(list, "Noto Sans Ogham")
	case ScriptItal:
		list = append(list, "Noto Sans Old Italic")
	case ScriptOrya:
		list = append(list, "Oriya Sangam MN")
	case ScriptRunr:
		list = append(list, "Noto Sans Runic")
	case ScriptSinh:
		list = append(list, "Sinhala Sangam MN")
	case ScriptSyrc:
		list = append(list, "Noto Sans Syriac")
	case ScriptTaml:
		list = append(list, "Tamil MN")
	case ScriptTelu:
		list = append(list, "Telugu MN")
	case ScriptThaa:
		list = append(list, "Noto Sans Thaana")
	case ScriptThai:
		list = append(list, "Thonburi")
	case ScriptTibt:
		list = append(list, "Kailasa")
	case ScriptCans:
		list = append(list, "Euphemia UCAS")
	case ScriptYiii:
		list = append(list, "Noto Sans Yi", "STHeiti")
	case ScriptTglg:
		list = append(list, "Noto Sans Tagalog")
	case ScriptHano:
		list = append(list, "Noto Sans Hanunoo")
	case ScriptBuhd:
		list = append(list, "Noto Sans Buhid")
	case ScriptTagb:
		list = append(list, "Noto Sans Tagbanwa")
	case ScriptBrai:
		list = append(list, "Apple Braille")
	case ScriptCprt:
		list = append(list, "Noto Sans Cypriot")
	case ScriptLimb:
		list = append(list, "Noto Sans Limbu")
	case ScriptLinb:
		list = append(list, "Noto Sans Linear B")
	case ScriptOsma:
		list = append(list, "Noto Sans Osmanya")
	case ScriptShaw:
		list = append(list, "Noto Sans Shavian")
	case ScriptTale:
		list = append(list, "Noto Sans Tai Le")
	case ScriptUgar:
		list = append(list, "Noto Sans Ugaritic")
	case ScriptBugi:
		list = append(list, "Noto Sans Buginese")
	case ScriptGlag:
		list = append(list, "Noto Sans Glagolitic")
	case ScriptKhar:
		list = append(list, "Noto Sans Kharoshthi")
	case ScriptSylo:
		list = append(list, "Noto Sans Syloti Nagri")
	case ScriptTalu:
		list = append(list, "Noto Sans New Tai Lue")
	case ScriptTfng:
		list = append(list, "Noto Sans Tifinagh")
	case ScriptXpeo:
		list = append(list, "Noto Sans Old Persian")
	case ScriptBali:
		list = append(list, "Noto Sans Balinese")
	case ScriptBatk:
		list = append(list, "Noto Sans Batak")
	case ScriptBrah:
		list = append(list, "Noto Sans Brahmi")
	case ScriptCham:
		list = append(list, "Noto Sans Cham")
	case ScriptEgyp:
		list = append(list, "Noto Sans Egyptian Hieroglyphs")
	case ScriptHmng:
		list = append(list, "Noto Sans Pahawh Hmong")
	case ScriptHung:
		list = append(list, "Noto Sans Old Hungarian")
	case ScriptJava:
		list = append(list, "Noto Sans Javanese")
	case ScriptKali:
		list = append(list, "Noto Sans Kayah Li")
	case ScriptLepc:
		list = append(list, "Noto Sans Lepcha")
	case ScriptLina:
		list = append(list, "Noto Sans Linear A")
	case ScriptMand:
		list = append(list, "Noto Sans Mandaic")
	case ScriptNkoo:
		list = append(list, "Noto Sans NKo")
	case ScriptOrkh:
		list = append(list, "Noto Sans Old Turkic")
	case ScriptPerm:
		list = append(list, "Noto Sans Old Permic")
	case ScriptPhag:
		list = append(list, "Noto Sans PhagsPa")
	case ScriptPhnx:
		list = append(list, "Noto Sans Phoenician")
	case ScriptPlrd:
		list = append(list, "Noto Sans Miao")
	case ScriptVaii:
		list = append(list, "Noto Sans Vai")
	case ScriptXsux:
		list = append(list, "Noto Sans Cuneiform")
	case ScriptCari:
		list = append(list, "Noto Sans Carian")
	case ScriptLana:
		list = append(list, "Noto Sans Tai Tham")
	case ScriptLyci:
		list = append(list, "Noto Sans Lycian")
	case ScriptLydi:
		list = append(list, "Noto Sans Lydian")
	case ScriptOlck:
		list = append(list, "Noto Sans Ol Chiki")
	case ScriptRjng:
		list = append(list, "Noto Sans Rejang")
	case ScriptSaur:
		list = append(list, "Noto Sans Saurashtra")
	case ScriptSund:
		list = append(list, "Noto Sans Sundanese")
	case ScriptMtei:
		list = append(list, "Noto Sans Meetei Mayek")
	case ScriptArmi:
		list = append(list, "Noto Sans Imperial Aramaic")
	case ScriptAvst:
		list = append(list, "Noto Sans Avestan")
	case ScriptCakm:
		list = append(list, "Noto Sans Chakma")
	case ScriptKthi:
		list = append(list, "Noto Sans Kaithi")
	case ScriptMani:
		list = append(list, "Noto Sans Manichaean")
	case ScriptPhli:
		list = append(list, "Noto Sans Inscriptional Pahlavi")
	case ScriptPhlp:
		list = append(list, "Noto Sans Psalter Pahlavi")
	case ScriptPrti:
		list = append(list, "Noto Sans Inscriptional Parthian")
	case ScriptSamr:
		list = append(list, "Noto Sans Samaritan")
	case ScriptTavt:
		list = append(list, "Noto Sans Tai Viet")
	case ScriptBamu:
		list = append(list, "Noto Sans Bamum")
	case ScriptLisu:
		list = append(list, "Noto Sans Lisu")
	case ScriptSarb:
		list = append(list, "Noto Sans Old South Arabian")
	case ScriptBass:
		list = append(list, "Noto Sans Bassa Vah")
	case ScriptDupl:
		list = append(list, "Noto Sans Duployan")
	case ScriptElba:
		list = append(list, "Noto Sans Elbasan")
	case ScriptGran:
		list = append(list, "Noto Sans Grantha")
	case ScriptMend:
		list = append(list, "Noto Sans Mende Kikakui")
	case ScriptMerc, ScriptMero:
		list = append(list, "Noto Sans Meroitic")
	case ScriptNarb:
		list = append(list, "Noto Sans Old North Arabian")
	case ScriptNbat:
		list = append(list, "Noto Sans Nabataean")
	case ScriptPalm:
		list = append(list, "Noto Sans Palmyrene")
	case ScriptSind:
		list = append(list, "Noto Sans Khudawadi")
	case ScriptWara:
		list = append(list, "Noto Sans Warang Citi")
	case ScriptMroo:
		list = append(list, "Noto Sans Mro")
	case ScriptShrd:
		list = append(list, "Noto Sans Sharada")
	case ScriptSora:
		list = append(list, "Noto Sans Sora Sompeng")
	case ScriptTakr:
		list = append(list, "Noto Sans Takri")
	case ScriptKhoj:
		list = append(list, "Noto Sans Khojki")
	case ScriptTirh:
		list = append(list, "Noto Sans Tirhuta")
	case ScriptAghb:
		list = append(list, "Noto Sans Caucasian Albanian")
	case ScriptMahj:
		list = append(list, "Noto Sans Mahajani")
	case ScriptAhom:
		list = append(list, "Noto Serif Ahom")
	case ScriptHatr:
		list = append(list, "Noto Sans Hatran")
	case ScriptModi:
		list = append(list, "Noto Sans Modi")
	case ScriptMult:
		list = append(list, "Noto Sans Multani")
	case ScriptPauc:
		list = append(list, "Noto Sans Pau Cin Hau")
	case ScriptSidd:
		list = append(list, "Noto Sans Siddham")
	case ScriptAdlm:
		list = append(list, "Noto Sans Adlam")
	case ScriptBhks:
		list = append(list, "Noto Sans Bhaiksuki")
	case ScriptMarc:
		list = append(list, "Noto Sans Marchen")
	case ScriptNewa:
		list = append(list, "Noto Sans Newa")
	case ScriptOsge:
		list = append(list, "Noto Sans Osage")
	case ScriptRohg:
		list = append(list, "Noto Sans Hanifi Rohingya")
	case ScriptWcho:
		list = append(list, "Noto Sans Wancho")
	default:
		break
	}

	// b := uint32(ch) >> 8
	// if script == ScriptZyyy ||
	// 	(b >= 0x20 && b <= 0x2b) || b == 0x2e {
	// 	if b == 0x27 {
	// 		list = append(list, "Zapf Dingbats")
	// 	}
	// 	list = append(list, "Geneva", "STIXGeneral", "Apple Symbols")
	// 	list = append(list, "Hiragino Sans", "Hiragino Kaku Gothic ProN")
	// }

	// Arial Unicode MS has lots of glyphs for obscure characters; try it as a
	// last resort.
	list = append(list, "Arial Unicode MS")

	list = append(list, "Apple Color Emoji")

	return list
}
