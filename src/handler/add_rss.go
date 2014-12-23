package handler

import (
	"database/sql"
	"fmt"
	"github.com/ysugimoto/husky"
)

func AddRss(d *husky.Dispatcher) {
	db := husky.NewDb(GetDSN())

	req := d.Input.GetRequest()
	req.ParseForm()

	categoryId := req.FormValue("category")
	url := req.FormValue("url")
	title := req.FormValue("title")

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
	exists := db.Select("id").Where("url", "=", url).Where("user_id", "=", userId).GetRow("pb_rss_urls")
	var id int
	if err := exists.Scan(&id); err != sql.ErrNoRows {
		message := "URL is already registered!"
		SendError(d, message)
		return
	}

	db.TransBegin()

	// Insert URL
	insert := map[string]interface{}{
		"url":         url,
		"title":       title,
		"user_id":     userId,
		"category_id": categoryId,
		"created_at":  getCurrentDateTime(),
	}
	if _, err := db.Insert("pb_rss_urls", insert); err != nil {
		db.TransRollback()
		message := fmt.Sprintf("Query Error: %v\n", err)
		SendError(d, message)
		return
	}

	db.TransCommit()
	SendOK(d, "RSS Registered!")
}
