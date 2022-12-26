package playground

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

//go:embed docker-compose.yml.tmpl
var templateFS embed.FS

type Relay struct {
	Name         string
	DatabaseName string
	DatabasePort int
	Port         int
}

type Template struct {
	Relays []Relay
}

var relayPort = 2700
var dbPort = 5432

func Generate(count int) error {
	var relays []Relay
	for i := 0; i < count; i++ {
		relays = append(relays, Relay{
			Name:         fmt.Sprintf("relay_%v", i),
			DatabaseName: fmt.Sprintf("relay_%v_pg", i),
			DatabasePort: dbPort + i,
			Port:         relayPort + i,
		})
	}

	data := Template{Relays: relays}

	tmpl, err := template.ParseFS(templateFS, "docker-compose.yml.tmpl")
	if err != nil {
		return err
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	fkitDir := filepath.Join(home, ".fkit")
	if err := os.RemoveAll(fkitDir); err != nil {
		return err
	}

	if err := os.Mkdir(filepath.Join(home, ".fkit"), os.ModePerm); err != nil {
		return err
	}

	out, err := os.Create(filepath.Join(fkitDir, "docker-compose.yml"))
	if err != nil {
		return err
	}

	return tmpl.Execute(out, data)
}
