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
4. Create a config file named `config.json` in your home directory with the following content:
```
{
  "database_url": "postgresql://username:password@localhost:5432/dbname"
}
```

After installation, you can run gator from any directory. Simply use `gator <command> [args...]` in terminal.
