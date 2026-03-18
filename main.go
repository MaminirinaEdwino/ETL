package main

import (
	"fmt"
	"os"

	"github.com/MaminirinaEdwino/etl/src/cmd"
	// "github.com/MaminirinaEdwino/etl/src/model"
	tea "github.com/charmbracelet/bubbletea"
)

type FilterModel struct {
	Choices     []string
	cursor      int
	SelectedMap map[int]string
}

func InitialModel(Choices []string) FilterModel {
	return FilterModel{
		Choices: Choices,
		SelectedMap: make(map[int]string),
	}
}

func (m FilterModel) Init() tea.Cmd{
	return nil
}

func (m FilterModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {

    case tea.KeyMsg:
        switch msg.String() {
        case "ctrl+c", "q":
            return m, tea.Quit
        case "up", "k":
            if m.cursor > 0 {
                m.cursor--
            }
        case "down", "j":
            if m.cursor < len(m.Choices)-1 {
                m.cursor++
            }
        case "enter", "space":
            _, ok := m.SelectedMap[m.cursor]
            if ok {
                delete(m.SelectedMap, m.cursor)
            } else {
                m.SelectedMap[m.cursor] = m.Choices[m.cursor]
            }
        }
    }

    return m, nil
}

func (m FilterModel) View() string {
    s := "Extract Transform Load\n\n"
	s+="What do you want to extract\n\n ? "
    for i, choice := range m.Choices {

        cursor := " "
        if m.cursor == i {
            cursor = ">"
        }

        checked := " " 
        if _, ok := m.SelectedMap[i]; ok {
            checked = "x" 
        }
        s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
    }
	for _, el := range m.SelectedMap{
		s+=fmt.Sprintf("%s\n", el)
	}
    s += "\nPress q to quit.\n"

    return s
}

func main() {
	inputFile := "road_accident_data.csv"
	// outputFile := "accidents_clean.json"

	extractor, rowList,err := cmd.NewExtractor(inputFile)
	if err != nil {
		fmt.Printf("Erreur setup: %v\n", err)
		return
	}

	rawRows := make(chan []string, 100)
	// transformedData := make(chan model.RawAccident, 100)
	go extractor.Run(rawRows)
	p := tea.NewProgram(InitialModel(rowList))
	if _, err := p.Run(); err != nil {
        fmt.Printf("Alas, there's been an error: %v", err)
        os.Exit(1)
    }
	// go func() {
	// 	for row := range rawRows {
	// 		acc, err := cmd.TransformRow(row, extractor)
	// 		if err != nil {
	// 			continue
	// 		}
	// 		myConfig := model.FilterConfig{
	// 			MinVehicles: 12,
	// 		}
	// 		if cmd.ShouldKeep(acc, myConfig) {
	// 			transformedData <- acc
	// 		}
	// 	}
	// 	close(transformedData)
	// }()
	// err = cmd.LoadToJSON(outputFile, transformedData)
	// if err != nil {
	// 	fmt.Printf("Erreur lors du chargement: %v\n", err)
	// }
}
