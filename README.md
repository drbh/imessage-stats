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

### Command List
- number [phone number]
- all 
- counts

## Per Number
### Get Profile Stats
```bash
go run main.go number +12223334444 | jq '.' > example.json
```

## Per Database
### Get Profile Stats
```bash
go run main.go all | jq '.' > example.json
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


## Example Data

### Get Profile Stats
```JSON
{
  "Messages": [
    {
      "Year": 2018,
      "Month": 8,
      "Day": 12,
      "Wkday": "Sunday",
      "Hour": 20,
      "Len": 12,
      "Positve": 0,
      "Timestamp": "555813885000000000"
    },
  ],
  "EmojiMap": {
    "“": 1,
    "”": 1,
    "❤": 5,
    "🌲": 1,
    "🎄": 1,
    "👍": 1,
    "💕": 3,
    "😂": 1,
    "😉": 1,
    "😊": 5,
    "😘": 15,
    "🙏": 3,
    "🤗": 1,
    "🤣": 1
  },
  "WkHr": {
    "Friday_0": 0,
    "Friday_1": 0,
    "Friday_10": 0,
    "Friday_11": 0,
  },
  "SentimentScore": 1,
  "MessageCount": 75,
  "FirstSeen": "2000-12-31T19:00:00-05:00",
  "AverageResponseSeconds": 1649333,
  "ResponseTimes": [
    {
      "IsentTime": "2018-09-20T10:33:41-04:00",
      "TheyRespondTime": "2018-08-12T20:44:57-04:00",
      "Diff": 3332924000000000
    },
  ]
}
```

### Get Word Frequencies

```JSON
{
  "!": 14,
  "!!": 9,
  "!!!": 13,
  "!!!!": 3,
  "!!!!!": 7,
  "!!!!!!": 1,
  "!!!!!!!": 3,
  "!!!!!!!!!!!!!!!!!!!!": 1,
  "!!!!!!!!!!!!!!!!!!!!!!!!!!": 1,
  "!!***": 1,
  "!**": 2,
  "!***": 1,
  "!=": 3,
  "!?": 1,
  "!countdown": 1,
  "\"": 7,
}
```


### Gotchas

If a number is in more then one chat, like a direct message and a group chat; there chat stats will not seperate those conversation from the statistics (so if the dates look like they dont match up - find the group chat!)

