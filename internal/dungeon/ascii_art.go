package dungeon

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
