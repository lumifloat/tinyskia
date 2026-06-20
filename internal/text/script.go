package text

import (
	"sync"
	"unicode"
)

type Script uint16

const (
	ScriptPcun Script = 15
	ScriptPelm Script = 16
	ScriptXsux Script = 20
	ScriptXpeo Script = 30
	ScriptUgar Script = 40
	ScriptEgyp Script = 50
	ScriptEgyh Script = 60
	ScriptEgyd Script = 70
	ScriptHluw Script = 80
	ScriptNkdb Script = 85
	ScriptMaya Script = 90
	ScriptSgnw Script = 95
	ScriptMero Script = 100
	ScriptMerc Script = 101
	ScriptPsin Script = 103
	ScriptSarb Script = 105
	ScriptNarb Script = 106
	ScriptChrs Script = 109
	ScriptPhnx Script = 115
	ScriptLydi Script = 116
	ScriptTfng Script = 120
	ScriptSamr Script = 123
	ScriptArmi Script = 124
	ScriptHebr Script = 125
	ScriptPalm Script = 126
	ScriptHatr Script = 127
	ScriptElym Script = 128
	ScriptPrti Script = 130
	ScriptPhli Script = 131
	ScriptPhlp Script = 132
	ScriptPhlv Script = 133
	ScriptAvst Script = 134
	ScriptSyrc Script = 135
	ScriptSyrn Script = 136
	ScriptSyrj Script = 137
	ScriptSyre Script = 138
	ScriptMani Script = 139
	ScriptMand Script = 140
	ScriptSogd Script = 141
	ScriptSogo Script = 142
	ScriptOugr Script = 143
	ScriptMong Script = 145
	ScriptNbat Script = 159
	ScriptArab Script = 160
	ScriptAran Script = 161
	ScriptGara Script = 164
	ScriptNkoo Script = 165
	ScriptAdlm Script = 166
	ScriptRohg Script = 167
	ScriptThaa Script = 170
	ScriptOrkh Script = 175
	ScriptHung Script = 176
	ScriptSidt Script = 180
	ScriptYezi Script = 192
	ScriptGrek Script = 200
	ScriptCari Script = 201
	ScriptLyci Script = 202
	ScriptCopt Script = 204
	ScriptGoth Script = 206
	ScriptItal Script = 210
	ScriptRunr Script = 211
	ScriptOgam Script = 212
	ScriptLatn Script = 215
	ScriptLatg Script = 216
	ScriptLatf Script = 217
	ScriptMoon Script = 218
	ScriptOsge Script = 219
	ScriptCyrl Script = 220
	ScriptCyrs Script = 221
	ScriptGlag Script = 225
	ScriptElba Script = 226
	ScriptPerm Script = 227
	ScriptVith Script = 228
	ScriptTodr Script = 229
	ScriptArmn Script = 230
	ScriptAghb Script = 239
	ScriptGeor Script = 240
	ScriptGeok Script = 241
	ScriptDsrt Script = 250
	ScriptBerf Script = 258
	ScriptBass Script = 259
	ScriptOsma Script = 260
	ScriptOlck Script = 261
	ScriptWara Script = 262
	ScriptPauc Script = 263
	ScriptMroo Script = 264
	ScriptMedf Script = 265
	ScriptSunu Script = 274
	ScriptTnsa Script = 275
	ScriptVisp Script = 280
	ScriptShaw Script = 281
	ScriptPlrd Script = 282
	ScriptWcho Script = 283
	ScriptJamo Script = 284
	ScriptBopo Script = 285
	ScriptHang Script = 286
	ScriptKore Script = 287
	ScriptKits Script = 288
	ScriptTeng Script = 290
	ScriptCirt Script = 291
	ScriptSara Script = 292
	ScriptPiqd Script = 293
	ScriptToto Script = 294
	ScriptNagm Script = 295
	ScriptOnao Script = 296
	ScriptChis Script = 298
	ScriptTols Script = 299
	ScriptBrah Script = 300
	ScriptSidd Script = 302
	ScriptRanj Script = 303
	ScriptKhar Script = 305
	ScriptGuru Script = 310
	ScriptNand Script = 311
	ScriptGong Script = 312
	ScriptGonm Script = 313
	ScriptMahj Script = 314
	ScriptDeva Script = 315
	ScriptSylo Script = 316
	ScriptKthi Script = 317
	ScriptSind Script = 318
	ScriptShrd Script = 319
	ScriptGujr Script = 320
	ScriptTakr Script = 321
	ScriptKhoj Script = 322
	ScriptMult Script = 323
	ScriptModi Script = 324
	ScriptBeng Script = 325
	ScriptTirh Script = 326
	ScriptOrya Script = 327
	ScriptDogr Script = 328
	ScriptSoyo Script = 329
	ScriptTibt Script = 330
	ScriptPhag Script = 331
	ScriptMarc Script = 332
	ScriptNewa Script = 333
	ScriptBhks Script = 334
	ScriptLepc Script = 335
	ScriptLimb Script = 336
	ScriptMtei Script = 337
	ScriptAhom Script = 338
	ScriptZanb Script = 339
	ScriptTelu Script = 340
	ScriptTutg Script = 341
	ScriptDiak Script = 342
	ScriptGran Script = 343
	ScriptSaur Script = 344
	ScriptKnda Script = 345
	ScriptTaml Script = 346
	ScriptMlym Script = 347
	ScriptSinh Script = 348
	ScriptCakm Script = 349
	ScriptMymr Script = 350
	ScriptLana Script = 351
	ScriptThai Script = 352
	ScriptTale Script = 353
	ScriptTalu Script = 354
	ScriptKhmr Script = 355
	ScriptLaoo Script = 356
	ScriptKali Script = 357
	ScriptCham Script = 358
	ScriptTavt Script = 359
	ScriptBali Script = 360
	ScriptJava Script = 361
	ScriptSund Script = 362
	ScriptRjng Script = 363
	ScriptLeke Script = 364
	ScriptBatk Script = 365
	ScriptMaka Script = 366
	ScriptBugi Script = 367
	ScriptKawi Script = 368
	ScriptTglg Script = 370
	ScriptHano Script = 371
	ScriptBuhd Script = 372
	ScriptTagb Script = 373
	ScriptTayo Script = 380
	ScriptKrai Script = 396
	ScriptGukh Script = 397
	ScriptSora Script = 398
	ScriptLisu Script = 399
	ScriptLina Script = 400
	ScriptLinb Script = 401
	ScriptCpmn Script = 402
	ScriptCprt Script = 403
	ScriptHira Script = 410
	ScriptKana Script = 411
	ScriptHrkt Script = 412
	ScriptJpan Script = 413
	ScriptNkgb Script = 420
	ScriptEthi Script = 430
	ScriptBamu Script = 435
	ScriptKpel Script = 436
	ScriptLoma Script = 437
	ScriptMend Script = 438
	ScriptAfak Script = 439
	ScriptCans Script = 440
	ScriptCher Script = 445
	ScriptHmng Script = 450
	ScriptHmnp Script = 451
	ScriptYiii Script = 460
	ScriptVaii Script = 470
	ScriptWole Script = 480
	ScriptNshu Script = 499
	ScriptHani Script = 500
	ScriptHans Script = 501
	ScriptHant Script = 502
	ScriptHanb Script = 503
	ScriptHntl Script = 504
	ScriptKitl Script = 505
	ScriptJurc Script = 510
	ScriptTang Script = 520
	ScriptShui Script = 530
	ScriptBlis Script = 550
	ScriptBrai Script = 570
	ScriptSeal Script = 590
	ScriptInds Script = 610
	ScriptRoro Script = 620
	ScriptDupl Script = 755
	ScriptQaaa Script = 900
	ScriptQabx Script = 949
	ScriptZsye Script = 993
	ScriptZinh Script = 994
	ScriptZmth Script = 995
	ScriptZsym Script = 996
	ScriptZxxx Script = 997
	ScriptZyyy Script = 998
	ScriptZzzz Script = 999
)

var Scripts = map[string]Script{
	"Pcun": ScriptPcun,
	"Pelm": ScriptPelm,
	"Xsux": ScriptXsux,
	"Xpeo": ScriptXpeo,
	"Ugar": ScriptUgar,
	"Egyp": ScriptEgyp,
	"Egyh": ScriptEgyh,
	"Egyd": ScriptEgyd,
	"Hluw": ScriptHluw,
	"Nkdb": ScriptNkdb,
	"Maya": ScriptMaya,
	"Sgnw": ScriptSgnw,
	"Mero": ScriptMero,
	"Merc": ScriptMerc,
	"Psin": ScriptPsin,
	"Sarb": ScriptSarb,
	"Narb": ScriptNarb,
	"Chrs": ScriptChrs,
	"Phnx": ScriptPhnx,
	"Lydi": ScriptLydi,
	"Tfng": ScriptTfng,
	"Samr": ScriptSamr,
	"Armi": ScriptArmi,
	"Hebr": ScriptHebr,
	"Palm": ScriptPalm,
	"Hatr": ScriptHatr,
	"Elym": ScriptElym,
	"Prti": ScriptPrti,
	"Phli": ScriptPhli,
	"Phlp": ScriptPhlp,
	"Phlv": ScriptPhlv,
	"Avst": ScriptAvst,
	"Syrc": ScriptSyrc,
	"Syrn": ScriptSyrn,
	"Syrj": ScriptSyrj,
	"Syre": ScriptSyre,
	"Mani": ScriptMani,
	"Mand": ScriptMand,
	"Sogd": ScriptSogd,
	"Sogo": ScriptSogo,
	"Ougr": ScriptOugr,
	"Mong": ScriptMong,
	"Nbat": ScriptNbat,
	"Arab": ScriptArab,
	"Aran": ScriptAran,
	"Gara": ScriptGara,
	"Nkoo": ScriptNkoo,
	"Adlm": ScriptAdlm,
	"Rohg": ScriptRohg,
	"Thaa": ScriptThaa,
	"Orkh": ScriptOrkh,
	"Hung": ScriptHung,
	"Sidt": ScriptSidt,
	"Yezi": ScriptYezi,
	"Grek": ScriptGrek,
	"Cari": ScriptCari,
	"Lyci": ScriptLyci,
	"Copt": ScriptCopt,
	"Goth": ScriptGoth,
	"Ital": ScriptItal,
	"Runr": ScriptRunr,
	"Ogam": ScriptOgam,
	"Latn": ScriptLatn,
	"Latg": ScriptLatg,
	"Latf": ScriptLatf,
	"Moon": ScriptMoon,
	"Osge": ScriptOsge,
	"Cyrl": ScriptCyrl,
	"Cyrs": ScriptCyrs,
	"Glag": ScriptGlag,
	"Elba": ScriptElba,
	"Perm": ScriptPerm,
	"Vith": ScriptVith,
	"Todr": ScriptTodr,
	"Armn": ScriptArmn,
	"Aghb": ScriptAghb,
	"Geor": ScriptGeor,
	"Geok": ScriptGeok,
	"Dsrt": ScriptDsrt,
	"Berf": ScriptBerf,
	"Bass": ScriptBass,
	"Osma": ScriptOsma,
	"Olck": ScriptOlck,
	"Wara": ScriptWara,
	"Pauc": ScriptPauc,
	"Mroo": ScriptMroo,
	"Medf": ScriptMedf,
	"Sunu": ScriptSunu,
	"Tnsa": ScriptTnsa,
	"Visp": ScriptVisp,
	"Shaw": ScriptShaw,
	"Plrd": ScriptPlrd,
	"Wcho": ScriptWcho,
	"Jamo": ScriptJamo,
	"Bopo": ScriptBopo,
	"Hang": ScriptHang,
	"Kore": ScriptKore,
	"Kits": ScriptKits,
	"Teng": ScriptTeng,
	"Cirt": ScriptCirt,
	"Sara": ScriptSara,
	"Piqd": ScriptPiqd,
	"Toto": ScriptToto,
	"Nagm": ScriptNagm,
	"Onao": ScriptOnao,
	"Chis": ScriptChis,
	"Tols": ScriptTols,
	"Brah": ScriptBrah,
	"Sidd": ScriptSidd,
	"Ranj": ScriptRanj,
	"Khar": ScriptKhar,
	"Guru": ScriptGuru,
	"Nand": ScriptNand,
	"Gong": ScriptGong,
	"Gonm": ScriptGonm,
	"Mahj": ScriptMahj,
	"Deva": ScriptDeva,
	"Sylo": ScriptSylo,
	"Kthi": ScriptKthi,
	"Sind": ScriptSind,
	"Shrd": ScriptShrd,
	"Gujr": ScriptGujr,
	"Takr": ScriptTakr,
	"Khoj": ScriptKhoj,
	"Mult": ScriptMult,
	"Modi": ScriptModi,
	"Beng": ScriptBeng,
	"Tirh": ScriptTirh,
	"Orya": ScriptOrya,
	"Dogr": ScriptDogr,
	"Soyo": ScriptSoyo,
	"Tibt": ScriptTibt,
	"Phag": ScriptPhag,
	"Marc": ScriptMarc,
	"Newa": ScriptNewa,
	"Bhks": ScriptBhks,
	"Lepc": ScriptLepc,
	"Limb": ScriptLimb,
	"Mtei": ScriptMtei,
	"Ahom": ScriptAhom,
	"Zanb": ScriptZanb,
	"Telu": ScriptTelu,
	"Tutg": ScriptTutg,
	"Diak": ScriptDiak,
	"Gran": ScriptGran,
	"Saur": ScriptSaur,
	"Knda": ScriptKnda,
	"Taml": ScriptTaml,
	"Mlym": ScriptMlym,
	"Sinh": ScriptSinh,
	"Cakm": ScriptCakm,
	"Mymr": ScriptMymr,
	"Lana": ScriptLana,
	"Thai": ScriptThai,
	"Tale": ScriptTale,
	"Talu": ScriptTalu,
	"Khmr": ScriptKhmr,
	"Laoo": ScriptLaoo,
	"Kali": ScriptKali,
	"Cham": ScriptCham,
	"Tavt": ScriptTavt,
	"Bali": ScriptBali,
	"Java": ScriptJava,
	"Sund": ScriptSund,
	"Rjng": ScriptRjng,
	"Leke": ScriptLeke,
	"Batk": ScriptBatk,
	"Maka": ScriptMaka,
	"Bugi": ScriptBugi,
	"Kawi": ScriptKawi,
	"Tglg": ScriptTglg,
	"Hano": ScriptHano,
	"Buhd": ScriptBuhd,
	"Tagb": ScriptTagb,
	"Tayo": ScriptTayo,
	"Krai": ScriptKrai,
	"Gukh": ScriptGukh,
	"Sora": ScriptSora,
	"Lisu": ScriptLisu,
	"Lina": ScriptLina,
	"Linb": ScriptLinb,
	"Cpmn": ScriptCpmn,
	"Cprt": ScriptCprt,
	"Hira": ScriptHira,
	"Kana": ScriptKana,
	"Hrkt": ScriptHrkt,
	"Jpan": ScriptJpan,
	"Nkgb": ScriptNkgb,
	"Ethi": ScriptEthi,
	"Bamu": ScriptBamu,
	"Kpel": ScriptKpel,
	"Loma": ScriptLoma,
	"Mend": ScriptMend,
	"Afak": ScriptAfak,
	"Cans": ScriptCans,
	"Cher": ScriptCher,
	"Hmng": ScriptHmng,
	"Hmnp": ScriptHmnp,
	"Yiii": ScriptYiii,
	"Vaii": ScriptVaii,
	"Wole": ScriptWole,
	"Nshu": ScriptNshu,
	"Hani": ScriptHani,
	"Hans": ScriptHans,
	"Hant": ScriptHant,
	"Hanb": ScriptHanb,
	"Hntl": ScriptHntl,
	"Kitl": ScriptKitl,
	"Jurc": ScriptJurc,
	"Tang": ScriptTang,
	"Shui": ScriptShui,
	"Blis": ScriptBlis,
	"Brai": ScriptBrai,
	"Seal": ScriptSeal,
	"Inds": ScriptInds,
	"Roro": ScriptRoro,
	"Dupl": ScriptDupl,
	"Qaaa": ScriptQaaa,
	"Qabx": ScriptQabx,
	"Zsye": ScriptZsye,
	"Zinh": ScriptZinh,
	"Zmth": ScriptZmth,
	"Zsym": ScriptZsym,
	"Zxxx": ScriptZxxx,
	"Zyyy": ScriptZyyy,
	"Zzzz": ScriptZzzz,
}

type Plane [256 * 256]Script

type ScriptCache struct {
	planes [17]*Plane
	locks  [17]sync.Once
}

var table = map[Script]*unicode.RangeTable{
	ScriptAdlm: unicode.Adlam,
	ScriptAhom: unicode.Ahom,
	ScriptHluw: unicode.Anatolian_Hieroglyphs,
	ScriptArab: unicode.Arabic,
	ScriptArmn: unicode.Armenian,
	ScriptAvst: unicode.Avestan,
	ScriptBali: unicode.Balinese,
	ScriptBamu: unicode.Bamum,
	ScriptBass: unicode.Bassa_Vah,
	ScriptBatk: unicode.Batak,
	ScriptBeng: unicode.Bengali,
	ScriptBhks: unicode.Bhaiksuki,
	ScriptBopo: unicode.Bopomofo,
	ScriptBrah: unicode.Brahmi,
	ScriptBrai: unicode.Braille,
	ScriptBugi: unicode.Buginese,
	ScriptBuhd: unicode.Buhid,
	ScriptCans: unicode.Canadian_Aboriginal,
	ScriptCari: unicode.Carian,
	ScriptAghb: unicode.Caucasian_Albanian,
	ScriptCakm: unicode.Chakma,
	ScriptCham: unicode.Cham,
	ScriptCher: unicode.Cherokee,
	ScriptChrs: unicode.Chorasmian,
	ScriptZyyy: unicode.Common,
	ScriptCopt: unicode.Coptic,
	ScriptXsux: unicode.Cuneiform,
	ScriptCprt: unicode.Cypriot,
	ScriptCpmn: unicode.Cypro_Minoan,
	ScriptCyrl: unicode.Cyrillic,
	ScriptDsrt: unicode.Deseret,
	ScriptDeva: unicode.Devanagari,
	ScriptDiak: unicode.Dives_Akuru,
	ScriptDogr: unicode.Dogra,
	ScriptDupl: unicode.Duployan,
	ScriptEgyp: unicode.Egyptian_Hieroglyphs,
	ScriptElba: unicode.Elbasan,
	ScriptElym: unicode.Elymaic,
	ScriptEthi: unicode.Ethiopic,
	ScriptGeor: unicode.Georgian,
	ScriptGlag: unicode.Glagolitic,
	ScriptGoth: unicode.Gothic,
	ScriptGran: unicode.Grantha,
	ScriptGrek: unicode.Greek,
	ScriptGujr: unicode.Gujarati,
	ScriptGong: unicode.Gunjala_Gondi,
	ScriptGuru: unicode.Gurmukhi,
	ScriptHani: unicode.Han,
	ScriptHang: unicode.Hangul,
	ScriptRohg: unicode.Hanifi_Rohingya,
	ScriptHano: unicode.Hanunoo,
	ScriptHatr: unicode.Hatran,
	ScriptHebr: unicode.Hebrew,
	ScriptHira: unicode.Hiragana,
	ScriptArmi: unicode.Imperial_Aramaic,
	ScriptZinh: unicode.Inherited,
	ScriptPhli: unicode.Inscriptional_Pahlavi,
	ScriptPrti: unicode.Inscriptional_Parthian,
	ScriptJava: unicode.Javanese,
	ScriptKthi: unicode.Kaithi,
	ScriptKnda: unicode.Kannada,
	ScriptKana: unicode.Katakana,
	ScriptKawi: unicode.Kawi,
	ScriptKali: unicode.Kayah_Li,
	ScriptKhar: unicode.Kharoshthi,
	ScriptKits: unicode.Khitan_Small_Script,
	ScriptKhmr: unicode.Khmer,
	ScriptKhoj: unicode.Khojki,
	ScriptSind: unicode.Khudawadi,
	ScriptLaoo: unicode.Lao,
	ScriptLatn: unicode.Latin,
	ScriptLepc: unicode.Lepcha,
	ScriptLimb: unicode.Limbu,
	ScriptLina: unicode.Linear_A,
	ScriptLinb: unicode.Linear_B,
	ScriptLisu: unicode.Lisu,
	ScriptLyci: unicode.Lycian,
	ScriptLydi: unicode.Lydian,
	ScriptMahj: unicode.Mahajani,
	ScriptMaka: unicode.Makasar,
	ScriptMlym: unicode.Malayalam,
	ScriptMand: unicode.Mandaic,
	ScriptMani: unicode.Manichaean,
	ScriptMarc: unicode.Marchen,
	ScriptGonm: unicode.Masaram_Gondi,
	ScriptMedf: unicode.Medefaidrin,
	ScriptMtei: unicode.Meetei_Mayek,
	ScriptMend: unicode.Mende_Kikakui,
	ScriptMerc: unicode.Meroitic_Cursive,
	ScriptMero: unicode.Meroitic_Hieroglyphs,
	ScriptPlrd: unicode.Miao,
	ScriptModi: unicode.Modi,
	ScriptMong: unicode.Mongolian,
	ScriptMroo: unicode.Mro,
	ScriptMult: unicode.Multani,
	ScriptMymr: unicode.Myanmar,
	ScriptNbat: unicode.Nabataean,
	ScriptNagm: unicode.Nag_Mundari,
	ScriptNand: unicode.Nandinagari,
	ScriptTalu: unicode.New_Tai_Lue,
	ScriptNewa: unicode.Newa,
	ScriptNkoo: unicode.Nko,
	ScriptNshu: unicode.Nushu,
	ScriptHmnp: unicode.Nyiakeng_Puachue_Hmong,
	ScriptOgam: unicode.Ogham,
	ScriptOlck: unicode.Ol_Chiki,
	ScriptHung: unicode.Old_Hungarian,
	ScriptItal: unicode.Old_Italic,
	ScriptNarb: unicode.Old_North_Arabian,
	ScriptPerm: unicode.Old_Permic,
	ScriptXpeo: unicode.Old_Persian,
	ScriptSogo: unicode.Old_Sogdian,
	ScriptSarb: unicode.Old_South_Arabian,
	ScriptOrkh: unicode.Old_Turkic,
	ScriptOugr: unicode.Old_Uyghur,
	ScriptOrya: unicode.Oriya,
	ScriptOsge: unicode.Osage,
	ScriptOsma: unicode.Osmanya,
	ScriptHmng: unicode.Pahawh_Hmong,
	ScriptPalm: unicode.Palmyrene,
	ScriptPauc: unicode.Pau_Cin_Hau,
	ScriptPhag: unicode.Phags_Pa,
	ScriptPhnx: unicode.Phoenician,
	ScriptPhlv: unicode.Psalter_Pahlavi,
	ScriptRjng: unicode.Rejang,
	ScriptRunr: unicode.Runic,
	ScriptSamr: unicode.Samaritan,
	ScriptSaur: unicode.Saurashtra,
	ScriptShrd: unicode.Sharada,
	ScriptShaw: unicode.Shavian,
	ScriptSidd: unicode.Siddham,
	ScriptSgnw: unicode.SignWriting,
	ScriptSinh: unicode.Sinhala,
	ScriptSogd: unicode.Sogdian,
	ScriptSora: unicode.Sora_Sompeng,
	ScriptSoyo: unicode.Soyombo,
	ScriptSund: unicode.Sundanese,
	ScriptSylo: unicode.Syloti_Nagri,
	ScriptSyrc: unicode.Syriac,
	ScriptTglg: unicode.Tagalog,
	ScriptTagb: unicode.Tagbanwa,
	ScriptTale: unicode.Tai_Le,
	ScriptLana: unicode.Tai_Tham,
	ScriptTavt: unicode.Tai_Viet,
	ScriptTakr: unicode.Takri,
	ScriptTaml: unicode.Tamil,
	ScriptTnsa: unicode.Tangsa,
	ScriptTang: unicode.Tangut,
	ScriptTelu: unicode.Telugu,
	ScriptThaa: unicode.Thaana,
	ScriptThai: unicode.Thai,
	ScriptTibt: unicode.Tibetan,
	ScriptTfng: unicode.Tifinagh,
	ScriptTirh: unicode.Tirhuta,
	ScriptToto: unicode.Toto,
	ScriptUgar: unicode.Ugaritic,
	ScriptVaii: unicode.Vai,
	ScriptVith: unicode.Vithkuqi,
	ScriptWcho: unicode.Wancho,
	ScriptWara: unicode.Warang_Citi,
	ScriptYezi: unicode.Yezidi,
	ScriptYiii: unicode.Yi,
	ScriptZanb: unicode.Zanabazar_Square,
}

func (sc *ScriptCache) load(p uint32) {
	sc.locks[p].Do(func() {
		plane := &Plane{}

		planeStart := p << 16
		planeEnd := planeStart | 0xFFFF

		if p == 0 {
			for id, rt := range table {
				for _, r16 := range rt.R16 {
					if uint32(r16.Lo) > planeEnd {
						continue
					}
					stride := uint32(r16.Stride)
					for cp := uint32(r16.Lo); cp <= uint32(r16.Hi); cp += stride {
						plane[cp&0xFFFF] = id
					}
				}
			}
		} else {
			for id, rt := range table {
				for _, r32 := range rt.R32 {
					if r32.Hi < planeStart || r32.Lo > planeEnd {
						continue
					}
					stride := uint32(r32.Stride)
					for cp := r32.Lo; cp <= r32.Hi; cp += stride {
						if cp>>16 == p {
							plane[cp&0xFFFF] = id
						}
					}
				}
			}
		}

		sc.planes[p] = plane
	})
}

func (sc *ScriptCache) GetScript(r rune) Script {
	cp := uint32(r)
	if cp > 0x10FFFF {
		return ScriptZzzz
	}

	p := cp >> 16
	idx := cp & 0xFFFF

	if sc.planes[p] == nil {
		sc.load(p)
	}

	return sc.planes[p][idx]
}

func (sc *ScriptCache) Segement(s []rune) ([][]rune, []Script) {
	if len(s) == 0 {
		return nil, nil
	}

	var fragments [][]rune
	var scripts []Script

	runes := s

	start := 0
	curr := sc.GetScript(runes[0])

	for i := 1; i < len(runes); i++ {
		next := sc.GetScript(runes[i])

		if next != curr {
			fragments = append(fragments, runes[start:i])
			scripts = append(scripts, curr)

			start = i
			curr = next
		}
	}

	fragments = append(fragments, runes[start:])
	scripts = append(scripts, curr)

	return fragments, scripts
}
