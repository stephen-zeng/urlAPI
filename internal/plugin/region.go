package plugin

import (
	"encoding/json"
	"io"
	"net/http"
)

func GetRegion(dat Config) (PluginResponse, error) {
	url := "https://api.vore.top/api/IPdata?ip=" + dat.IP
	resp, err := http.Get(url)
	if err != nil {
		return PluginResponse{
			Region: "Unknown",
		}, err
	}
	defer resp.Body.Close()
	jsonResponse, err := io.ReadAll(resp.Body)
	if err != nil {
		return PluginResponse{
			Region: "Unknown",
		}, err
	}
	var response map[string]interface{}
	err = json.Unmarshal(jsonResponse, &response)
	if err != nil {
		return PluginResponse{
			Region: "Unknown",
		}, err
	} else {
		return PluginResponse{
			Region: response["ipdata"].(map[string]interface{})["info1"].(string),
		}, nil
	}
}
