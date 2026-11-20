// Package screenpack loads and validates an OpenKakutou screenpack: the
// format of docs/format.md, checked against the schemas of package schema
// and across files.
package screenpack

// Path is an asset path, relative to the screenpack folder (or, in a sprite
// sheet, to the sheet's folder).
type Path string

// StyleID names a text style of the screen or of the root file.
type StyleID string

// FontID names a font of the root file.
type FontID string

// Screenpack is a loaded screenpack.
type Screenpack struct {
	Name       string
	Author     string
	Version    string
	Fonts      map[FontID]Path
	TextStyles map[StyleID]TextStyle
	Sounds     Sounds
	Music      *Music
	// WarningColor tints what the game flags as a problem; empty means the
	// format's default, #dc3c3c.
	WarningColor string
	// Screens holds the screens the pack has; a nil one falls back to the
	// default screenpack's.
	Screens Screens
	// Sheets are the sprite sheets the screens use, by path.
	Sheets map[Path]*SpriteSheet
	// Warnings are problems that do not stop loading, such as an unknown
	// file in screens/.
	Warnings Problems
}

type rootFile struct {
	Name       string                `json:"name"`
	Author     string                `json:"author"`
	Version    string                `json:"version"`
	Fonts      map[FontID]Path       `json:"fonts"`
	TextStyles map[StyleID]TextStyle `json:"textStyles"`
	Sounds     Sounds                `json:"sounds"`
	Music      *Music                `json:"music"`
	// WarningColor is the root file's "warningColor".
	WarningColor string `json:"warningColor"`
}

// Sounds maps a sound name (move, confirm, cancel, error, and a screen's
// extra ones) to its file.
type Sounds map[string]Path

// Music is a looping music track.
type Music struct {
	File        Path     `json:"file"`
	Volume      *float64 `json:"volume"`
	LoopStartMs int      `json:"loopStartMs"`
}

// TextStyle is how a text is drawn.
type TextStyle struct {
	Font      FontID   `json:"font"`
	Size      float64  `json:"size"`
	Color     string   `json:"color"`
	Align     string   `json:"align"`
	Outline   *Outline `json:"outline"`
	Shadow    *Shadow  `json:"shadow"`
	Uppercase bool     `json:"uppercase"`
}

// Outline is a text style's outline.
type Outline struct {
	Color string  `json:"color"`
	Width float64 `json:"width"`
}

// Shadow is a text style's drop shadow.
type Shadow struct {
	Color  string     `json:"color"`
	Offset [2]float64 `json:"offset"`
}

// Layout places an element or a part: Pivot is set on Anchor, moved by
// Offset, in units of 1/100 of the screen height. A nil Size side keeps the
// image's aspect ratio.
type Layout struct {
	Anchor string       `json:"anchor"`
	Pivot  string       `json:"pivot"`
	Offset [2]float64   `json:"offset"`
	Size   *[2]*float64 `json:"size"`
}

// Effect is one half of a transition.
type Effect struct {
	Kind    string `json:"kind"`
	Ms      int    `json:"ms"`
	DelayMs int    `json:"delayMs"`
	Easing  string `json:"easing"`
	From    any    `json:"from"`
}

// Transition is how a screen or an element enters and leaves.
type Transition struct {
	In  *Effect `json:"in"`
	Out *Effect `json:"out"`
}

// Element is one drawn thing: an image, a text or a rectangle (Type).
type Element struct {
	Layout
	Type       string      `json:"type"`
	Visible    *bool       `json:"visible"`
	Opacity    *float64    `json:"opacity"`
	Transition *Transition `json:"transition"`
	// Image is a still image; Sheet with Animation or Frame shows a sprite
	// sheet instead.
	Image     Path        `json:"image"`
	Sheet     Path        `json:"sheet"`
	Animation string      `json:"animation"`
	Frame     string      `json:"frame"`
	Tile      string      `json:"tile"`
	Scroll    *[2]float64 `json:"scroll"`
	Text      string      `json:"text"`
	Style     StyleID     `json:"style"`
	Color     string      `json:"color"`
	Radius    float64     `json:"radius"`
}

// TextPart is a part showing text the game supplies.
type TextPart struct {
	Layout
	Style StyleID `json:"style"`
}

// Menu styles and arranges the entries the game offers.
type Menu struct {
	Layout
	Direction      string      `json:"direction"`
	Spacing        float64     `json:"spacing"`
	VisibleEntries int         `json:"visibleEntries"`
	Style          StyleID     `json:"style"`
	ActiveStyle    StyleID     `json:"activeStyle"`
	DisabledStyle  StyleID     `json:"disabledStyle"`
	FlaggedStyle   StyleID     `json:"flaggedStyle"`
	Cursor         *Element    `json:"cursor"`
	Order          []string    `json:"order"`
	Hidden         []string    `json:"hidden"`
	Sounds         Sounds      `json:"sounds"`
	ValueStyle     StyleID     `json:"valueStyle"`
	ValueOffset    *[2]float64 `json:"valueOffset"`
	// ValueOffsets places each value column of an entry with several values.
	ValueOffsets [][2]float64 `json:"valueOffsets"`
}

// Sort is a select screen's order rule.
type Sort struct {
	By         string `json:"by"`
	Descending bool   `json:"descending"`
}

// Sides holds one part per player side.
type Sides[T any] struct {
	P1 *T `json:"p1"`
	P2 *T `json:"p2"`
}

// Portraits places both sides' large portraits.
type Portraits struct {
	Sides[Layout]
	MirrorP2 *bool `json:"mirrorP2"`
}

// Screen holds the fields every screen file has. Silent is set when the
// file asks for no music ("music": null), as opposed to leaving it out.
type Screen struct {
	Background []Element             `json:"background"`
	Foreground []Element             `json:"foreground"`
	Music      *Music                `json:"music"`
	Silent     bool                  `json:"-"`
	Sounds     Sounds                `json:"sounds"`
	TextStyles map[StyleID]TextStyle `json:"textStyles"`
	Transition *Transition           `json:"transition"`
}

// Title is the title and main menu screen.
type Title struct {
	Screen
	Logo       *Element    `json:"logo"`
	Menu       *Menu       `json:"menu"`
	PressStart *PressStart `json:"pressStart"`
	Version    *TextPart   `json:"version"`
}

// PressStart is the title's "press start" prompt.
type PressStart struct {
	TextPart
	Text    string `json:"text"`
	BlinkMs *int   `json:"blinkMs"`
}

// CharacterSelect is the character select screen.
type CharacterSelect struct {
	Screen
	Grid     *Grid            `json:"grid"`
	Sort     *Sort            `json:"sort"`
	Cursor   *Cursors         `json:"cursor"`
	Portrait *Portraits       `json:"portrait"`
	Name     *Sides[TextPart] `json:"name"`
	Palette  *Sides[TextPart] `json:"palette"`
	Title    *TextPart        `json:"title"`
}

// Grid is the character select's grid of cells.
type Grid struct {
	Layout
	Columns   int         `json:"columns"`
	Rows      int         `json:"rows"`
	CellSize  *[2]float64 `json:"cellSize"`
	Spacing   *[2]float64 `json:"spacing"`
	Cell      *Element    `json:"cell"`
	EmptyCell *Element    `json:"emptyCell"`
	Random    *Element    `json:"random"`
}

// Cursors are the select cursors of each side, and the one shown when both
// share a cell.
type Cursors struct {
	P1   *Element `json:"p1"`
	P2   *Element `json:"p2"`
	Both *Element `json:"both"`
}

// StageSelect is the stage select screen.
type StageSelect struct {
	Screen
	Sort    *Sort     `json:"sort"`
	Preview *Layout   `json:"preview"`
	Name    *TextPart `json:"name"`
	Title   *TextPart `json:"title"`
	Arrows  *Arrows   `json:"arrows"`
}

// Arrows show that more items lie on a side.
type Arrows struct {
	Left  *Element `json:"left"`
	Right *Element `json:"right"`
}

// MatchSetup is the match setup screen.
type MatchSetup struct {
	Screen
	Menu  *Menu     `json:"menu"`
	Title *TextPart `json:"title"`
}

// Versus is the VS screen.
type Versus struct {
	Screen
	Portrait   *Portraits       `json:"portrait"`
	Name       *Sides[TextPart] `json:"name"`
	Stage      *TextPart        `json:"stage"`
	DurationMs *int             `json:"durationMs"`
}

// Options is the options screen.
type Options struct {
	Screen
	Tabs     *Menu         `json:"tabs"`
	Menu     *Menu         `json:"menu"`
	Controls *ControlSlots `json:"controls"`
	Hint     *TextPart     `json:"hint"`
	Title    *TextPart     `json:"title"`
}

// ControlSlots style the controls page's bindings.
type ControlSlots struct {
	TextPart
	SlotStyle    StyleID `json:"slotStyle"`
	WaitingStyle StyleID `json:"waitingStyle"`
}

// Pause is the pause screen.
type Pause struct {
	Screen
	Dim   *Dim      `json:"dim"`
	Title *TextPart `json:"title"`
	Menu  *Menu     `json:"menu"`
}

// Dim is drawn over the frozen match.
type Dim struct {
	Color string `json:"color"`
}

// MoveList is the move list screen.
type MoveList struct {
	Screen
	Title     *TextPart    `json:"title"`
	List      *MoveLines   `json:"list"`
	Buttons   *ButtonIcons `json:"buttons"`
	Scrollbar *Scrollbar   `json:"scrollbar"`
}

// MoveLines is the move list itself.
type MoveLines struct {
	TextPart
	VisibleLines int     `json:"visibleLines"`
	LineSpacing  float64 `json:"lineSpacing"`
}

// ButtonIcons map button ids to frames of one sprite sheet.
type ButtonIcons struct {
	Sheet  Path              `json:"sheet"`
	Frames map[string]string `json:"frames"`
}

// Scrollbar is shown when a list is longer than its visible lines.
type Scrollbar struct {
	Track *Element `json:"track"`
	Thumb *Element `json:"thumb"`
}

// Results is the results screen.
type Results struct {
	Screen
	Winner  *Winner   `json:"winner"`
	Message *TextPart `json:"message"`
	Score   *TextPart `json:"score"`
	Menu    *Menu     `json:"menu"`
}

// Winner shows the match winner.
type Winner struct {
	Portrait *Layout   `json:"portrait"`
	Name     *TextPart `json:"name"`
}

// Continue is the continue screen.
type Continue struct {
	Screen
	Portrait  *Layout   `json:"portrait"`
	Countdown *TextPart `json:"countdown"`
	Message   *TextPart `json:"message"`
	Menu      *Menu     `json:"menu"`
}

// Victory is the ending screen.
type Victory struct {
	Screen
	Portrait   *Layout   `json:"portrait"`
	Name       *TextPart `json:"name"`
	Message    *TextPart `json:"message"`
	DurationMs *int      `json:"durationMs"`
}

// Error is the screen showing a problem the game reports.
type Error struct {
	Screen
	Title   *TextPart     `json:"title"`
	Message *ErrorMessage `json:"message"`
	Hint    *TextPart     `json:"hint"`
}

// ErrorMessage is the error text, wrapped inside Width units.
type ErrorMessage struct {
	TextPart
	Width float64 `json:"width"`
}

// Screens are a screenpack's screens; nil means the pack does not have it.
type Screens struct {
	Title           *Title
	CharacterSelect *CharacterSelect
	StageSelect     *StageSelect
	MatchSetup      *MatchSetup
	Versus          *Versus
	Options         *Options
	Pause           *Pause
	MoveList        *MoveList
	Results         *Results
	Continue        *Continue
	Victory         *Victory
	Error           *Error
}

// slot returns a new screen value for a screen name and stores it in s, or
// nil for an unknown name.
func (s *Screens) slot(name string) any {
	switch name {
	case "title":
		s.Title = &Title{}
		return s.Title
	case "characterSelect":
		s.CharacterSelect = &CharacterSelect{}
		return s.CharacterSelect
	case "stageSelect":
		s.StageSelect = &StageSelect{}
		return s.StageSelect
	case "matchSetup":
		s.MatchSetup = &MatchSetup{}
		return s.MatchSetup
	case "versus":
		s.Versus = &Versus{}
		return s.Versus
	case "options":
		s.Options = &Options{}
		return s.Options
	case "pause":
		s.Pause = &Pause{}
		return s.Pause
	case "moveList":
		s.MoveList = &MoveList{}
		return s.MoveList
	case "results":
		s.Results = &Results{}
		return s.Results
	case "continue":
		s.Continue = &Continue{}
		return s.Continue
	case "victory":
		s.Victory = &Victory{}
		return s.Victory
	case "error":
		s.Error = &Error{}
		return s.Error
	}
	return nil
}

// SpriteSheet is an image cut into frames, and animations over them.
type SpriteSheet struct {
	// Image is relative to the sheet's own folder.
	Image      Path                 `json:"image"`
	Frames     map[string][4]int    `json:"frames"`
	Animations map[string]Animation `json:"animations"`
}

// Animation plays sheet frames in order.
type Animation struct {
	Frames []AnimationFrame `json:"frames"`
	// Loop is true when unset.
	Loop *bool `json:"loop"`
}

// AnimationFrame is one step of an animation.
type AnimationFrame struct {
	Frame string `json:"frame"`
	Ms    int    `json:"ms"`
}
