package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func randToken(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}

func main() {
	dsn := flag.String("dsn", "root:@tcp(127.0.0.1:3306)/teamgram?charset=utf8mb4&parseTime=true&loc=UTC", "mysql dsn")
	creator := flag.Int64("creator", 777000, "creator user id")
	botID := flag.Int64("bot_id", 0, "bot user id (required)")
	username := flag.String("username", "", "bot username (required)")
	firstName := flag.String("first_name", "NewBot", "bot first name")
	token := flag.String("token", "", "bot token (optional, auto-generate if empty)")
	flag.Parse()

	if *botID <= 0 || strings.TrimSpace(*username) == "" {
		log.Fatal("bot_id and username are required")
	}
	if *token == "" {
		*token = fmt.Sprintf("%d:%s", *botID, randToken(16))
	}

	db, err := sql.Open("mysql", *dsn)
	if err != nil { log.Fatal(err) }
	defer db.Close()

	now := time.Now()
	tx, err := db.Begin()
	if err != nil { log.Fatal(err) }
	defer tx.Rollback()

	_, err = tx.Exec(`INSERT INTO users (id, user_type, access_hash, secret_key_id, first_name, last_name, username, phone, country_code, verified, support, scam, fake, premium, about, state, is_bot, account_days_ttl, photo_id, restricted, restriction_reason, archive_and_mute_new_noncontact_peers, emoji_status_document_id, emoji_status_until, deleted, delete_reason, created_at, updated_at)
VALUES (?, 4, 0, 0, ?, '', ?, '', '', 1, 0, 0, 0, 0, '', 0, 1, 180, 0, 0, '', 0, 0, 0, 0, '', ?, ?)
ON DUPLICATE KEY UPDATE first_name=VALUES(first_name), username=VALUES(username), is_bot=1, updated_at=VALUES(updated_at)`, *botID, *firstName, *username, now, now)
	if err != nil { log.Fatal(err) }

	_, err = tx.Exec(`INSERT INTO bots (bot_id, bot_type, creator_user_id, token, description, bot_chat_history, bot_nochats, bot_inline_geo, bot_info_version, bot_inline_placeholder, attach_menu_enabled, bot_attach_menu, bot_business, bot_has_main_app, bot_active_users, has_menu_button, menu_button_text, menu_button_url, bot_can_edit, has_preview_medias, description_photo_id, description_document_id, main_app_url, has_app_settings, placeholder_path, background_color, background_dark_color, header_color, header_dark_color, privacy_policy_url)
VALUES (?, 0, ?, ?, '', 0, 0, 0, 1, '', 0, 0, 0, 0, 0, 0, '', '', 0, 0, 0, 0, '', 0, '', 0, 0, 0, 0, '')
ON DUPLICATE KEY UPDATE creator_user_id=VALUES(creator_user_id), token=VALUES(token)`, *botID, *creator, *token)
	if err != nil { log.Fatal(err) }

	if err := tx.Commit(); err != nil { log.Fatal(err) }
	fmt.Printf("BOT_CREATED id=%d username=%s token=%s\n", *botID, *username, *token)
}
