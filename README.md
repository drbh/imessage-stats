# Welcome to Imessage Stats

This is a command line tool that allows you to get statistics from your Messages data! You'll need Go installed and the following dependencies

```bash
go get github.com/cdipaolo/sentiment
go get github.com/drbh/gomoji-counter
go get github.com/mattn/go-sqlite3
go get golang.org/x/sync/syncmap
```

also i'd suggest installing `jq` so you can manage the large JSON repsonses in the Terminal. `| jq '.' >` sends the program output to a file, while pretty printing the data.

#### Apple Security Settings Setup
```bash
open -b com.apple.systempreferences /System/Library/PreferencePanes/Security.prefPane
```

Now you'll want to allow Terminal to have `Full Disk Access` this is needed to read the Message db from the Terminal. 


# Get Some Stats

## Per Number
### Get Profile Stats
```bash
go run main.go +12223334444 | jq '.' > example.json
```

## Per Database
### Get Profile Stats
```bash
go run main.go --all | jq '.' > example.json
```


### Get Word Frequencies
```bash
go run main.go counts | jq '.' > example.json
```


# What Stats Are Included

### Get Profile Stats 
- EmojiMap 
- SentimentScore 
- MessageCount 
- FirstSeen 
- AverageResponseSeconds 
  
#### Stats Per Message 
- Year 
- Month 
- Day 
- Wkday 
- Hour 
- Len 
- Positve 
- Timestamp 

#### Frequency Counts Per 
- Weekday By HourOfDay. 


### Gotchas

If a number is in more then one chat, like a direct message and a group chat; there chat stats will not seperate those conversation from the statistics (so if the dates look like they dont match up - find the group chat!)

