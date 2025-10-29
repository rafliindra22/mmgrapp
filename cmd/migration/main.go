package main

import config "mmgrapp/internal/configs"

func main() {
	config.ConnectDB()
	config.Migrate()
}
