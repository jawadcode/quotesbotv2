# QuotesBot V2

## What happened to V1?
We do not speak of V1

## What is this?
QuotesBot is a bot that allows you to save and recall quotes (messages) ~~to blackmail your friends~~ for future use.

## Features:
- [x] Help
- [x] Save quote
- [x] Get recent quotes
- [x] Get quote by ID
- [x] Search through quotes (Full-Text search)
- [ ] Delete quote

## Made With:
- Golang
  - [discordgo](https://github.com/bwmarrin/discordgo) - A golang wrapper for the discord API
  - [sqlx](https://github.com/jmoiron/sqlx) - An extension of `database/sql`
- PostgreSQL - ***The superior database***
  - Full Text search, i.e. `to_tsvector()` and `to_tsquery()`

## Todo:
- Separate the handlers into separate files
- Create the rest of the basic handlers
