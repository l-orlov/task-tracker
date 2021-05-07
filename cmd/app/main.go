package main

import (
	"os"

	"github.com/l-orlov/task-tracker/internal/app"
)

const envConfigPath = "CONFIG_PATH"

func main() {
	app.Run(os.Getenv(envConfigPath))
}
