package main

import (
	"encoding/json"
	"fmt"
	"github.com/ysugimoto/husky"
	"handler"
	"io/ioutil"
	"os"
)

type AppConfig struct {
	Host   string `json:"host"`
	Port   int    `json:"port"`
	Path   string `json:"path"`
	DbUser string `json:"db_user"`
	DbPass string `json:"db_pass"`
	DbHost string `json:"db_host"`
	DbPort int    `json:"db_port"`
	DbName string `json:"db_name"`
}

func main() {
	cwd, _ := os.Getwd()
	app := husky.NewApp()

	// CLI Options
	app.Command.Alias("c", "config", cwd+"/config.json")
	app.Command.Parse(os.Args[1:])

	config, _ := app.Command.GetOption("config")
	if _, err := os.Stat(config.(string)); err == nil {
		if buffer, err := ioutil.ReadFile(config.(string)); err == nil {
			conf := AppConfig{}
			if err := json.Unmarshal(buffer, &conf); err == nil {
				app.Config.Set("host", conf.Host)
				app.Config.Set("port", conf.Port)
				app.Config.Set("path", conf.Path)

				dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
					conf.DbUser,
					conf.DbPass,
					conf.DbHost,
					conf.DbPort,
					conf.DbName,
				)

				handler.SetDSN(dsn)
			} else {
				fmt.Printf("%v", err)
				return
			}
		}
	}

	app.AcceptCORS([]string{"X-Requested-With", "X-LAP-Token"})
	app.Post("/add_rss_category", handler.AddRssCategory)
	app.Post("/add_rss", handler.AddRss)
	app.Post("/add", handler.Add)
	app.Get("/accept", handler.Accept)
	app.Get("/search", handler.Search)

	app.Serve()
}
