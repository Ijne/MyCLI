package start

import (
	"flag"
)

type Config struct {
	VFSPath    string
	ScriptPath string
}

func parseCommandLineFlags() Config {
	var config Config

	flag.StringVar(&config.VFSPath, "vfs", "", "Путь к файлу конфигурации VFS (CSV)")
	flag.StringVar(&config.VFSPath, "v", "", "Путь к файлу конфигурации VFS (CSV) (короткая версия)")
	flag.StringVar(&config.ScriptPath, "script", "", "Путь к стартовому скрипту для выполнения")
	flag.StringVar(&config.ScriptPath, "s", "", "Путь к стартовому скрипту для выполнения (короткая версия)")

	flag.Parse()

	return config
}
