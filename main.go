package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/MaminirinaEdwino/etl/src/cmd"
	"github.com/MaminirinaEdwino/etl/src/model"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

var (
	fieldType      = []string{"int", "float", "string", "date"}
	fieldOperation = []string{"equal", "less than", "bigger than", "different"}
)
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#01F70D")).
			Border(lipgloss.RoundedBorder())
	inputStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Width(25)
	headingStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#01F70D")).
			Bold(true)
	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Width(30).
			PaddingLeft(1)
	checkedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#01F70D")).
			Background(lipgloss.Color("#01F70D"))
	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#01F70D"))
)

type Filter struct {
	Type      string
	Value     string
	Operation string
}

type FilterModel struct {
	Choices               []string
	cursor                int
	SelectedMap           map[int]string
	SelectedForFilter     string
	OutputFile            string
	Extractor             model.Extractor
	Tab                   int
	TabList               []string
	Filter                map[string]Filter
	ValueInput            textinput.Model
	TypeCursor            int
	OperationCursor       int
	cursorType            string
	Message               string
	ChoicePagination      int
	ChoicePaginationEnd   int
	ChoicePaginationStart int
	ChoiceLimit           int
	TotalPage             int
}

func InitialModel(Choices []string, outputFile string, extractor model.Extractor) FilterModel {
	ti := textinput.New()

	ti.Placeholder = "Enter value here ..."
	ti.CharLimit = 256
	ti.Width = 20
	choiceLimit := 10
	TotalPage := len(Choices) / choiceLimit
	return FilterModel{
		Choices:               Choices,
		SelectedMap:           make(map[int]string),
		OutputFile:            outputFile,
		Extractor:             extractor,
		Tab:                   0,
		Filter:                make(map[string]Filter),
		TabList:               []string{"choices", "filter", "extract"},
		SelectedForFilter:     "",
		ValueInput:            ti,
		ChoiceLimit:           choiceLimit,
		ChoicePagination:      1,
		TotalPage:             TotalPage,
		ChoicePaginationStart: 0,
		ChoicePaginationEnd:   choiceLimit - 1,
	}
}

func (m FilterModel) Init() tea.Cmd {
	return nil
}

func (m FilterModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			if !m.ValueInput.Focused() {
				return m, tea.Quit
			}
		case "up", "k":
			if m.TabList[m.Tab] == "choices" {
				if m.cursor > 0 {
					m.cursor--
					if m.ChoicePaginationStart > 0 {
						m.ChoicePaginationEnd--
						m.ChoicePaginationStart--
					}
				}
			}
			if m.TabList[m.Tab] == "filter" {
				if m.cursorType == "type" && m.TypeCursor > 0 {
					m.TypeCursor--
				}
				if m.cursorType == "operation" && m.OperationCursor > 0 {
					m.OperationCursor--
				}
			}
		case "t":
			if m.TabList[m.Tab] == "filter" && !m.ValueInput.Focused() {
				m.cursorType = "type"
			}
		case "o":
			if m.TabList[m.Tab] == "filter" && !m.ValueInput.Focused() {
				m.cursorType = "operation"
			}
		case "down", "j":
			if m.TabList[m.Tab] == "choices" {
				if m.cursor < len(m.Choices)-1 {
					m.cursor++
					if m.ChoicePaginationEnd < len(m.Choices) {
						m.ChoicePaginationEnd++
						m.ChoicePaginationStart++
					}
				}
			}
			if m.TabList[m.Tab] == "filter" {
				if m.cursorType == "type" && m.TypeCursor < len(fieldType)-1 {
					m.TypeCursor++
				}
				if m.cursorType == "operation" && m.OperationCursor < len(fieldOperation)-1 {
					m.OperationCursor++
				}
			}
		case "enter", "space":
			if m.TabList[m.Tab] == "choices" {
				_, ok := m.SelectedMap[m.cursor]
				if ok {
					delete(m.SelectedMap, m.cursor)
				} else {
					m.SelectedMap[m.cursor] = m.Choices[m.cursor]
				}
			}
		case "f":
			if m.TabList[m.Tab] == "choices" {
				if value, ok := m.SelectedMap[m.cursor]; ok {
					if _, ok = m.Filter[m.SelectedForFilter]; ok {
						delete(m.Filter, value)
						m.SelectedForFilter = ""
					} else {
						m.Filter[value] = Filter{}
						m.SelectedForFilter = value
					}
				}
			}
			if m.TabList[m.Tab] == "filter" {
				if m.SelectedForFilter != "" {
					m.Filter[m.SelectedForFilter] = Filter{
						Value:     m.ValueInput.Value(),
						Type:      fieldType[m.TypeCursor],
						Operation: fieldOperation[m.OperationCursor],
					}
				}
			}
		case "ctrl+e":
			m.Message = ""
			ExtractData(&m.Extractor, m.OutputFile, m)
			m.Message = "Extraction complete . . ."
		case "ctrl+left":
			if m.Tab > 0 {
				m.Tab--
			}
		case "tab":
			if m.TabList[m.Tab] == "filter" && !m.ValueInput.Focused() {
				m.ValueInput.Focus()
			} else {
				m.ValueInput.Blur()
			}
		case "ctrl+right":
			if m.Tab < len(m.TabList)-1 {
				m.Tab++
			}
		}
	}
	var cm tea.Cmd
	if m.ValueInput.Focused() {
		m.ValueInput, cm = m.ValueInput.Update(msg)
	}
	return m, cm
}

func (m FilterModel) View() string {
	var s strings.Builder
	s.WriteString("\n")
	s.WriteString(titleStyle.Render("Extract Transform Load"))
	s.WriteString("\n")

	switch m.TabList[m.Tab] {
	case "choices":
		s.WriteString("What do you want to extract ? \n\n")
		for i, choice := range m.Choices {
			if i >= m.ChoicePaginationStart && i <= m.ChoicePaginationEnd {
				cursor := " "
				content := choice
				if m.cursor == i {
					cursor = ">"
					content = selectedStyle.Render(content)
				}
				checked := " "
				if _, ok := m.SelectedMap[i]; ok {
					checked = checkedStyle.Render(" ")
				}
				filtered := " "
				if _, ok := m.Filter[choice]; ok {
					filtered = "F"
				}
				fmt.Fprintf(&s, "%s [%s] [%s] %s\n", cursor, checked, filtered, content)
			}
		}
		fmt.Fprintln(&s, "up/down : navigate inside the field lists\nenter: for selecting the field that you want to extract\nf: for adding a filter for the fields\nctrl+right: navigate to the filter section\nctrl+c/q : quit the program")
	case "filter":
		var t, o strings.Builder
		fmt.Fprintf(&s, "Filter %s \n", headingStyle.Render(m.SelectedForFilter))
		fmt.Fprintf(&s, "%s", inputStyle.Render(m.ValueInput.View()))
		fmt.Fprintln(&s, "")

		fmt.Fprint(&t, "Field Type\n")
		for i, element := range fieldType {
			cursor := " "
			content := element
			if m.TypeCursor == i {
				cursor = ">"
				content = selectedStyle.Render(content)
			}
			fmt.Fprintf(&t, "%s %s\n", cursor, content)
		}
		// fmt.Fprintln(&s, boxStyle.Render(t.String()))
		fmt.Fprintln(&o, "Filter Operation")
		for i, element := range fieldOperation {
			cursor := " "
			content := element
			if m.OperationCursor == i {
				cursor = ">"
				content = selectedStyle.Render(content)
			}
			fmt.Fprintf(&o, "%s %s\n", cursor, content)
		}
		// fmt.Fprintln(&s, boxStyle.Render(o.String()))
		row := lipgloss.JoinHorizontal(lipgloss.Top, boxStyle.Render(t.String()), boxStyle.Render(o.String()))
		fmt.Fprintln(&s, row)
		fmt.Fprintln(&s, "Actual Filter ")
		for i, value := range m.Filter {
			fmt.Fprintf(&s, "%s %s %s %s\n", i, value.Value, value.Type, value.Operation)
		}
		fmt.Fprintln(&s, "up/down : choose the type and the operation for the filter(you don't need to check this time your choices this tilme, so just place the cursor near your choices)\ntab: focus and un-focus the input\nf: add the filter parameters to the filter list\nctrl+right: navitage to the extract section\n ctrl+left: go back to the filter section\n ctrl+c/q : quit the program")
	case "extract":
		fmt.Fprint(&s, "Extract\n")
		fmt.Fprintln(&s, m.Message)
		var header []string
		var rows [][]string
		count := 0
		outFile, _ := os.Open(m.OutputFile)

		defer outFile.Close()
		// scanner := bufio.NewScanner(outFile)
		decoder := json.NewDecoder(outFile)
		tmp := make(map[string]string)
		for _, value := range m.SelectedMap {
			header = append(header, value)
		}
		for count < 15 {
			decoder.Decode(&tmp)
			var tmpTab []string
			found := false
			for _, el := range header {
				tmpTab = append(tmpTab, tmp[el])
			}
			for _, row := range rows {
				if slices.Equal(row, tmpTab) {
					found = true
					break
				}
			}
			if !found {
				rows = append(rows, tmpTab)
			}
			count++
		}

		t := table.New().
			Border(lipgloss.NormalBorder()).
			BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
			Headers(header...).
			Rows(rows...)
		fmt.Fprintln(&s, t.Render())
	}

	//s.WriteString("\nPress \nq to quit.\ne to see extract tab\nc to switch to choice tab\n")
	fmt.Fprintln(&s, "ctrl+e: load the result into a json file\nctrl+e/q :quit the program")
	return s.String()
}
func ShouldKeep(acc map[string]string, filter map[string]Filter) bool {
	okMap := make(map[string]bool)
	for i, value := range acc {
		if _, ok := filter[i]; ok {
			switch filter[i].Type {
			case "string":
				switch filter[i].Operation {
				case "equal":
					if value == filter[i].Value {
						okMap[i] = true
					}
				}
			case "int":
				filterValue, _ := strconv.Atoi(filter[i].Value)
				realValue, _ := strconv.Atoi(value)
				switch filter[i].Operation {
				case "equal":
					if filterValue == realValue {
						okMap[i] = true
					}
				case "less than":
					if filterValue < realValue {
						okMap[i] = true
					}
				case "bigger than":
					if filterValue > realValue {
						okMap[i] = true
					}
				case "different":
					if filterValue != realValue {
						okMap[i] = true
					}
				}
			case "float":
			case "date":
				filterValue, _ := time.Parse("01/02/2006", filter[i].Value)
				realValue, _ := time.Parse("01/02/2006", value)
				switch filter[i].Operation {
				case "equal":
					if filterValue.Equal(realValue) {
						okMap[i] = true
					}
				case "less than":
					if filterValue.Before(realValue) {
						okMap[i] = true
					}
				case "bigger than":
					if filterValue.After(realValue) {
						okMap[i] = true
					}
				case "different":
					if filterValue != realValue {
						okMap[i] = true
					}
				}

			}
		}
	}
	for i := range filter {
		if _, ok := okMap[i]; !ok {
			return false
		}
	}
	return true
}
func ExtractData(extractor *model.Extractor, outputFile string, m FilterModel) {
	rawRows := make(chan []string, 100)
	transformedData := make(chan map[string]string, 100)
	go extractor.Run(rawRows)

	go func() {
		for row := range rawRows {
			acc, err := cmd.TransformRow(row, extractor, m.SelectedMap)
			if err != nil {
				continue
			}

			if len(m.Filter) > 0 {
				if ShouldKeep(acc, m.Filter) {
					transformedData <- acc
				}
			} else {
				transformedData <- acc
			}
		}
		close(transformedData)
	}()
	err := cmd.LoadToJSON(outputFile, transformedData)
	if err != nil {
		fmt.Printf("Erreur lors du chargement: %v\n", err)
	}
}

func main() {
	inputFile := flag.String("inputfile", "", "The source file")
	outputFile := flag.String("outputfile", "", "The name of the file for the loaading the result")
	// inputFile := "road_accident_data.csv"
	// outputFile := "accidents_clean.json"
	flag.Parse()
	switch {
	case *inputFile != "" && *outputFile != "":
		extractor, rowList, err := cmd.NewExtractor(*inputFile)
		if err != nil {
			fmt.Printf("Erreur setup: %v\n", err)
			return
		}
		p := tea.NewProgram(InitialModel(rowList, *outputFile, *extractor))
		if _, err := p.Run(); err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
			os.Exit(1)
		}
	default:
		fmt.Printf("Use the following command:\netl --inputfile=\"yourfile.csv\" --outputfile=\"yourfile.json\"")
	}
}
