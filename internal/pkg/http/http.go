package http

import (
	"awesomeProject/internal/pkg/domain"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

const mainURL = "https://api.chucknorris.io/jokes/random"
const categoriesURL = "https://api.chucknorris.io/jokes/categories"
const maximumAttemptsByAmount = 3

func GetRandomJoke() {
	req, err := http.NewRequest("GET", mainURL, nil)
	if err != nil {
		log.Error(err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
	}

	var joke domain.Joke
	err = json.Unmarshal(bodyBytes, &joke)
	if err != nil {
		log.Error(err)
	}

	fmt.Println(joke.Value)

	gracefulShutdown()
}

func GetJokesByCategories(amount int) {
	categories := GetCategoriesList()
	for _, category := range categories {
		log.Infof("Processing category: %s", category)
		var uniqIds []string
		var jokes []string
		attempts := 1

		for {
			id, value := GetJokeByCategory(category)
			attempts++

			if !contains(uniqIds, id) {
				uniqIds = append(uniqIds, id)
				jokes = append(jokes, value)
			}

			if len(uniqIds) == amount {
				err := os.WriteFile(fmt.Sprintf("%s.txt", category), []byte(strings.Join(jokes[:], "\n")), 0644)
				if err != nil {
					log.Error(err)
				}
				break
			}

			if attempts > amount*maximumAttemptsByAmount {
				err := os.WriteFile(fmt.Sprintf("%s.txt", category), []byte(strings.Join(jokes[:], "\n")), 0644)
				if err != nil {
					log.Error(err)
				}
				break
			}
		}
	}
	gracefulShutdown()
}

func GetJokeByCategory(category string) (string, string) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s?category=%s", mainURL, category), nil)
	if err != nil {
		log.Error(err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
	}

	var joke domain.Joke
	err = json.Unmarshal(bodyBytes, &joke)
	if err != nil {
		log.Error(err)
	}

	return joke.Id, joke.Value
}

func GetCategoriesList() []string {
	req, err := http.NewRequest("GET", categoriesURL, nil)
	if err != nil {
		log.Error(err)
		return nil
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err)
		return nil
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return nil
	}

	var result []string

	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		log.Error(err)
		return nil
	}
	return result
}

func gracefulShutdown() {
	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt)
	signal.Notify(s, syscall.SIGTERM)
	go func() {
		<-s
		os.Exit(0)
	}()
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}
