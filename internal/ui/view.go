package ui

// --- Constantes de tamanho mínimo ---
const (
	minWidth  = 60
	minHeight = 18
)

// layout contém todas as dimensões computadas para o render responsivo.
type layout struct {
	contentWidth  int
	avatarWidth   int
	avatarHeight  int
	statsWidth    int
	statsHeight   int
	barWidth      int
	inputWidth    int
	gameOverWidth int
	gameBoxWidth  int
	dungeonWidth  int
	dungeonHeight int
	tooSmall      bool
}

func computeLayout(termWidth, termHeight int) layout {
	if termWidth < minWidth || termHeight < minHeight {
		return layout{tooSmall: true}
	}

	contentWidth := clampInt(termWidth-4, 0, 200)
	avatarWidth := clampInt(contentWidth*35/100, 26, 40)
	statsWidth := clampInt(contentWidth-avatarWidth-4, 30, 200)
	boxHeight := clampInt(termHeight-12, 10, 18)
	barWidth := clampInt(statsWidth-22, 10, 40)
	inputWidth := contentWidth
	gameOverWidth := clampInt(contentWidth*60/100, 40, 60)
	gameBoxWidth := clampInt(contentWidth-4, 40, 200)
	dungeonWidth := clampInt(contentWidth-2, 36, 80)
	dungeonHeight := clampInt(termHeight-8, 20, 40)

	return layout{
		contentWidth:  contentWidth,
		avatarWidth:   avatarWidth,
		avatarHeight:  boxHeight,
		statsWidth:    statsWidth,
		statsHeight:   boxHeight,
		barWidth:      barWidth,
		inputWidth:    inputWidth,
		gameOverWidth: gameOverWidth,
		gameBoxWidth:  gameBoxWidth,
		dungeonWidth:  dungeonWidth,
		dungeonHeight: dungeonHeight,
		tooSmall:      false,
	}
}

func clampInt(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
