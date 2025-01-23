package data

import (
	"database/sql"
	"fmt"
	uuid2 "github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func connect() {
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
		fmt.Println("Error connecting to database")
	}
	//defer db.Close()
	fmt.Println("Connected to database")
}

func dbInit() {
	exec := "CREATE TABLE IF NOT EXISTS task (UUID TEXT, time TEXT, IP TEXT, type TEXT, status TEXT, target TEXT, return TEXT, PRIMARY KEY(UUID))"
	_, err := db.Exec(exec)
	if err != nil {
		fmt.Println("Error creating table")
		log.Fatal(err)
		return
	}
	fmt.Println("Initialized database")
}

func Disconnect() {
	defer db.Close()
	fmt.Println("Disconnected from database")
}

func add(data map[string]string) {
	fmt.Println("Adding data to database...")
	uuid := uuid2.New().String()
	insert := "INSERT INTO task (UUID, time, IP, type, status, target) VALUES (?, ?, ?, ?, ?, ?)"
	_, err := db.Exec(insert, uuid, data["time"], data["ip"], data["type"], data["status"], data["target"])
	if err != nil {
		fmt.Println("Error inserting data")
		log.Fatal(err)
		return
	}
	fmt.Println(data)
}

func del(data map[string]string) {
	fmt.Println("Deleting data from database...")
	if uuid, exists := data["UUID"]; exists {
		_, err := db.Exec("DELETE FROM task WHERE UUID = ?", uuid)
		if err != nil {
			fmt.Println("Error deleting data")
			log.Fatal(err)
			return
		}
	} else {
		fmt.Println("Data doesn't exists")
	}
}

func edit(data map[string]string) {
	fmt.Println("Editing data from database...")
	if uuid, exists := data["UUID"]; exists {
		change := "UPDATE task SET "
		if time, exists := data["time"]; exists {
			change += "time=" + time
		}
		if ip, exists := data["ip"]; exists {
			change += ", ip=" + ip
		}
		if typ, exists := data["type"]; exists {
			change += ", type=" + typ
		}
		if status, exists := data["status"]; exists {
			change += ", status=" + status
		}
		if target, exists := data["target"]; exists {
			change += ", target=" + target
		}
		if ret, exists := data["return"]; exists {
			change += ", return=" + ret
		}
		change += " WHERE UUID =" + uuid
		_, err := db.Exec(change)
		if err != nil {
			fmt.Println("Error updating data")
			log.Fatal(err)
			return
		}
	} else {
		fmt.Println("Data doesn't exists")
	}
}

func get(data map[string]string) []map[string]string {
	fmt.Println("Getting data from database...")
	var result *sql.Rows
	var err error
	if uuid, exists := data["uuid"]; exists {
		fmt.Println("getting by UUID")
		result, err = db.Query("SELECT * FROM task WHERE UUID = ?", uuid)
	} else if ip, exists := data["ip"]; exists {
		fmt.Println("getting by IP")
		result, err = db.Query("SELECT * FROM task WHERE IP = ?", ip)
	} else if time, exists := data["time"]; exists {
		fmt.Println("getting by time")
		result, err = db.Query("SELECT * FROM task WHERE time = ?", time)
	} else if typ, exists := data["type"]; exists {
		fmt.Println("getting by type")
		result, err = db.Query("SELECT * FROM task WHERE type = ?", typ)
	} else if status, exists := data["status"]; exists {
		fmt.Println("getting by status")
		result, err = db.Query("SELECT * FROM task WHERE status = ?", status)
	} else if target, exists := data["target"]; exists {
		fmt.Println("getting by target")
		result, err = db.Query("SELECT * FROM task WHERE target = ?", target)
	} else {
		fmt.Println("Data doesn't exists")
		return []map[string]string{}
	}
	if err != nil {
		fmt.Println("Error getting data")
		log.Fatal(err)
		return []map[string]string{}
	}
	var ret []map[string]string
	for result.Next() { // map的键值对不能传指针
		var uuid, time, ip, typ, status, target string
		var retItem sql.NullString
		err := result.Scan(&uuid, &time, &ip, &typ, &status, &target, &retItem)
		if err != nil {
			fmt.Println("Error reading data")
			log.Fatal(err)
			return []map[string]string{}
		}
		item := map[string]string{
			"uuid":   uuid,
			"time":   time,
			"ip":     ip,
			"type":   typ,
			"status": status,
			"target": target,
			"return": "",
		}
		if retItem.Valid {
			item["return"] = retItem.String
		}
		fmt.Println(item)
		ret = append(ret, item)
	}
	return ret
}
