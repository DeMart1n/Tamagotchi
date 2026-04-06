package ui

// Sprites half-block para inimigos e elementos das dungeons.
// Usa o mesmo sistema de renderização de sprites.go.
//
// Paleta estendida para inimigos:
//   K = contorno    W = branco     G = verde       L = verde claro
//   Y = amarelo     R = vermelho   P = rosa        B = azul
//   C = ciano       O = laranja    D = cinza       M = marrom
//   . = transparente
//
// Cada inimigo tem design tematicamente coerente com seu nome.

// EnemySprite retorna o sprite pixel art renderizado para um inimigo.
func EnemySprite(enemyName string) string {
	frames, ok := enemySpriteMap[enemyName]
	if !ok {
		return RenderSprite(defaultEnemySprite)
	}
	return RenderSprite(frames)
}

// ============================================================
// INIMIGOS — Grid 11x10 (renderiza como 11x5 chars)
// ============================================================

var defaultEnemySprite = spriteFrame{
	rows: []string{
		"...KKKKK...",
		"..KDDDDDK..",
		".KDDDDDDDK.",
		".KDWDKWDDDK",
		".KDDDDDDDK.",
		".KKDDDDDKK.",
		"..KDDDDDK..",
		"..KKDDDKK..",
		"...KK.KK...",
		"...........",
	},
}

// SLIME — bolha gelatinosa verde translúcida, olhos simples
var slimeSprite = spriteFrame{
	rows: []string{
		"...........",
		"...LLLLL...",
		"..LGGGGGL..",
		".LGWGKWGGL.",
		".LGGGGGGL..",
		".LGGGGGL...",
		"..LGGGGL...",
		"..LLLLLL...",
		"...........",
		"...........",
	},
}

// RATO GIGANTE — corpo marrom, orelhas grandes, cauda conectada
var ratoSprite = spriteFrame{
	rows: []string{
		".M.......M.",
		"KMK.KKK.KMK",
		".KMMMMMMMK.",
		".KMWMKMWMK.",
		".KMMMMMMK..",
		".KKMPMMMKK.",
		"..KMMMMMKM.",
		"..KKMMKKKMK",
		"...K.K..KM.",
		"...........",
	},
}

// ESQUELETO — ossos brancos, olhos vazios (vermelhos)
var esqueletoSprite = spriteFrame{
	rows: []string{
		"...KKKKK...",
		"..KWWWWWK..",
		".KWWWWWWWK.",
		".KWRKWRKWK.",
		".KWWWWWWWK.",
		".KKWKWKWKK.",
		"...KWWWK...",
		"..KKWKWKK..",
		"..KW.K.WK..",
		"...........",
	},
}

// MORCEGO VAMPIRO — asas abertas, olhos vermelhos, presas (11 cols)
var morcegoSprite = spriteFrame{
	rows: []string{
		"K.........K",
		"KK.KKKKK.KK",
		"KKDDDDDDDKK",
		"KDDRWKWRDDK",
		".KDDDDDDDK.",
		".KKDWKWDKK.",
		"..KDDDDDK..",
		"..KKDDDKK..",
		"...........",
		"...........",
	},
}

// GOBLIN GUERREIRO — pele verde, elmo marrom visível, cara agressiva
var goblinSprite = spriteFrame{
	rows: []string{
		"..KKKKKKK..",
		"..KMMMMMK..",
		".KKGGGGGKK.",
		".KGRRKRRGK.",
		".KGGGGGGGK.",
		".KKGYKYGKK.",
		"..KGGGGGK..",
		".KKGGGGGKK.",
		"..KK.K.KK..",
		"...........",
	},
}

// ARANHA VENENOSA — corpo escuro, 8 patas, olhos múltiplos vermelhos
var aranhaSprite = spriteFrame{
	rows: []string{
		"D..KKKKK..D",
		".DKDDDDDK..",
		"DKDRRRRRDKD",
		".KDDDDDDDK.",
		"DKKDDDDDKKD",
		"..KDDDDDK..",
		".D.KDKDK.D.",
		"D..K.K.K..D",
		"...........",
		"...........",
	},
}

// CAVALEIRO NEGRO — armadura escura, elmo com visor vermelho (11 cols fixo)
var cavaleiroSprite = spriteFrame{
	rows: []string{
		"...KDKDK...",
		"..KKDDDKK..",
		".KDDDDDDDK.",
		".KDDRRDDDK.",
		".KDDDDDDDK.",
		"KKDDDDDDDKK",
		"KKDDDDDDDKK",
		".KKDDKDDKK.",
		".KKK.K.KKK.",
		"...........",
	},
}

// MAGO SOMBRIO — manto roxo, cajado, olhos brilhantes
var magoSprite = spriteFrame{
	rows: []string{
		"....KYK....",
		"...KKKKK...",
		"..KBBBBBK..",
		".KBBYYBBBK.",
		".KBBBBBBBK.",
		".KKBBBBBKK.",
		"YKBBBBBBBKY",
		".KKBBBBBKK.",
		"..KK.K.KK..",
		"...........",
	},
}

// DRAGAO ANCIAO (BOSS) — Grid maior: 15x12 (renderiza 15x6)
// Corpo vermelho/laranja, asas, chamas, presença imponente
var dragaoSprite = spriteFrame{
	rows: []string{
		"..K...........K..",
		".KRK...KKK...KRK",
		"..KKKKRRRRRKKK..",
		".KRRRRRRRRRRRRK.",
		".KRRYWKRYRWYKRK.",
		".KRRRRRRRRRRRRK.",
		"KKRRRRRRRRRRRKK.",
		"YKRRROORROORRKOY",
		".KKRRRRRRRRRKK..",
		"..KKRRRKRRRK...",
		"...KKK.K.KKK...",
		"................",
	},
}

var enemySpriteMap = map[string]spriteFrame{
	"Slime":            slimeSprite,
	"Rato Gigante":     ratoSprite,
	"Esqueleto":        esqueletoSprite,
	"Morcego Vampiro":  morcegoSprite,
	"Goblin Guerreiro": goblinSprite,
	"Aranha Venenosa":  aranhaSprite,
	"Cavaleiro Negro":  cavaleiroSprite,
	"Mago Sombrio":     magoSprite,
	"Dragao Anciao":    dragaoSprite,
}

// ============================================================
// BIOME BACKGROUNDS — pinturas de fundo estilo GBA (Pokemon Fire Red)
// Grid: 35x12 (renderiza 35x6) — panorama largo
// ============================================================

// BiomeSprite retorna a pintura de fundo pixel art do bioma.
func BiomeSprite(biomeName string) string {
	frame, ok := biomeSpriteMap[biomeName]
	if !ok {
		return ""
	}
	return RenderSprite(frame)
}

var biomeSpriteMap = map[string]spriteFrame{
	// FLORESTAL — cena de floresta: céu, copas, troncos, chão com flores
	// 25 cols x 10 rows (renderiza 25x5)
	"Florestal": {rows: []string{
		"....CCCCCCCCCCCCCCCCC....",
		"..CCCCCCCCCCCCCCCCCCCC..",
		".GGLGGGGGLGGGGGLGGGGGLG.",
		"GGGLGGGGLGGGGGLGGGGLGGG",
		"GGGGGGLGGGGGGGGGLGGGGGG",
		"GGGLGGGGGGLGGGGGGLGGGGG",
		"..MGLM..MGLM..MGLM..MG",
		"LLMLLMLLMLLMLLMLLMLLMLL",
		"LLLPLLLLLYLLLPLLLLLYLLL",
		"LLLLLLLLLLLLLLLLLLLLLLLLL",
	}},

	// GELIDO — cena de tundra: céu azul, montanhas de gelo contrastadas, cristais
	// 25 cols x 10 rows
	"Gelido": {rows: []string{
		"BBBBBBBBBBBBBBBBBBBBBBBBB",
		"BBBBBBBCBBBBBBCBBBBBBBCBB",
		"..CK..CKWC..CK..CKWC..C",
		".CKWC.CKWWC.CKWC.CKWWC.",
		"CKWWWCCKWWWCCKWWWCCKWWWC",
		"KWWWCKWWWWCKWWWWCKWWWWCKW",
		"WWWWWWWWWWWWWWWWWWWWWWWWW",
		"WCWWWWWCWWWWWCWWWWWCWWWWW",
		"WWWWCWWWWWCWWWWWCWWWWWCWW",
		"WWWWWWWWWWWWWWWWWWWWWWWWW",
	}},

	// VULCANICO — cena vulcânica: céu vermelho, vulcão, lava, rochas
	// 25 cols x 10 rows
	"Vulcanico": {rows: []string{
		"KKKKKKKKYRKKKKKKKKYRKKKKK",
		"KKKKKKKYOYRKKKKKKYOYRKKK",
		"......DDODDD......DDODDD",
		".....DDDRODDD....DDDRODD",
		"....DDDDDRDDDDD.DDDDDRDD",
		"RRDDDDDDDDDDDRRDDDDDDDD",
		"ORRDDDDDDDDDORRDDDDDDDD",
		"RORRRDDDDRORRRDDDDRORRD",
		"RRRORRRRRRRRORRRRRRRRORRR",
		"ORRRRORRRRORRRRORRRRORRRO",
	}},

	// ABISSAL — cena abissal: profundezas com muita bioluminescência
	// 25 cols x 10 rows
	"Abissal": {rows: []string{
		"KKKKKKKKKKKKKKKKKKKKKKKK",
		"KKBKKCKKKBKKKKKCKKKBKKKK",
		"KKKKBKKKCKKBKKKKKBKCKKBK",
		"DKDDDKDDDDKDDDDKDDDDDKD",
		"DDCDDDDDDCDDDDDDCDDDDDC",
		"DDDDBDDDDDDDCDDDDDBDDDD",
		"DCDDDDCDDDBDDDDDCDDDDBD",
		"DDDCDDDDDDDDCDDDDDDCDDD",
		"DBDDDDBDDDDDDDDBDDDDDDDB",
		"KKKCKKKKKKKCKKKKKKKCKKKK",
	}},
}

// ============================================================
// EFEITOS VISUAIS DE COMBATE — estilo GBA
// ============================================================

// CombatEffect retorna uma animação visual para ações de combate.
func CombatEffect(effectType string) string {
	frame, ok := combatEffectMap[effectType]
	if !ok {
		return ""
	}
	return RenderSprite(frame)
}

var combatEffectMap = map[string]spriteFrame{
	// ATAQUE — corte rápido (7x4)
	"attack": {rows: []string{
		"......W",
		"....WK.",
		"..WK...",
		"WK.....",
	}},

	// ATAQUE CRITICO — explosão (7x6)
	"critical": {rows: []string{
		"..YOY..",
		".YOROY.",
		"YORRROY",
		"YORRROY",
		".YOROY.",
		"..YOY..",
	}},

	// DEFESA — escudo (7x6)
	"defend": {rows: []string{
		".CCCCC.",
		"CBBBBBC",
		"CBBWBBC",
		"CBBBBBC",
		".CBBBC.",
		"..CBC..",
	}},

	// CURA — brilho verde (7x4)
	"heal": {rows: []string{
		"..LGL..",
		".LGGGL.",
		"..LGL..",
		".L...L.",
	}},

	// SKILL — energia (7x4)
	"skill": {rows: []string{
		".BCCCB.",
		"BCCYCCCB",
		"BCCYCCCB",
		".BCCCB.",
	}},

	// FUGA — poeira (7x4)
	"flee": {rows: []string{
		"D......",
		".DD....",
		"..DDD..",
		"...DD..",
	}},

	// DANO — flash (7x4)
	"damage": {rows: []string{
		"R....R.",
		".R..R..",
		".R..R..",
		"R....R.",
	}},

	// MISS — vento (7x4)
	"miss": {rows: []string{
		"DDD....",
		"..DDD..",
		"....DDD",
		"..DDD..",
	}},
}

// ============================================================
// SPRITES DE UI — baú, poção, fogueira, etc.
// Grid: 9x8 (renderiza 9x4)
// ============================================================

// ItemSprite retorna o sprite de um item/elemento de UI.
func ItemSprite(itemType string) string {
	frame, ok := itemSpriteMap[itemType]
	if !ok {
		return ""
	}
	return RenderSprite(frame)
}

var itemSpriteMap = map[string]spriteFrame{
	// BAÚ DE TESOURO — madeira com detalhes dourados
	"chest": {rows: []string{
		".KKKKKKK.",
		"KYYYYYYY.",
		"KYYKYKKY.",
		"KYYYYYYY.",
		"KKKKKKKKK",
		"KMMYMMMK.",
		"KMMMMMMK.",
		"KKKKKKKKK",
	}},

	// POÇÃO DE CURA
	"potion": {rows: []string{
		"...KKK...",
		"...KWK...",
		"..KKKKK..",
		".KKRRRK..",
		".KRRRRRK.",
		".KRRLRRK.",
		".KKRRRK..",
		"..KKKKK..",
	}},

	// DESCANSO (fogueira com pedras)
	"rest": {rows: []string{
		"...YOY...",
		"..YOYOY..",
		"..OYOYO..",
		".YOYOYOY.",
		"..ROROR..",
		".DDDRDDDD",
		"DDDDDDDDD",
		"DDDDDDDDD",
	}},

	// VITÓRIA (troféu com brilho)
	"victory": {rows: []string{
		"Y.KKKKK.Y",
		".YKYYYKKY",
		".KYYYYYK.",
		"..KYYYKY.",
		"...KYK...",
		"..KKKKK..",
		".KKKKKKK.",
		"KKKKKKKKK",
	}},

	// DERROTA (caveira)
	"defeat": {rows: []string{
		"..KKKKK..",
		".KWWWWWK.",
		"KWWWWWWWK",
		"KWRKWRKWK",
		"KWWWWWWWK",
		".KKWKWKK.",
		"..KWKWK..",
		"...KKK...",
	}},

	// LOJA (barraca com toldo)
	"shop": {rows: []string{
		"KKKKKKKKK",
		"KYRYRYRKYY",
		"KYYYYYYY.",
		"KKKKKKKKK",
		"KM.WWW.MK",
		"KM.WYW.MK",
		"KM.WWW.MK",
		"KKKKKKKKK",
	}},
}
