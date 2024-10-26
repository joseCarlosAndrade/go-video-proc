package converter

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"time"
)

// returns if the video was already processed
func IsProcessed(db *sql.DB, videoID int) bool {
	var isProcessed bool

	query := "SELECT EXISTS(SELECT 1 FROM processed_videos WHERE video_id = $1 AND status='success')"

	// query row queries the db and returns rows
	err := db.QueryRow(query, videoID).Scan(&isProcessed) // executes the query and then scans the result into the isProcessed var

	if err != nil {
		slog.Error("Error checking video processing", slog.Int("video_id", videoID), slog.String("error", err.Error()))
		return false
	}

	return isProcessed
}

// registers in the db that this video was processed
func MarkProcessed(db *sql.DB, videoId int) error {
	query := "INSERT INTO processed_videos (video_id, status, processed_at) values ($1, $2, $3)"

	_, err := db.Exec(query, videoId, "success", time.Now()) // exec statements without returning any rows

	if err != nil {
		slog.Error("error registering video as processed", slog.Int("video_id", videoId), slog.String("error", err.Error()))
		return err
	}

	return nil
}

func RegisterError(db *sql.DB, errorData map[string]any, err error) {
	serializedError, _ := json.Marshal(errorData)

	query := "INSERT INTO process_errors_log (error_details, created_at) VALUES ($1, $2)"

	_, dbErr := db.Exec(query, serializedError)
	if dbErr != nil {
		slog.Error("error registering error", slog.String("error_details", string(serializedError)))
		return 
	}
	
}