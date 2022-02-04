package app

import (
	"awesomeProject/internal/pkg/http"
	"awesomeProject/internal/pkg/version"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"strings"
)

// Run ...
func Run() {
	defer func() {
		if r := recover(); r != nil {
			log.Fatal("Recovered in main", r)
		}
	}()

	application := kingpin.New("Joker app. ", "Awesome application to get jokes by category or random.")

	log.Debug("Joker app start\nTo interrupt execution press Ctrl+C...")
	application.Version(fmt.Sprintf("Joker version: %s\nCopyright (C) 2021  AwesomeCo\n\n"+
		"Compiler: %s\nOn: %s\nSystem: %s\n", version.Version,
		strings.ReplaceAll(version.Compiler, "_", " "),
		strings.ReplaceAll(version.BuildTime, "_", " "),
		strings.ReplaceAll(version.System, "_", " ")))

	random := application.Command("random", "Get one random joke.")

	dump := application.Command("dump", "Get N jokes by each category.")
	dumpNum := dump.Flag("n", "Amount of jokes in each category").Default("5").Short('n').Int()

	command, err := application.Parse(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	log.Info("Joker app started")

	switch command {
	case random.FullCommand():
		http.GetRandomJoke()
	case dump.FullCommand():
		http.GetJokesByCategories(*dumpNum)
	}

	log.Info("Joker app stop")
}
