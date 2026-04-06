package dungeon

import (
	"fmt"
	"strings"
)

// EnemyArt retorna a ASCII art de um inimigo pelo ID.
func EnemyArt(enemyName string) string {
	art, ok := enemyArtMap[enemyName]
	if !ok {
		return defaultArt
	}
	return art
}

var defaultArt = `    .---.
   /     \
  |  ? ?  |
   \ --- /
    '---'`

var enemyArtMap = map[string]string{
	"Slime": `    .-"""-.
   /       \
  |  o   o  |
  |    ~    |
   \       /
    '-----'`,

	"Rato Gigante": `       /\  /\
      {  \/ }
      {  () }
   __ { /\ } __
  /  '-\  /-'  \
  \____/\/\____/
     /||    ||\
      ""    ""`,

	"Esqueleto": `     .-.
    (o.o)
   __|=|__
  //.=|=.\\
 // .=|=. \\
 \\ .=|=. //
  \\(_=_)//
   (:| |:)
    || ||
    () ()`,

	"Morcego Vampiro": `   /\  _  /\
  / \\/ \// \
  \  (oo)  /
   '--  --'
   /|    |\
  / |    | \`,

	"Goblin Guerreiro": `    .-^^^-.
   /  o o  \
  |  { ^ }  |
  | \  =  / |
  /\|'---'|/\
 / /|     |\ \
   ||     ||
   ""     ""`,

	"Aranha Venenosa": `  \  \._./ /
   \ (o o)/
  --( /V\ )--
   / (   ) \
  /  /'-'\  \
     /   \`,

	"Cavaleiro Negro": `     ,  ,
    |\ /|
    ( O O )
     \_=_/
   __|   |__
  [_________]
   /||   ||\
  / ||   || \
    ||   ||
   _||   ||_`,

	"Mago Sombrio": `     /\
    /  \
   / ** \
  |  ()  |
  | /--\ |
  |/ /\ \|
   |/  \|
   ||  ||
   ""  ""`,

	"Dragao Anciao": `          /\    .-" /
         /  ; .'  .'
        :   :/  .'
         \  ;-.'
    .--""""/   \""""--.
   /  .-'"/.';. \""-.  \
  :  :  //    \\  :  :
  '._\ //  /\  \\ /_.'
      \_/  :  :  \_/
           '..'`,
}

// HPBar gera uma barra de HP visual para o combate.
func HPBar(current, max, width int) string {
	if max <= 0 {
		max = 1
	}
	pct := float64(current) / float64(max)
	if pct < 0 {
		pct = 0
	}
	if pct > 1 {
		pct = 1
	}
	filled := int(pct * float64(width))
	bar := ""
	for i := 0; i < width; i++ {
		if i < filled {
			bar += "▇"
		} else {
			bar += "░"
		}
	}
	return bar
}

// MPBar gera uma barra de MP visual para o combate.
func MPBar(current, max, width int) string {
	if max <= 0 {
		max = 1
	}
	pct := float64(current) / float64(max)
	if pct < 0 {
		pct = 0
	}
	if pct > 1 {
		pct = 1
	}
	filled := int(pct * float64(width))
	bar := ""
	for i := 0; i < width; i++ {
		if i < filled {
			bar += "▣"
		} else {
			bar += "▢"
		}
	}
	return bar
}

// BiomeArt retorna a ASCII art de introdução para um bioma
func BiomeArt(biome Biome) string {
	switch biome {
	case BiomeForest:
		return forestArt
	case BiomeIcy:
		return icyArt
	case BiomeVolcanic:
		return volcanicArt
	case BiomeAbyssal:
		return abyssalArt
	default:
		return defaultBiomeArt
	}
}

// Artes dos biomas — estilo GBA/retro sem emojis
var forestArt = `
    ,,  ,,    /\    ,,  ,,    /\    ,,
   /||\//\\  /||\  /||\//\\  /||\  /||\
  / || \  \\/ || \/ || \  \\/ || \/ || \
  |_||_|  ||  ||_|  ||_|  ||  ||_|  ||_|
  ,.,.,.,.,.,.,.,.,.,.,.,.,.,.,.,.,.,.,.,
  :.:.:.:.:.:.:.:.:.:.:.:.:.:.:.:.:.:.:.:.

    +---------------------------------+
    |       FLORESTA  ENCANTADA       |
    +---------------------------------+

    Cogumelos brilham no chao.
    O ar e puro e fresco...`

var icyArt = `
     *    .  *  .    *  .  *   .  *
    /\   *  /\  . * /\  *  /\   *  /\
   /  \ . /  \  * /  \ . /  \  * /  \
  /    \/    \  /    \/    \  /    \/
  ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
  ......................................

    +---------------------------------+
    |       TUNDRA  CONGELANTE        |
    +---------------------------------+

    O vento gelado corta a pele.
    Seu halito congela no ar...`

var volcanicArt = `
        /\          /\
       /##\    .   /##\    .
      /####\ /#\ /####\ /#\
     /######\####\######\####\
    ~^~^~^~^~^~^~^~^~^~^~^~^~^~^~
    ~.~.~.~.~.~.~.~.~.~.~.~.~.~.~

    +---------------------------------+
    |      MONTANHAS  DE  FOGO        |
    +---------------------------------+

    Lava borbulha ao redor...
    O calor e sufocante!`

var abyssalArt = `
    ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
    ~  .    .        .    .        .    .  ~
    ~     .    .  .     .    .  .     .    ~
    ~  .     .       .     .       .     .~
    ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
    ........................................

    +---------------------------------+
    |       ABISSO  MISTERIOSO        |
    +---------------------------------+

    Escuridao infinita...
    Algo observa das sombras...`

var defaultBiomeArt = `
     _____     _____     _____     _____
    |     |   |     |   |     |   |     |
    | [ ] |   | [ ] |   | [ ] |   | [ ] |
    |_____|   |_____|   |_____|   |_____|
    ========================================
    ........................................

    +---------------------------------+
    |        MASMORRA  SOMBRIA        |
    +---------------------------------+

    Pedras antigas te observam.
    O perigo espreita...`

// GetBiomeIntroMessage retorna a string formatada com a arte do bioma e todos os efeitos
func GetBiomeIntroMessage(biome Biome) string {
	art := BiomeArt(biome)
	modifier := ModifierForBiome(biome)

	var sb strings.Builder
	sb.WriteString(art)
	sb.WriteString("\n\n")
	sb.WriteString(strings.Repeat("-", 40))
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("  BIOMA: %s\n", biome.String()))
	sb.WriteString(fmt.Sprintf("  Efeitos: %s\n", FormatBiomeEffects(modifier)))
	sb.WriteString(strings.Repeat("-", 40))
	sb.WriteString("\n")
	sb.WriteString("\n  Pressione ENTER para entrar...")

	return sb.String()
}

// FormatBiomeEffects formata os efeitos do bioma em uma string legível.
func FormatBiomeEffects(m BiomeModifier) string {
	var effects []string
	if m.HungerRecoveryBonus > 0 {
		effects = append(effects, fmt.Sprintf("+%d recuperação fome", m.HungerRecoveryBonus))
	}
	if m.FriendlyEncounterChance > 0 {
		effects = append(effects, fmt.Sprintf("+%d%% encontros amistosos", m.FriendlyEncounterChance))
	}
	if m.PlayerFireVulnerability > 0 {
		effects = append(effects, fmt.Sprintf("+%d%% vulnerabilidade fogo", m.PlayerFireVulnerability))
	}
	if m.PlayerSpeedPenaltyPct > 0 {
		effects = append(effects, fmt.Sprintf("-%d%% velocidade", m.PlayerSpeedPenaltyPct))
	}
	if m.PlayerIceDefensePct > 0 {
		effects = append(effects, fmt.Sprintf("+%d%% defesa gelo", m.PlayerIceDefensePct))
	}
	if m.FreezeChancePct > 0 {
		effects = append(effects, fmt.Sprintf("%d%% chance congelamento", m.FreezeChancePct))
	}
	if m.EnemyAttackBonusPct > 0 {
		effects = append(effects, fmt.Sprintf("+%d%% ATK inimigo", m.EnemyAttackBonusPct))
	}
	if m.PlayerFireDotPctMaxHP > 0 {
		effects = append(effects, fmt.Sprintf("DOT fogo %d%% HP", m.PlayerFireDotPctMaxHP))
	}
	if m.HappinessPenaltyTick > 0 {
		effects = append(effects, fmt.Sprintf("-%d felicidade/tick", m.HappinessPenaltyTick))
	}
	if m.EnemyLuckBonusPct > 0 {
		effects = append(effects, fmt.Sprintf("+%d%% sorte inimiga", m.EnemyLuckBonusPct))
	}
	if m.RareLootChanceBonusPct > 0 {
		effects = append(effects, fmt.Sprintf("+%d%% loot raro", m.RareLootChanceBonusPct))
	}
	if m.PlayerHPRegenPenaltyPct > 0 {
		effects = append(effects, fmt.Sprintf("-%d%% regen HP", m.PlayerHPRegenPenaltyPct))
	}

	if len(effects) == 0 {
		return "Sem efeitos especiais"
	}
	return strings.Join(effects, " | ")
}
