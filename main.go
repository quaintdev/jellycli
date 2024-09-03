package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"

	"code.rohanrd.xyz/jellycli/jellyapi"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Config struct {
	Host    string
	UserId  string
	AuthKey string
}

var jellyServer *jellyapi.Server

const version = "1.0.0"

func main() {
	versionFlag := flag.Bool("version", false, "Display version information")
	flag.Parse()

	if *versionFlag {
		fmt.Println("Version:", version)
		return
	}

	var conf Config
	confPath := path.Join(os.Getenv("HOME"), ".config/jellycli.conf")
	confBytes, err := os.Open(confPath)
	if err != nil {
		log.Fatal(err)
	}
	json.NewDecoder(confBytes).Decode(&conf)

	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()
	}

	jellyServer = jellyapi.NewServer(conf.Host, conf.AuthKey, conf.UserId)

	collections, err := jellyServer.GetCollections()
	if err != nil {
		log.Fatal(err)
	}
	var items []list.Item
	for _, c := range collections {
		items = append(items, c)
	}

	m := model{list: list.New(items, list.NewDefaultDelegate(), 0, 0)}
	m.list.Title = "JellyCli"
	m.stack.Push("home")

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

}

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type model struct {
	stack Stack
	list  list.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "ctrl+left":
			id := m.stack.Pop()
			switch id {
			case "home":
				collections, _ := jellyServer.GetCollections()
				var items []list.Item
				for _, c := range collections {
					items = append(items, c)
				}
				m.list.SetItems(items)
			case "":
				break
			default:
				collections, _ := jellyServer.GetChildItems(id)
				var items []list.Item
				for _, c := range collections {
					items = append(items, c)
				}
				m.list.SetItems(items)
			}
		case "enter", " ":
			selectedCollection := m.list.SelectedItem().(jellyapi.Collection)
			parentId := selectedCollection.Id

			var items []list.Item
			if selectedCollection.IsFolder {
				collections, _ := jellyServer.GetChildItems(parentId)
				for _, c := range collections {
					items = append(items, c)
				}
				m.stack.Push(parentId)
				m.list.SetItems(items)
			}
			if selectedCollection.VideoType == "VideoFile" {
				openWithVLC(jellyServer.Host+"/Items/"+selectedCollection.Id+"/Download?api_key="+jellyServer.AuthKey, selectedCollection.Name)
			}
		}
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return docStyle.Render(m.list.View())
}

func openWithVLC(url, name string) error {
	cmd := exec.Command("vlc", url, "--no-video-title-show", "--input-title-format", name)
	return cmd.Start()
}
