package command

import (
	"encoding/json"
	"log"
	"urlAPI/database"
)

func Arg(args []string) {
	for index, arg := range args {
		if index == 1 {
			continue
		}
		switch arg {
		case "port":
			Port = args[index+1]
		case "repwd":
			repwd()
			log.Println("Password has been reset to 123456, please change it ASAP.")
		case "clear":
			database.ClearTask()
			log.Println("Cleared")
		case "restore":
			database.ClearSession()
			database.ClearSetting()
			log.Println("Restored")
		case "clear_ip_restriction":
			log.Println("Cleared IP restriction")
		}
	}
}

func repwd() {
	dbSettingList := database.SettingMap["pwd"]
	dbSettingList[0] = "8d9f6a89e5e1daab9225e92650a8caf918e38161a4ce23fea07de1bc8fc378a9"
	jsonList, _ := json.Marshal(dbSettingList)
	dbWriter := database.Setting{
		Name:  "pwd",
		Value: string(jsonList),
	}
	if err := dbWriter.Update(); err != nil {
		log.Fatal(err)
	}
}

func clearIPRestrict() {
	dbWriter := database.Setting{
		Name:  "dashallowedip",
		Value: `["*"]`,
	}
	if err := dbWriter.Update(); err != nil {
		log.Fatal(err)
	}
}
