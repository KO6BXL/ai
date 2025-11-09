package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/KO6BXL/ai"
	openrouter "github.com/KO6BXL/ai/OpenRouter"
)

type env struct {
	Key string `json:"key"`
}

func gimmeKey() string {
	file, err := os.ReadFile("../.env")
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
	or := openrouter.NewOR("deepseek/deepseek-v3.2-exp", gimmeKey())
	or.SetProviders(openrouter.Provider{DataCollection: "allow", Only: []string{"deepseek"}})
	AI := ai.NewAI(or)
	resp, err := AI.Message("What is the square root of pi?")
	if err != nil || resp.Id == "" {
		log.Fatal(err)
	}
	fmt.Println(resp.Outputs[0].Message.Content)
}
