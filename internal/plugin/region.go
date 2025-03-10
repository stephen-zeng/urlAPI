package plugin

import (
	"encoding/json"
	"io"
	"urlAPI/internal/client"
)

func GetRegion(dat Config) (PluginResponse, error) {
	url := "https://api.vore.top/api/IPdata?ip=" + dat.IP
	resp, err := client.GlobalHTTPClient.Get(url)
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
	var response regionResp
	err = json.Unmarshal(jsonResponse, &response)
	if err != nil {
		return PluginResponse{
			Region: "Unknown",
		}, err
	} else {
		return PluginResponse{
			Region: response.IPData.Info1,
		}, nil
	}
}
