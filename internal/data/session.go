package data

import "log"

func InitSession(data Config) error {
	err := db.AutoMigrate(&Session{})
	if err != nil {
		return err
	}
	if data.Type == "restore" {
		err := db.Where("1 = 1").Delete(&Session{})
		if err.Error != nil {
			return err.Error
		} else {
			log.Println("Initialized Session")
			return nil
		}
	}
	return nil
}
func AddSession(data Config) error {
	err := addSession(
		data.SessionToken,
		data.SessionTerm,
		data.SessionExpire)
	if err != nil {
		log.Println(err)
		return err
	} else {
		return nil
	}
}
func DelSession(data Config) error {
	err := delSession(data.SessionToken)
	if err != nil {
		log.Println(err)
		return err
	} else {
		return nil
	}
}
func FetchSession(data Config) ([]Session, error) {
	return fetchSession(data.SessionToken)
}
