package session

import (
	"backend/internal/data"
	"log"
	"math/rand"
	"time"
)

func login(dat Config) (string, error) {
	rand.Seed(time.Now().UnixNano())
	acsii := []int{10, 26, 26}
	acsiiPlus := []int{48, 65, 97}
	tk := ""
	for i := 1; i <= 64; i++ {
		choose := rand.Int() % len(acsii)
		tk += string(rand.Int()%acsii[choose] + acsiiPlus[choose])
	}
	var exp time.Time
	if dat.Term == true {
		exp = time.Now().AddDate(0, 0, 7)
	} else {
		exp = time.Now().Add(time.Hour * 12)
	}
	err := data.AddSession(data.DataConfig(
		data.WithToken(tk),
		data.WithExpire(exp),
		data.WithTerm(dat.Term)))
	if err != nil {
		return "", err
	} else {
		return tk, nil
	}
}
func logout(dat Config) error {
	return data.DelSession(data.DataConfig(data.WithToken(dat.Token)))
}
func exit(dat Config) error {
	sessions, err := data.FetchSession(data.DataConfig(data.WithToken(dat.Token)))
	if err != nil {
		return err
	}
	if sessions[0].Term == false {
		log.Println("Temporary session logout")
		return data.DelSession(data.DataConfig(data.WithToken(dat.Token)))
	} else {
		return nil
	}
}
