package plugin

import (
	"math/rand"
	"time"
	"urlAPI/internal/data"
)

func random(API, Info string) (PluginResponse, error) {
	content, err := data.FetchRepo(data.DataConfig(
		data.WithBy("api&info"),
		data.WithAPI(API),
		data.WithRepoInfo(Info)))
	if err != nil {
		return PluginResponse{}, err
	}
	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(content[0].Content))
	return PluginResponse{
		URL: content[0].Content[index],
	}, err
}
