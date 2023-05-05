package todo

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// FilterState describes the current filtering state on the model.
type FilterState int

// Possible filter states.
const (
	Unfiltered    FilterState = iota // no filter set
	Filtering                        // user is actively setting a filter
	FilterApplied                    // a filter is applied and user is not editing filter
)

const (
	bullet   = "•"
	ellipsis = "…"
)

type Paginator struct {
	n_pages      int
	current_page int
	pages        [][]Memo
}
type Keys struct {
	// Keybindings used when browsing the list.
	CursorUp    key.Binding
	CursorDown  key.Binding
	NextPage    key.Binding
	PrevPage    key.Binding
	GoToStart   key.Binding
	GoToEnd     key.Binding
	Filter      key.Binding
	ClearFilter key.Binding

	// Keybindings used when setting a filter.
	CancelWhileFiltering key.Binding
	AcceptWhileFiltering key.Binding

	// Help toggle keybindings.
	ShowFullHelp  key.Binding
	CloseFullHelp key.Binding

	// The quit keybinding. This won't be caught when filtering.
	Quit key.Binding

	// The quit-no-matter-what keybinding. This will be caught when filtering.
	ForceQuit key.Binding
}

// DefaultKeyMap returns a default set of keybindings.
func DefaultKeys() Keys {
	return Keys{
		// Browsing.
		CursorUp: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up"),
		),
		CursorDown: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "down"),
		),
		PrevPage: key.NewBinding(
			key.WithKeys("left", "h", "pgup", "b", "u"),
			key.WithHelp("←/h/pgup", "prev page"),
		),
		NextPage: key.NewBinding(
			key.WithKeys("right", "l", "pgdown", "f", "d"),
			key.WithHelp("→/l/pgdn", "next page"),
		),
		GoToStart: key.NewBinding(
			key.WithKeys("home", "g"),
			key.WithHelp("g/home", "go to start"),
		),
		GoToEnd: key.NewBinding(
			key.WithKeys("end", "G"),
			key.WithHelp("G/end", "go to end"),
		),
		Filter: key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "filter"),
		),
		ClearFilter: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "clear filter"),
		),

		// Filtering.
		CancelWhileFiltering: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "cancel"),
		),
		AcceptWhileFiltering: key.NewBinding(
			key.WithKeys("enter", "tab", "shift+tab", "ctrl+k", "up", "ctrl+j", "down"),
			key.WithHelp("enter", "apply filter"),
		),

		// Toggle help.
		ShowFullHelp: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "more"),
		),
		CloseFullHelp: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "close help"),
		),

		// Quitting.
		Quit: key.NewBinding(
			key.WithKeys("q", "esc"),
			key.WithHelp("q", "quit"),
		),
		ForceQuit: key.NewBinding(key.WithKeys("ctrl+c")),
	}
}

type Style struct {
	Title                 lipgloss.Style
	FilterPrompt          lipgloss.Style
	FilterCursor          lipgloss.Style
	Status                lipgloss.Style
	ActivePaginationDot   lipgloss.Style
	InactivePaginationDot lipgloss.Style
	HelpStyle             lipgloss.Style
	Content               lipgloss.Style
}

// Model contains the state of this component.
type Model struct {
	showHelp bool
	Title    string
	Styles   Style
	// Key mappings for navigating the list.
	KeyMap Keys
	// Filter is used to filter the list.
	width         int
	maxWidth      int
	height        int
	contentHeight int
	cursor        int
	Help          help.Model
	FilterInput   textinput.Model
	filterState   FilterState
	paginator     Paginator
	// The master set of items we're working with.
	items []Memo
}

// New returns a new model with sensible defaults.
func New(items []Memo, query string, base_color string) Model {
	styles := newDefaultStyle(base_color)

	filterInput := textinput.New()
	filterInput.Prompt = "Filter: "
	filterInput.PromptStyle = styles.FilterPrompt
	filterInput.Cursor.Style = styles.FilterCursor
	filterInput.CharLimit = 64
	filterInput.Focus()
	//handle current filter
	m := Model{
		showHelp:    true,
		KeyMap:      DefaultKeys(),
		Styles:      styles,
		Title:       "List",
		FilterInput: filterInput,
		filterState: Unfiltered,
		items:       items,
		Help:        help.New(),
		width:       50,
		cursor:      0,
	}
	m.Styles.Content = lipgloss.NewStyle().Width(m.width)
	//update pagination
	m.updateHeight()
	m.updatePaginator()
	return m
}
func newDefaultStyle(base_color string) Style {
	verySubduedColor := lipgloss.AdaptiveColor{Light: "#DDDADA", Dark: "#3C3C3C"}
	subduedColor := lipgloss.AdaptiveColor{Light: "#9B9B9B", Dark: "#5C5C5C"}
	var s Style
	s.Title = lipgloss.NewStyle().
		Background(lipgloss.Color("62")).
		Foreground(lipgloss.Color("230")).
		Margin(2, 0, 0, 0).
		Padding(0, 1)
	s.Status = lipgloss.NewStyle().Foreground(subduedColor).Margin(1, 0)
	s.ActivePaginationDot = lipgloss.NewStyle().
		Foreground(subduedColor).
		SetString(bullet)

	s.InactivePaginationDot = lipgloss.NewStyle().
		Foreground(verySubduedColor).
		SetString(bullet)
	s.HelpStyle = lipgloss.NewStyle().Padding(1, 0, 0, 2)
	s.ActivePaginationDot = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#847A85", Dark: "#979797"}).
		MarginTop(1).
		SetString(bullet)

	s.InactivePaginationDot = lipgloss.NewStyle().
		Foreground(verySubduedColor).
		MarginTop(1).
		SetString(bullet)
	return s
}
func (m Model) Init() tea.Cmd {
	return nil
}
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, m.KeyMap.ForceQuit) {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.maxWidth = msg.Width
		m.updateHeight()
		m.updatePaginator()
	}
	cmds = append(cmds, m.handleBrowsing(msg))
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	var sections []string
	sections = append(sections, m.titleView())
	sections = append(sections, m.statusView())
	// sections = append(sections, m.Styles.Content.Render(m.contentView()))
	//maybe joing paginator and help centrally
	sections = append(
		sections,
		lipgloss.JoinVertical(lipgloss.Center, m.contentView(), m.paginatorView()),
	)
	sections = append(sections, m.helpView())
	return lipgloss.Place(
		m.maxWidth,
		m.height,
		lipgloss.Center,
		lipgloss.Center, lipgloss.JoinVertical(lipgloss.Left, sections...))
}
func (m Model) titleView() string {
	return m.Styles.Title.Render("Memos")
}
func (m Model) statusView() string {
	return m.Styles.Status.Render(fmt.Sprintf("%d memos", len(m.items)))
}
func (m Model) contentView() string {
	//fill rest of the height with empty lines
	var (
		sections     []string
		excessHeight = m.contentHeight
	)
	for i, memo := range m.items {
		memoView := memo.RenderSelected(m.width)
		if m.cursor == i {
			memoView = memo.Render(m.width)
		}
		excessHeight -= lipgloss.Height(memoView)
		//add excess height
		sections = append(sections, memoView)
	}
	// sections = append(sections, strings.Repeat("\n", excessHeight))
	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func (m *Model) handleBrowsing(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		// Note: we match clear filter before quit because, by default, they're
		// both mapped to escape.
		// case key.Matches(msg, m.KeyMap.ClearFilter):
		// 	m.resetFiltering()

		case key.Matches(msg, m.KeyMap.Quit):
			return tea.Quit

		case key.Matches(msg, m.KeyMap.CursorUp):
			m.CursorUp()

		case key.Matches(msg, m.KeyMap.CursorDown):
			m.CursorDown()

		// case key.Matches(msg, m.KeyMap.PrevPage):
		// 	m.Paginator.PrevPage()
		//
		// case key.Matches(msg, m.KeyMap.NextPage):
		// 	m.Paginator.NextPage()
		//
		// case key.Matches(msg, m.KeyMap.GoToStart):
		// 	m.Paginator.Page = 0
		// 	m.cursor = 0
		//
		// case key.Matches(msg, m.KeyMap.GoToEnd):
		// 	m.Paginator.Page = m.Paginator.TotalPages - 1
		// 	m.cursor = m.Paginator.ItemsOnPage(numItems) - 1
		//
		// case key.Matches(msg, m.KeyMap.Filter):
		// 	m.hideStatusMessage()
		// 	if m.FilterInput.Value() == "" {
		// 		// Populate filter with all items only if the filter is empty.
		// 		m.filteredItems = m.itemsAsFilterItems()
		// 	}
		// 	m.Paginator.Page = 0
		// 	m.cursor = 0
		// 	m.filterState = Filtering
		// 	m.FilterInput.CursorEnd()
		// 	m.FilterInput.Focus()
		// 	m.updateKeybindings()
		// 	return textinput.Blink

		case key.Matches(msg, m.KeyMap.ShowFullHelp):
			fallthrough
		case key.Matches(msg, m.KeyMap.CloseFullHelp):
			m.Help.ShowAll = !m.Help.ShowAll
			// m.updatePagination()
		}
	}
	return tea.Batch(cmds...)
}

func (m Model) helpView() string {
	return m.Styles.Content.Render(m.Styles.HelpStyle.Render(m.Help.View(m)))
}

func (m Model) paginatorView() string {
	var dots []string
	for i := 0; i < m.paginator.n_pages; i++ {
		if i == m.paginator.current_page {
			dots = append(dots, m.Styles.ActivePaginationDot.String())
		} else {
			dots = append(dots, m.Styles.InactivePaginationDot.String())
		}
	}
	return lipgloss.JoinHorizontal(lipgloss.Center, dots...)
}

func (m *Model) updateHeight() {
	var v string
	m.contentHeight = m.height
	//header
	v = m.titleView()
	m.contentHeight -= lipgloss.Height(v)
	//status
	v = m.statusView()
	m.contentHeight -= lipgloss.Height(v)
	//paginator
	v = m.paginatorView()
	m.contentHeight -= lipgloss.Height(v)
	//help
	v = m.helpView()
	m.contentHeight -= lipgloss.Height(v)
}

// Set keybindings according to the filter state.
func (m *Model) updateKeybindings() {
	switch m.filterState {
	case Filtering:
		m.KeyMap.CursorUp.SetEnabled(false)
		m.KeyMap.CursorDown.SetEnabled(false)
		m.KeyMap.NextPage.SetEnabled(false)
		m.KeyMap.PrevPage.SetEnabled(false)
		m.KeyMap.GoToStart.SetEnabled(false)
		m.KeyMap.GoToEnd.SetEnabled(false)
		m.KeyMap.Filter.SetEnabled(false)
		m.KeyMap.ClearFilter.SetEnabled(false)
		m.KeyMap.CancelWhileFiltering.SetEnabled(true)
		m.KeyMap.AcceptWhileFiltering.SetEnabled(m.FilterInput.Value() != "")
		m.KeyMap.Quit.SetEnabled(false)
		m.KeyMap.ShowFullHelp.SetEnabled(false)
		m.KeyMap.CloseFullHelp.SetEnabled(false)

	default:
		hasItems := len(m.items) != 0
		m.KeyMap.CursorUp.SetEnabled(hasItems)
		m.KeyMap.CursorDown.SetEnabled(hasItems)

		m.KeyMap.NextPage.SetEnabled(true)
		m.KeyMap.PrevPage.SetEnabled(true)

		m.KeyMap.GoToStart.SetEnabled(hasItems)
		m.KeyMap.GoToEnd.SetEnabled(hasItems)

		m.KeyMap.Filter.SetEnabled(true)
		m.KeyMap.ClearFilter.SetEnabled(m.filterState == FilterApplied)
		m.KeyMap.CancelWhileFiltering.SetEnabled(false)
		m.KeyMap.AcceptWhileFiltering.SetEnabled(false)

		if m.Help.ShowAll {
			m.KeyMap.ShowFullHelp.SetEnabled(true)
			m.KeyMap.CloseFullHelp.SetEnabled(true)
		} else {
			m.KeyMap.ShowFullHelp.SetEnabled(false)
			m.KeyMap.CloseFullHelp.SetEnabled(false)
		}
	}
}

func (m Model) ShortHelp() []key.Binding {
	kb := []key.Binding{
		m.KeyMap.CursorUp,
		m.KeyMap.CursorDown,
	}

	// If the delegate implements the help.KeyMap interface add the short help
	// items to the short help after the cursor movement keys.
	kb = append(kb,
		m.KeyMap.Filter,
		m.KeyMap.ClearFilter,
		m.KeyMap.AcceptWhileFiltering,
		m.KeyMap.CancelWhileFiltering,
	)

	return append(kb,
		m.KeyMap.Quit,
		m.KeyMap.ShowFullHelp,
	)
}

// // FullHelp returns bindings to show the full help view. It's part of the
// // help.KeyMap interface.
func (m Model) FullHelp() [][]key.Binding {
	kb := [][]key.Binding{{
		m.KeyMap.CursorUp,
		m.KeyMap.CursorDown,
		m.KeyMap.NextPage,
		m.KeyMap.PrevPage,
		m.KeyMap.GoToStart,
		m.KeyMap.GoToEnd,
	}}

	listLevelBindings := []key.Binding{
		m.KeyMap.Filter,
		m.KeyMap.ClearFilter,
		m.KeyMap.AcceptWhileFiltering,
		m.KeyMap.CancelWhileFiltering,
	}

	return append(kb,
		listLevelBindings,
		[]key.Binding{
			m.KeyMap.Quit,
			m.KeyMap.CloseFullHelp,
		})
}
func (m *Model) updatePaginator() {
	var sublists [][]Memo
	var currentList []Memo
	currentHeight := 0

	for _, memo := range m.items {
		//get height of memos
		memoHeight := lipgloss.Height(memo.Render(m.width))
		if currentHeight+memoHeight <= m.contentHeight {
			currentList = append(currentList, memo)
			currentHeight += memoHeight
		} else {
			sublists = append(sublists, currentList)
			currentList = []Memo{memo}
			currentHeight = memoHeight
		}
	}

	if len(currentList) > 0 {
		sublists = append(sublists, currentList)
	}
	//update paginator
	m.paginator.pages = sublists
	m.paginator.n_pages = len(sublists)
	m.paginator.current_page = 0
	//reset cursor
	m.cursor = 0
}

func (m *Model) CursorUp() {
	m.cursor--
	// If we're at the start, stop
	if m.cursor < 0 && m.paginator.OnFirstPage() {
		m.cursor = 0
		return
	}
	// Move the cursor as normal
	if m.cursor >= 0 {
		return
	}
	// Go to the last item of the the previous page
	m.paginator.PrevPage()
	m.cursor = m.paginator.ItemsOnPage() - 1
}

// CursorDown moves the cursor down. This can also advance the state to the
// next page.
func (m *Model) CursorDown() {
	itemsOnPage := m.paginator.ItemsOnPage()
	m.cursor++
	// If we're not at the end continue
	if m.cursor < itemsOnPage {
		return
	}
	// Go to the next page
	if !m.paginator.OnLastPage() {
		m.paginator.NextPage()
		m.cursor = 0
		return
	}
	// During filtering the cursor position can exceed the number of
	// itemsOnPage. It's more intuitive to start the cursor at the
	// topmost position when moving it down in this scenario.
	if m.cursor > itemsOnPage {
		m.cursor = 0
		return
	}
	m.cursor = itemsOnPage - 1
}
func (p Paginator) OnLastPage() bool {
	return p.current_page == p.n_pages-1
}
func (p Paginator) OnFirstPage() bool {
	return p.current_page == 0
}
func (p *Paginator) NextPage() {
	if !p.OnLastPage() {
		p.current_page++
	}
}
func (p *Paginator) PrevPage() {
	if !p.OnFirstPage() {
		p.current_page--
	}
}
func (p Paginator) ItemsOnPage() int {
	return len(p.pages[p.current_page])
}
