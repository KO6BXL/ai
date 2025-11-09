package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type env struct {
	Key string `json:"key"`
}

func gimmeKey() string {
	file, err := os.ReadFile(".env")
	if err != nil {
		log.Fatal(err)
	}
	env := &env{}
	err = json.Unmarshal(file, env)
	if err != nil {
		log.Fatal(errors.Join(err, errors.New("Check .env to include {'key':'or key'}")))
	}
	return env.Key
}

func main() {
	provider := Provider{Only: []string{"deepseek"}}
	bod := NewOR("deepseek/deepseek-v3.2-exp", provider)
	bod.Context.NewUserMessage("Who are you?")

	client := http.Client{}

	req, err := http.NewRequest(http.MethodPost, "https://openrouter.ai/api/v1/chat/completions", bod)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", gimmeKey()))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	or := &ORResponse{}
	err = json.Unmarshal(body, or)
	if err != nil {
		log.Fatal(err)
	}
	//body2, err := io.ReadAll(bod)
	fmt.Printf("%s\n", or.Choices[0].Message.Content)
}
