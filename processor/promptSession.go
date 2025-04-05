package processor

import (
	"encoding/json"
	"github.com/pkg/errors"
	"urlAPI/database"
	"urlAPI/util"
)

func promptDBGetter(For string, Letter int) ([]string, error) {
	dbReader := database.Prompt{
		For:    For,
		Letter: Letter,
	}
	dbAllList, err := dbReader.Read()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var prompts []string
	_ = json.Unmarshal([]byte(dbAllList.PromptList[0].Prompts), &prompts)
	return prompts, nil
}

func promptDBWriter(For string, Letter int, Prompts []string) error {
	jsonData, _ := json.Marshal(Prompts)
	dbWriter := database.Prompt{
		For:     For,
		Letter:  Letter,
		Prompts: string(jsonData),
	}
	return errors.WithStack(dbWriter.Update())
}

func fetchPrompt(info *Session, data *database.Session) error {
	var err error
	info.PromptData, err = promptDBGetter(info.PromptFor, info.PromptIndex)
	return errors.WithStack(err)
}

func newPrompt(info *Session, data *database.Session) error {
	index := util.GetKeyIndex(info.PromptData[0])
	prompts, err := promptDBGetter(info.PromptFor, index)
	if err != nil {
		return errors.WithStack(err)
	}
	prompts = append(prompts, info.PromptData[0])
	if err = promptDBWriter(info.PromptFor, index, prompts); err != nil {
		return errors.WithStack(err)
	}
	info.SettingPart = info.PromptFor + "Propmt"
	info.SettingData = [][]string{prompts}
	if err = editSetting(info, data); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func editPrompt(info *Session, data *database.Session) error {
	if err := promptDBWriter(info.PromptFor, info.PromptIndex, info.PromptData); err != nil {
		return errors.WithStack(err)
	}
	info.SettingPart = info.PromptFor + "Propmt"
	info.SettingData = [][]string{info.PromptData}
	if err := editSetting(info, data); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
