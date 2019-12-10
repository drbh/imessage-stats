package imessagehooks

import (
	"database/sql"
	"log"
	"os/user"
	"strconv"
	"time"

	"golang.org/x/sync/syncmap"
)

var (
	recentlySeen = syncmap.Map{}
)

func RunPoller(callback func(string)) {
	ticker := time.NewTicker(1000 * time.Millisecond)
	quit := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:

				now := time.Now()
				recentlySeen.Range(func(key interface{}, value interface{}) bool {
					wastime := AppleTimestampToTime(value.(string))
					diff := now.Sub(wastime)
					if diff > 4*time.Second {
						recentlySeen.Delete(key)
					}
					return true
				})

				poll(callback)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

// DateOffset is the time needed to offset an apple timestamp
const DateOffset = int64(978307200)

func MakeAppleTimestamp() int {
	a := time.Now().UnixNano() / int64(time.Millisecond)
	b := a / 1000
	c := b - (DateOffset + 2)
	d := c * 1000 * 1000 * 1000
	return int(d)
}

func AppleTimestampToTime(timestamp string) time.Time {
	n, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
	}
	x := n
	y := x / (1000 * 1000 * 1000)
	z := y + DateOffset
	k := z * 1000
	return time.Unix(0, int64(k)*int64(time.Millisecond))
}

type IMessageRow struct {
	GUID           string `db:"guid"`
	Handle         string `db:"handle"`
	HandleID       string `db:"handle_id"`
	Text           string `db:"text"`
	Date           string `db:"date"`
	DateRead       string `db:"date_read"`
	IsFromMe       string `db:"is_from_me"`
	CacheRoomnames string `db:"cache_roomnames"`
	IsRead         string `db:"is_read"`
}

func Fetch(handle string, latest string) []IMessageRow {
	args := "?mode=ro&_mutex=no&_journal_mode=WAL&_query_only=1&_synchronous=2"
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	connectionString := usr.HomeDir + "/Library/Messages/chat.db" + args
	database, err := sql.Open("sqlite3", connectionString)
	defer database.Close()
	if err != nil {
		log.Fatal("Connection Failed ", err)
	}
	defer database.Close()
	database.SetMaxOpenConns(1)
	rows, qerr := database.Query(`
			 SELECT
			    guid,
			    id as handle,
			    handle_id,
			    text,
			    date,
			    date_read,
			    is_from_me,
				cache_roomnames,
				is_read
			FROM message
			LEFT OUTER JOIN handle ON message.handle_id = handle.ROWID
			WHERE date >= ` + latest + `
			AND handle == "` + handle + `"
        `)
	defer rows.Close()
	if qerr != nil {
		log.Fatal("Query Failed ", qerr)
	}

	var allRows []IMessageRow

	for rows.Next() {
		var imrow IMessageRow
		rows.Scan(&imrow.GUID,
			&imrow.Handle,
			&imrow.HandleID,
			&imrow.Text,
			&imrow.Date,
			&imrow.DateRead,
			&imrow.IsFromMe,
			&imrow.CacheRoomnames,
			&imrow.IsRead)

		allRows = append(allRows, imrow)
	}
	return allRows

}

func FetchFullDatabase(latest string) []IMessageRow {
	args := "?mode=ro&_mutex=no&_journal_mode=WAL&_query_only=1&_synchronous=2"
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	connectionString := usr.HomeDir + "/Library/Messages/chat.db" + args
	database, err := sql.Open("sqlite3", connectionString)
	defer database.Close()
	if err != nil {
		log.Fatal("Connection Failed ", err)
	}
	defer database.Close()
	database.SetMaxOpenConns(1)
	rows, qerr := database.Query(`
			 SELECT
			    guid,
			    id as handle,
			    handle_id,
			    text,
			    date,
			    date_read,
			    is_from_me,
				cache_roomnames,
				is_read
			FROM message
			LEFT OUTER JOIN handle ON message.handle_id = handle.ROWID
			WHERE date >= ` + latest + `
        `)
	defer rows.Close()
	if qerr != nil {
		log.Fatal("Query Failed ", qerr)
	}

	var allRows []IMessageRow

	for rows.Next() {
		var imrow IMessageRow
		rows.Scan(&imrow.GUID,
			&imrow.Handle,
			&imrow.HandleID,
			&imrow.Text,
			&imrow.Date,
			&imrow.DateRead,
			&imrow.IsFromMe,
			&imrow.CacheRoomnames,
			&imrow.IsRead)

		allRows = append(allRows, imrow)
	}
	return allRows

}

func poll(callback func(string)) {
	args := "?mode=ro&_mutex=no&_journal_mode=WAL&_query_only=1&_synchronous=2"
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	connectionString := usr.HomeDir + "/Library/Messages/chat.db" + args
	database, err := sql.Open("sqlite3", connectionString)
	defer database.Close()

	if err != nil {
		log.Fatal("Connection Failed ", err)
	}
	defer database.Close()

	database.SetMaxOpenConns(1)
	latestHolder := MakeAppleTimestamp()
	latest := strconv.Itoa(latestHolder - 1)

	rows, qerr := database.Query(`
			 SELECT
			    guid,
			    id as handle,
			    handle_id,
			    text,
			    date,
			    date_read,
			    is_from_me,
				cache_roomnames,
				is_read
			FROM message
			LEFT OUTER JOIN handle ON message.handle_id = handle.ROWID
			WHERE date >= ` + latest + `
        `)
	defer rows.Close()
	if qerr != nil {
		log.Fatal("Query Failed ", qerr)
	}

	var guid string
	var handle string
	var handleID string
	var text string
	var date string
	var dateRead string
	var isFromMe string
	var cacheRoomnames string
	var isRead string

	for rows.Next() {
		rows.Scan(
			&guid,
			&handle,
			&handleID,
			&text,
			&date,
			&dateRead,
			&isFromMe,
			&cacheRoomnames,
			&isRead,
		)

		_, ok := recentlySeen.Load(guid)
		if ok {
			// fmt.Println("found - dont send again")
			return
		}
		// else we set the map values
		recentlySeen.Store(guid, date)
		toSend := guid + " " + handle + " " + " " + text
		callback(toSend)
	}
}

type FullMessageRow struct {
	ROWID                          int         `json:"ROWID"`
	GUID                           string      `json:"guid"`
	Text                           interface{} `json:"text"`
	Replace                        int         `json:"replace"`
	ServiceCenter                  interface{} `json:"service_center"`
	HandleID                       int         `json:"handle_id"`
	Subject                        interface{} `json:"subject"`
	Country                        interface{} `json:"country"`
	AttributedBody                 interface{} `json:"attributedBody"`
	Version                        int         `json:"version"`
	Type                           int         `json:"type"`
	Service                        string      `json:"service"`
	Account                        string      `json:"account"`
	AccountGUID                    string      `json:"account_guid"`
	Error                          int         `json:"error"`
	Date                           int64       `json:"date"`
	DateRead                       int         `json:"date_read"`
	DateDelivered                  int         `json:"date_delivered"`
	IsDelivered                    int         `json:"is_delivered"`
	IsFinished                     int         `json:"is_finished"`
	IsEmote                        int         `json:"is_emote"`
	IsFromMe                       int         `json:"is_from_me"`
	IsEmpty                        int         `json:"is_empty"`
	IsDelayed                      int         `json:"is_delayed"`
	IsAutoReply                    int         `json:"is_auto_reply"`
	IsPrepared                     int         `json:"is_prepared"`
	IsRead                         int         `json:"is_read"`
	IsSystemMessage                int         `json:"is_system_message"`
	IsSent                         int         `json:"is_sent"`
	HasDdResults                   int         `json:"has_dd_results"`
	IsServiceMessage               int         `json:"is_service_message"`
	IsForward                      int         `json:"is_forward"`
	WasDowngraded                  int         `json:"was_downgraded"`
	IsArchive                      int         `json:"is_archive"`
	CacheHasAttachments            int         `json:"cache_has_attachments"`
	CacheRoomnames                 interface{} `json:"cache_roomnames"`
	WasDataDetected                int         `json:"was_data_detected"`
	WasDeduplicated                int         `json:"was_deduplicated"`
	IsAudioMessage                 int         `json:"is_audio_message"`
	IsPlayed                       int         `json:"is_played"`
	DatePlayed                     int         `json:"date_played"`
	ItemType                       int         `json:"item_type"`
	OtherHandle                    int         `json:"other_handle"`
	GroupTitle                     interface{} `json:"group_title"`
	GroupActionType                int         `json:"group_action_type"`
	ShareStatus                    int         `json:"share_status"`
	ShareDirection                 int         `json:"share_direction"`
	IsExpirable                    int         `json:"is_expirable"`
	ExpireState                    int         `json:"expire_state"`
	MessageActionType              int         `json:"message_action_type"`
	MessageSource                  int         `json:"message_source"`
	AssociatedMessageGUID          interface{} `json:"associated_message_guid"`
	AssociatedMessageType          int         `json:"associated_message_type"`
	BalloonBundleID                interface{} `json:"balloon_bundle_id"`
	PayloadData                    interface{} `json:"payload_data"`
	ExpressiveSendStyleID          interface{} `json:"expressive_send_style_id"`
	AssociatedMessageRangeLocation int         `json:"associated_message_range_location"`
	AssociatedMessageRangeLength   int         `json:"associated_message_range_length"`
	TimeExpressiveSendPlayed       int         `json:"time_expressive_send_played"`
	MessageSummaryInfo             interface{} `json:"message_summary_info"`
	CkSyncState                    int         `json:"ck_sync_state"`
	CkRecordID                     interface{} `json:"ck_record_id"`
	CkRecordChangeTag              interface{} `json:"ck_record_change_tag"`
	DestinationCallerID            interface{} `json:"destination_caller_id"`
	SrCkSyncState                  int         `json:"sr_ck_sync_state"`
	SrCkRecordID                   interface{} `json:"sr_ck_record_id"`
	SrCkRecordChangeTag            interface{} `json:"sr_ck_record_change_tag"`
	IsCorrupt                      int         `json:"is_corrupt"`
}
