package cmd

import (
	"fmt"
	"log"
	"urlAPI/internal/auth"
	"urlAPI/internal/bootstrap"
	"urlAPI/internal/database"
	"urlAPI/util"
)

func admin(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("missing admin command")
	}
	if err := bootstrap.Init(); err != nil {
		return err
	}
	defer bootstrap.Release()
	switch args[0] {
	case "repwd":
		password, err := resetPassword()
		if err != nil {
			return err
		}
		fmt.Printf("The dashboard password has been reset to: %s\n"+
			"It is shown only once. All sessions were logged out; log in and change it from the dashboard.\n", password)
	case "clear":
		database.ClearTask()
		log.Println("Task Cleared")
	case "logout":
		database.ClearSession()
		log.Println("Session Restored")
	case "clear_ip_restriction":
		clearIPRestrict()
		log.Println("Cleared IP restriction")
	default:
		return fmt.Errorf("unknown admin command %q", args[0])
	}
	return nil
}

// resetPassword replaces the dashboard password with a freshly generated
// one, logs out every session and returns the new password.
func resetPassword() (string, error) {
	password, err := auth.GeneratePassword()
	if err != nil {
		return "", err
	}
	hash, err := auth.HashPassword(auth.ClientCredential(password))
	if err != nil {
		return "", err
	}
	err = database.UpdateAppSettings(func(settings *util.AppSettings) error {
		settings.Security.DashboardPasswordHash = hash
		return nil
	})
	if err != nil {
		return "", err
	}
	database.ClearSession()
	return password, nil
}

func clearIPRestrict() {
	err := database.UpdateAppSettings(func(settings *util.AppSettings) error {
		settings.Security.DashboardAllowedIPs = []string{"*"}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}
