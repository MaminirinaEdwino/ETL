package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/MaminirinaEdwino/etl/src/cmd"
	"github.com/MaminirinaEdwino/etl/src/model"
	tea "github.com/charmbracelet/bubbletea"
)

type FilterModel struct {
	Choices     []string
	cursor      int
	SelectedMap map[int]string
	OutputFile  string
	Extractor   model.Extractor
	Tab         string
}

func InitialModel(Choices []string, outputFile string, extractor model.Extractor) FilterModel {
	return FilterModel{
		Choices:     Choices,
		SelectedMap: make(map[int]string),
		OutputFile:  outputFile,
		Extractor:   extractor,
		Tab:         "choices",
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
			return m, tea.Quit
		case "up", "k":
			if m.Tab == "choices" {
				if m.cursor > 0 {
					m.cursor--
				}
			}
		case "down", "j":
			if m.Tab == "choices" {
				if m.cursor < len(m.Choices)-1 {
					m.cursor++
				}
			}
		case "enter", "space":
			if m.Tab == "choices" {
				_, ok := m.SelectedMap[m.cursor]
				if ok {
					delete(m.SelectedMap, m.cursor)
				} else {
					m.SelectedMap[m.cursor] = m.Choices[m.cursor]
				}
			}
		case "ctrl+e":
			ExtractData(&m.Extractor, m.OutputFile, m)
		case "c":
			m.Tab = "choices"
		case "e":
			m.Tab = "extract"
		}
	}

	return m, nil
}

func (m FilterModel) View() string {
	var s strings.Builder
	s.WriteString("Extract Transform Load\n\n")

	switch m.Tab {
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
			fmt.Fprintf(&s, "%s [%s] %s\n", cursor, checked, choice)
		}
	case "extract":
		fmt.Fprint(&s, "Extract")
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
