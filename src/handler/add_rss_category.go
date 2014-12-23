package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/ysugimoto/husky"
)

type RssCategory struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func AddRssCategory(d *husky.Dispatcher) {
	db := husky.NewDb(GetDSN())

	req := d.Input.GetRequest()
	req.ParseForm()

	category := req.FormValue("category")

	// check token
	token := req.Header.Get("X-LAP-Token")
	if token == "" {
		db.TransRollback()
		message := "Token not found"
		SendError(d, message)
		return
	}

	// Check token and get userId
	var userId int
	row := db.Select("id").Where("token", "=", token).GetRow("pb_users")
	if err := row.Scan(&userId); err != nil || err == sql.ErrNoRows {
		db.TransRollback()
		message := "Token not matched!"
		SendError(d, message)
		return
	}

	// Check duplicate
	exists := db.Select("id").Where("name", "=", category).Where("user_id", "=", userId).GetRow("pb_rss_categories")
	var id int
	if err := exists.Scan(&id); err != sql.ErrNoRows {
		message := category + " is already exists!"
		SendError(d, message)
		return
	}

	db.TransBegin()

	// Insert URL
	insert := map[string]interface{}{
		"name":       category,
		"user_id":    userId,
		"created_at": getCurrentDateTime(),
	}
	if _, err := db.Insert("pb_rss_categories", insert); err != nil {
		db.TransRollback()
		message := fmt.Sprintf("Query Error: %v\n", err)
		SendError(d, message)
		return
	}

	db.TransCommit()

	rows, _ := db.Select("id", "name").Where("user_id", "=", userId).Get("pb_rss_categories")
	var response []RssCategory
	for rows.Next() {
		rc := RssCategory{}
		rows.Scan(&rc.Id, &rc.Name)

		response = append(response, rc)
	}
	if str, err := json.Marshal(response); err != nil {
		message := fmt.Sprintf("Query Error: %v\n", err)
		SendError(d, message)
	} else {
		SendOK(d, string(str))
	}
}
