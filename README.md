# Quake Parser
 Parser for Quake server Log

## Usage
```bash
git clone https://github.com/luanrivello/quake_parser
cd quake_parser
```
- You can run the programm with no arguments defaulting the log file to ./data/qgames.log
```bash
go run main.go
```
- or you can give it a path like so
```bash
go run main.go $LOGPATH
```
- and a report will be generated in json format at the ./report folder