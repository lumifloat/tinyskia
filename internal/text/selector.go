package text

func IsGenericFamily(family string) bool {
	if family == "sans-serif" || family == "serif" || family == "monospace" ||
		family == "cursive" || family == "fantasy" || family == "math" || family == "emoji" {
		return true
	}
	return false
}

func GetPreferenceFonts(family string, script Script) []string {
	switch family {
	case "sans-serif":
		switch script {
		case ScriptArab:
			return FONT_SANS_SERIF_FONT_FAMILY_ARAB
		case ScriptArmi:
			return FONT_SANS_SERIF_FONT_FAMILY_ARMN
		case ScriptBeng:
			return FONT_SANS_SERIF_FONT_FAMILY_BENG
		case ScriptCans:
			return FONT_SANS_SERIF_FONT_FAMILY_CANS
		case ScriptCyrl:
			return FONT_SANS_SERIF_FONT_FAMILY_CYRL
		case ScriptDeva:
			return FONT_SANS_SERIF_FONT_FAMILY_DEVA
		case ScriptEthi:
			return FONT_SANS_SERIF_FONT_FAMILY_ETHI
		case ScriptBopo, ScriptHanb, ScriptHans, ScriptHani:
			return FONT_SANS_SERIF_FONT_FAMILY_HANS
		case ScriptHant, ScriptHntl:
			return FONT_SANS_SERIF_FONT_FAMILY_HANT
		case ScriptGeor, ScriptGeok, ScriptLisu:
			return FONT_SANS_SERIF_FONT_FAMILY_GEOR
		case ScriptGrek:
			return FONT_SANS_SERIF_FONT_FAMILY_GREK
		case ScriptGujr:
			return FONT_SANS_SERIF_FONT_FAMILY_GUJR
		case ScriptGuru:
			return FONT_SANS_SERIF_FONT_FAMILY_GURU
		case ScriptHebr:
			return FONT_SANS_SERIF_FONT_FAMILY_HEBR
		case ScriptHira, ScriptKana, ScriptHrkt, ScriptJpan:
			return FONT_SANS_SERIF_FONT_FAMILY_JPAN
		case ScriptKhmr:
			return FONT_SANS_SERIF_FONT_FAMILY_KHMR
		case ScriptKnda:
			return FONT_SANS_SERIF_FONT_FAMILY_KNDA
		case ScriptHang, ScriptJamo, ScriptKore:
			return FONT_SANS_SERIF_FONT_FAMILY_KORE
		case ScriptMlym:
			return FONT_SANS_SERIF_FONT_FAMILY_MLYM
		case ScriptOrya:
			return FONT_SANS_SERIF_FONT_FAMILY_ORYA
		case ScriptSinh:
			return FONT_SANS_SERIF_FONT_FAMILY_SINH
		case ScriptTaml:
			return FONT_SANS_SERIF_FONT_FAMILY_TAML
		case ScriptTelu:
			return FONT_SANS_SERIF_FONT_FAMILY_TELU
		case ScriptThai:
			return FONT_SANS_SERIF_FONT_FAMILY_THAI
		case ScriptTibt:
			return FONT_SANS_SERIF_FONT_FAMILY_TIBT
		}
		return FONT_SANS_SERIF_FONT_FAMILY_LATN
	case "serif":
		switch script {
		case ScriptArab:
			return FONT_SERIF_FONT_FAMILY_ARAB
		case ScriptArmi:
			return FONT_SERIF_FONT_FAMILY_ARMN
		case ScriptBeng:
			return FONT_SERIF_FONT_FAMILY_BENG
		case ScriptCans:
			return FONT_SERIF_FONT_FAMILY_CANS
		case ScriptCyrl:
			return FONT_SERIF_FONT_FAMILY_CYRL
		case ScriptDeva:
			return FONT_SERIF_FONT_FAMILY_DEVA
		case ScriptEthi:
			return FONT_SERIF_FONT_FAMILY_ETHI
		case ScriptBopo, ScriptHanb, ScriptHans, ScriptHani:
			return FONT_SERIF_FONT_FAMILY_HANS
		case ScriptHant, ScriptHntl:
			return FONT_SERIF_FONT_FAMILY_HANT
		case ScriptGeor, ScriptGeok, ScriptLisu:
			return FONT_SERIF_FONT_FAMILY_GEOR
		case ScriptGrek:
			return FONT_SERIF_FONT_FAMILY_GREK
		case ScriptGujr:
			return FONT_SERIF_FONT_FAMILY_GUJR
		case ScriptGuru:
			return FONT_SERIF_FONT_FAMILY_GURU
		case ScriptHebr:
			return FONT_SERIF_FONT_FAMILY_HEBR
		case ScriptHira, ScriptKana, ScriptHrkt, ScriptJpan:
			return FONT_SERIF_FONT_FAMILY_JPAN
		case ScriptKhmr:
			return FONT_SERIF_FONT_FAMILY_KHMR
		case ScriptKnda:
			return FONT_SERIF_FONT_FAMILY_KNDA
		case ScriptHang, ScriptJamo, ScriptKore:
			return FONT_SERIF_FONT_FAMILY_KORE
		case ScriptMlym:
			return FONT_SERIF_FONT_FAMILY_MLYM
		case ScriptOrya:
			return FONT_SERIF_FONT_FAMILY_ORYA
		case ScriptSinh:
			return FONT_SERIF_FONT_FAMILY_SINH
		case ScriptTaml:
			return FONT_SERIF_FONT_FAMILY_TAML
		case ScriptTelu:
			return FONT_SERIF_FONT_FAMILY_TELU
		case ScriptThai:
			return FONT_SERIF_FONT_FAMILY_THAI
		case ScriptTibt:
			return FONT_SERIF_FONT_FAMILY_TIBT
		}
		return FONT_SERIF_FONT_FAMILY_LATN
	case "monospace":
		switch script {
		case ScriptArab:
			return FONT_MONOSPACE_FONT_FAMILY_ARAB
		case ScriptArmi:
			return FONT_MONOSPACE_FONT_FAMILY_ARMN
		case ScriptBeng:
			return FONT_MONOSPACE_FONT_FAMILY_BENG
		case ScriptCans:
			return FONT_MONOSPACE_FONT_FAMILY_CANS
		case ScriptCyrl:
			return FONT_MONOSPACE_FONT_FAMILY_CYRL
		case ScriptDeva:
			return FONT_MONOSPACE_FONT_FAMILY_DEVA
		case ScriptEthi:
			return FONT_MONOSPACE_FONT_FAMILY_ETHI
		case ScriptBopo, ScriptHanb, ScriptHans, ScriptHani:
			return FONT_MONOSPACE_FONT_FAMILY_HANS
		case ScriptHant, ScriptHntl:
			return FONT_MONOSPACE_FONT_FAMILY_HANT
		case ScriptGeor, ScriptGeok, ScriptLisu:
			return FONT_MONOSPACE_FONT_FAMILY_GEOR
		case ScriptGrek:
			return FONT_MONOSPACE_FONT_FAMILY_GREK
		case ScriptGujr:
			return FONT_MONOSPACE_FONT_FAMILY_GUJR
		case ScriptGuru:
			return FONT_MONOSPACE_FONT_FAMILY_GURU
		case ScriptHebr:
			return FONT_MONOSPACE_FONT_FAMILY_HEBR
		case ScriptHira, ScriptKana, ScriptHrkt, ScriptJpan:
			return FONT_MONOSPACE_FONT_FAMILY_JPAN
		case ScriptKhmr:
			return FONT_MONOSPACE_FONT_FAMILY_KHMR
		case ScriptKnda:
			return FONT_MONOSPACE_FONT_FAMILY_KNDA
		case ScriptHang, ScriptJamo, ScriptKore:
			return FONT_MONOSPACE_FONT_FAMILY_KORE
		case ScriptMlym:
			return FONT_MONOSPACE_FONT_FAMILY_MLYM
		case ScriptOrya:
			return FONT_MONOSPACE_FONT_FAMILY_ORYA
		case ScriptSinh:
			return FONT_MONOSPACE_FONT_FAMILY_SINH
		case ScriptTaml:
			return FONT_MONOSPACE_FONT_FAMILY_TAML
		case ScriptTelu:
			return FONT_MONOSPACE_FONT_FAMILY_TELU
		case ScriptThai:
			return FONT_MONOSPACE_FONT_FAMILY_THAI
		case ScriptTibt:
			return FONT_MONOSPACE_FONT_FAMILY_TIBT
		}
		return FONT_MONOSPACE_FONT_FAMILY_LATN
	case "cursive":
		switch script {
		case ScriptArab:
			return FONT_CURSIVE_FONT_FAMILY_ARAB
		case ScriptArmi:
			return FONT_CURSIVE_FONT_FAMILY_ARMN
		case ScriptBeng:
			return FONT_CURSIVE_FONT_FAMILY_BENG
		case ScriptCans:
			return FONT_CURSIVE_FONT_FAMILY_CANS
		case ScriptCyrl:
			return FONT_CURSIVE_FONT_FAMILY_CYRL
		case ScriptDeva:
			return FONT_CURSIVE_FONT_FAMILY_DEVA
		case ScriptEthi:
			return FONT_CURSIVE_FONT_FAMILY_ETHI
		case ScriptBopo, ScriptHanb, ScriptHans, ScriptHani:
			return FONT_CURSIVE_FONT_FAMILY_HANS
		case ScriptHant, ScriptHntl:
			return FONT_CURSIVE_FONT_FAMILY_HANT
		case ScriptGeor, ScriptGeok, ScriptLisu:
			return FONT_CURSIVE_FONT_FAMILY_GEOR
		case ScriptGrek:
			return FONT_CURSIVE_FONT_FAMILY_GREK
		case ScriptGujr:
			return FONT_CURSIVE_FONT_FAMILY_GUJR
		case ScriptGuru:
			return FONT_CURSIVE_FONT_FAMILY_GURU
		case ScriptHebr:
			return FONT_CURSIVE_FONT_FAMILY_HEBR
		case ScriptHira, ScriptKana, ScriptHrkt, ScriptJpan:
			return FONT_CURSIVE_FONT_FAMILY_JPAN
		case ScriptKhmr:
			return FONT_CURSIVE_FONT_FAMILY_KHMR
		case ScriptKnda:
			return FONT_CURSIVE_FONT_FAMILY_KNDA
		case ScriptHang, ScriptJamo, ScriptKore:
			return FONT_CURSIVE_FONT_FAMILY_KORE
		case ScriptMlym:
			return FONT_CURSIVE_FONT_FAMILY_MLYM
		case ScriptOrya:
			return FONT_CURSIVE_FONT_FAMILY_ORYA
		case ScriptSinh:
			return FONT_CURSIVE_FONT_FAMILY_SINH
		case ScriptTaml:
			return FONT_CURSIVE_FONT_FAMILY_TAML
		case ScriptTelu:
			return FONT_CURSIVE_FONT_FAMILY_TELU
		case ScriptThai:
			return FONT_CURSIVE_FONT_FAMILY_THAI
		case ScriptTibt:
			return FONT_CURSIVE_FONT_FAMILY_TIBT
		}
		return FONT_CURSIVE_FONT_FAMILY_LATN
	case "fantasy":
		return FONT_FANTASY_FONT_FAMILY
	case "math":
		return FONT_MATH_FONT_FAMILY
	case "emoji":
		return FONT_EMOJI_FONT_FAMILY
	default:
		return nil
	}
}
