package source

import (
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"github.com/nicklaw5/helix"
	"log"
)

type Twitch struct {
	client *helix.Client
	db     *sql.DB

	usernameToIdMap map[string]string
}

func (t *Twitch) initCache() {
	createStudentTableSQL := `
		CREATE TABLE IF NOT EXISTS twitch_users (
			user_id TEXT UNIQUE,
			username TEXT
		);
	`

	_, err := t.db.Exec(createStudentTableSQL) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
}

func (t *Twitch) SaveActorName(userId string, name string) {
	query := `
		INSERT IGNORE INTO twitch_users(user_id, username) VALUES (?, ?)
		ON CONFLICT(user_id) DO UPDATE SET username = excluded.username;
	`
	_, err := t.db.Exec(query, userId, name)
	if err != nil {
		log.Fatal(err)
	}
}


func (t *Twitch) UserId(username string) (string, error) {
	if id, ok := t.usernameToIdMap[username]; ok {
		return id, nil
	}

	query := `
		SELECT user_id FROM twitch_users WHERE username = ? LIMIT 1
	`
	row, err := t.db.Query(query, username)
	if err != nil {
		return "", err
	}
	defer row.Close()
	for row.Next() { // Iterate and fetch the records from result cursor
		var userId string
		row.Scan(&userId)
		t.usernameToIdMap[username] = userId
		return userId, nil
	}

	users, err := t.client.GetUsers(&helix.UsersParams{
		Logins: []string{username},
	})
	if err != nil {
		return "", err
	}

	for _, user := range users.Data.Users {
		t.usernameToIdMap[username] = user.ID
		return user.ID, nil
	}

	return "", errors.New("unknown user")
}

func (t *Twitch) Videos(username string) []helix.Video {
	userId, err := t.UserId(username)
	if err != nil {
		panic(err)
	}

	videos, err := t.client.GetVideos(&helix.VideosParams{
		UserID: userId,
	})

	if err != nil {
		panic(err)
	}

	urls := make([]helix.Video, 0, 25)
	for _, v := range videos.Data.Videos {
		if v.Type != "archive" {
			continue
		}
		urls = append(urls, v)
	}

	return urls
}

func NewTwitch(clientId string, clientSecret string, db *sql.DB) *Twitch {
	client, err := helix.NewClient(&helix.Options{
		ClientID:     clientId,
		ClientSecret: clientSecret,
	})
	if err != nil {
		panic(err)
	}

	resp, err := client.RequestAppAccessToken([]string{"user:read:email"})
	if err != nil {
		panic(err)
	}

	// Set the access token on the client
	client.SetAppAccessToken(resp.Data.AccessToken)

	t := &Twitch{
		client:          client,
		db:              db,
		usernameToIdMap: make(map[string]string, 800),
	}
	t.initCache()

	return t
}
