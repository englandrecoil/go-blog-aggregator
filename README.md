# go-blog-aggregator

**Gator CLI** lets you add RSS feeds from across the internet, store posts in a PostgreSQL database, follow and unfollow from the other user's RSS feeds and view summaries of aggregated posts in terminal.

## 🚀 Prerequisites
You'll need **two** things installed to run the program:
1. Go 1.23 or later ([check official instructions](https://go.dev))
2. PostgreSQL 16 or later ([check how to install here](https://www.postgresql.org))

Run `go version` and `postgres --version` in your terminal to make sure the installation of Go and PostgreSQL worked.

## :minidisc: Installation
To use, follow these steps:
1. Clone the repository with `git clone https://github.com/englandrecoil/go-blog-aggregator`
2. Move to the created folder and build the project with `go build -o gator`
3. Create a symlink to make it available system-wide with `sudo ln -s $(pwd)/gator /usr/local/bin/gator`

## :spiral_notepad: Config
Create a config file named `config.json` in your home directory with the following content:
```json
{
  "db_url": "postgres://username:@localhost:5432/database?sslmode=disable"
}
```
Make sure to replace the values _username, password, database_ with your database connection string


## :keyboard: Usage
After installation, you can run gator from any directory. Simply use `gator <command> [args...]` in terminal. But first you need to:

Create a new User:
```bash
gator register <name>
```

Add a feed:
```bash
gator addfeed <url>
```

Start aggregator:
```bash
gator agg 30s
```

View available posts:
```bash
gator browse <limit>
```

Other available commands can be found below:
| Command  | Description |
| ------------- | ------------- |
| `gator login <name>`  | Log in as a user that already exists |
| `gator register <name>`  | Register a new user with given name |
| `gator feeds` | Displays all feeds |
| `gator users` | Displays all users |
| `gator follow <url>` | Follow an existing feed |
| `gator unfollow <url>` | Unfollow an existing feed |
| `gator reset` | Delete all users permanently |
