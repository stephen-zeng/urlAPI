package op

import (
	"fmt"
	"math/rand"
	"urlAPI/internal/database"
	"urlAPI/internal/model"
)

func generateRandom(task *model.Task) (GenerateResult, error) {
	content, ok := database.Repos.Get(task.API, task.Target)
	if !ok || len(content) == 0 {
		task.Status = "failed"
		task.Return = "Repo not found"
		return GenerateResult{}, fmt.Errorf("random repository not found")
	}
	result := GenerateResult{URL: content[rand.Intn(len(content))]}
	if err := setTaskResult(task, result); err != nil {
		return GenerateResult{}, err
	}
	return result, nil
}
