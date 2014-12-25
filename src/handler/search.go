package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/ysugimoto/husky"
	"strconv"
	"strings"
)

type SearchResult struct {
	Id    int    `json:"id"`
	Url   string `json:"url"`
	Title string `json:"title"`
	Tag   string `json:"tag"`
}

func Search(d *husky.Dispatcher) {
	db := husky.NewDb(GetDSN())
	req := d.Input.GetRequest()

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
		message := "Token not matched!"
		SendError(d, message)
		return
	}

	var qs []string
	var limit int

	param := req.URL.Query()
	if get, ok := param["q"]; !ok {
		qs = append(qs, "")
	} else {
		qs = get
	}
	if get, ok := param["l"]; !ok {
		limit = 10
	} else {
		if parse, err := strconv.Atoi(get[0]); err != nil {
			limit = parse
		} else {
			limit = 10
		}
	}

	// trim space and split search query
	q := string(bytes.TrimSpace([]byte(qs[0])))
	qq := strings.Split(q, " ")

	// Search Query
	query := "SELECT U.id, U.url, U.title, T.name FROM pb_tags as T "
	query += "JOIN pb_urls as U ON ( T.url_id = U.id ) "
	query += "WHERE U.user_id = ? AND "
	bind := []interface{}{userId}
	where := []string{}

	for _, l := range qq {
		where = append(where, "T.name LIKE ?")
		bind = append(bind, "%"+string(l)+"%")
	}
	query += strings.Join(where, " OR ")
	query += " LIMIT " + fmt.Sprint(limit)

	rows, err := db.Query(query, bind...)
	if err != nil {
		message := fmt.Sprintf("Query Error: %v", err)
		SendError(d, message)
		fmt.Printf("%v\n", err)
		return
	}

	var result []SearchResult
	for rows.Next() {
		r := SearchResult{}
		rows.Scan(&r.Id, &r.Url, &r.Title, &r.Tag)
		result = append(result, r)
	}

	result = filterResult(result, qq)

	if encode, err := json.Marshal(result); err != nil {
		SendError(d, fmt.Sprintf("Endode error: %v", err))
		return
	} else {
		d.Output.SetHeader("Content-Type", "application/json")
		d.Output.SetHeader("Access-Control-Allow-Origin", "*")
		d.Output.SetHeader("Access-Control-Allow-Headers", "X-LAP-Token")
		d.Output.SetStatus(200)
		d.Output.Send(encode)
	}
}

type FilterStack struct {
	hit int
	row SearchResult
}

func filterResult(result []SearchResult, tags []string) (filtered []SearchResult) {
	stack := map[int]FilterStack{}

	// count up tags hit
	for _, row := range result {
		if tmp, exists := stack[row.Id]; !exists {
			stack[row.Id] = FilterStack{
				hit: 1,
				row: row,
			}
		} else {
			tmp.hit++
			stack[row.Id] = tmp
		}
	}

	// Factory: hit times equals tags count
	size := len(tags)
	for _, filter := range stack {
		if filter.hit == size {
			filtered = append(filtered, filter.row)
		}
	}

	return
}
