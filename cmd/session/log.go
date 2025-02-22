package session

import (
	"log"
	"math/rand"
	"time"
	"urlAPI/internal/data"
)

func newLogin(Term bool) (string, error) {
	rand.Seed(time.Now().UnixNano())
	acsii := []int{10, 26, 26}
	acsiiPlus := []int{48, 65, 97}
	tk := ""
	for i := 1; i <= 64; i++ {
		choose := rand.Int() % len(acsii)
		tk += string(rand.Int()%acsii[choose] + acsiiPlus[choose])
	}
	var exp time.Time
	if Term == true {
		exp = time.Now().AddDate(0, 0, 7)
	} else {
		exp = time.Now().Add(time.Hour * 12)
	}
	err := data.AddSession(data.DataConfig(
		data.WithSessionToken(tk),
		data.WithSessionExpire(exp),
		data.WithSessionTerm(Term)))
	if err != nil {
		return "", err
	} else {
		return tk, nil
	}
}
func logout(Token string) error {
	return data.DelSession(data.DataConfig(data.WithSessionToken(Token)))
}
func exit(Token string) error {
	sessions, err := data.FetchSession(data.DataConfig(data.WithSessionToken(Token)))
	if err != nil {
		return err
	}
	if sessions[0].Term == false {
		log.Println("Temporary session logout")
		return data.DelSession(data.DataConfig(data.WithSessionToken(Token)))
	} else {
		return nil
	}
}
