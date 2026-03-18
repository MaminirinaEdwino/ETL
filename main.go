package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/MaminirinaEdwino/etl/src/cmd"
	"github.com/MaminirinaEdwino/etl/src/model"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	fieldType = []string{"int", "float", "string"}
	fieldOperation = []string{"equal", "less than", "higher than", "different"}
)

type Filter struct {
	Type      string
	Value     string
	Operation string
}

type FilterModel struct {
	Choices           []string
	cursor            int
	SelectedMap       map[int]string
	SelectedForFilter string
	OutputFile        string
	Extractor         model.Extractor
	Tab               int
	TabList           []string
	Filter            map[string]Filter
	ValueInput        textinput.Model
	TypeCursor        int
	OperationCursor   int
	cursorType        string
}

func InitialModel(Choices []string, outputFile string, extractor model.Extractor) FilterModel {
	ti := textinput.New()

	ti.Placeholder = "Enter value here ..."
	ti.CharLimit = 256
	ti.Width = 20

	return FilterModel{
		Choices:           Choices,
		SelectedMap:       make(map[int]string),
		OutputFile:        outputFile,
		Extractor:         extractor,
		Tab:               0,
		Filter:            make(map[string]Filter),
		TabList:           []string{"choices", "filter", "extract"},
		SelectedForFilter: "",
		ValueInput:        ti,
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
				}
			}
			if m.TabList[m.Tab] == "filter" {
				if m.cursorType == "type"{
					m.TypeCursor--
				}
				if m.cursorType == "operation"{
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
				}
			}
			if m.TabList[m.Tab] == "filter" {
				if m.cursorType == "type"{
					m.TypeCursor++
				}
				if m.cursorType == "operation"{
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
					if m.SelectedForFilter != "" {
						if _, ok = m.Filter[m.SelectedForFilter]; ok {
							delete(m.Filter, m.SelectedForFilter)
							m.SelectedForFilter = ""
						} else {
							m.Filter[value] = Filter{}
							m.SelectedForFilter = value
						}
					} else {
						m.Filter[value] = Filter{}
						m.SelectedForFilter = value
					}
				}
			}
			if  m.TabList[m.Tab] == "filter"  {
				if m.SelectedForFilter != "" {
					m.Filter[m.SelectedForFilter] = Filter{
						Value: m.ValueInput.Value(),
						Type: fieldType[m.TypeCursor],
						Operation: fieldOperation[m.OperationCursor],
					}
				}
			}
		case "ctrl+e":
			ExtractData(&m.Extractor, m.OutputFile, m)
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
	s.WriteString("Extract Transform Load\n\n")

	switch m.TabList[m.Tab] {
	case "choices":
		s.WriteString("What do you want to extract ? \n\n")
		for i, choice := range m.Choices {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
			}
			checked := " "
			if _, ok := m.SelectedMap[i]; ok {
				checked = "x"
			}
			filtered := " "
			if _, ok := m.Filter[choice]; ok {
				filtered = "F"
			}
			fmt.Fprintf(&s, "%s [%s] [%s] %s\n", cursor, checked, filtered, choice)
		}
	case "filter":
		fmt.Fprintf(&s, "Filter %s \n\n", m.SelectedForFilter)
		fmt.Fprintf(&s, "%s", m.ValueInput.View())
		fmt.Fprintln(&s, "")
		fmt.Fprint(&s, "Field Type\n")
		for i, element := range fieldType {
			cursor := " "
			if m.TypeCursor == i {
				cursor = ">"
			}
			fmt.Fprintf(&s, "%s %s", cursor, element)
		}
		for i, element := range fieldOperation {
			cursor := " "
			if m.OperationCursor == i {
				cursor = ">"
			}
			fmt.Fprintf(&s, "%s %s", cursor, element)
		}
	case "extract":
		fmt.Fprint(&s, "Extract\n")

	}

	s.WriteString("\nPress \nq to quit.\ne to see extract tab\nc to switch to choice tab\n")

	return s.String()
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

			transformedData <- acc
		}
		close(transformedData)
	}()
	err := cmd.LoadToJSON(outputFile, transformedData)
	if err != nil {
		fmt.Printf("Erreur lors du chargement: %v\n", err)
	}
}

func main() {
	inputFile := "road_accident_data.csv"
	outputFile := "accidents_clean.json"

	extractor, rowList, err := cmd.NewExtractor(inputFile)
	if err != nil {
		fmt.Printf("Erreur setup: %v\n", err)
		return
	}
	p := tea.NewProgram(InitialModel(rowList, outputFile, *extractor))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
