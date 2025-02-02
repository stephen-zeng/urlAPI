package data

import "log"

func InitSession(data Config) error {
	if data.Type != "restore" && db.Migrator().HasTable(&Session{}) {
		return nil
	} else {
		err := db.AutoMigrate(&Session{})
		if err != nil {
			return err
		} else {
			log.Println("Initialized Session")
			return nil
		}
	}
}
func AddSession(data Config) error {
	err := addSession(
		data.Token,
		data.Term,
		data.Expire)
	if err != nil {
		log.Println(err)
		return err
	} else {
		return nil
	}
}
func DelSession(data Config) error {
	err := delSession(data.Token)
	if err != nil {
		log.Println(err)
		return err
	} else {
		return nil
	}
}
func FetchSession(data Config) ([]Session, error) {
	return fetchSession(data.Token)
}
