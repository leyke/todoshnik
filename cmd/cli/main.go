package main

// импорты нужно разделять по разным соглашениям на 2 или 3 группы:
// - зависимости из стандартной библиотеки
// - зависимости из проекта
// - сторонние зависимости
// Кажется это есть в Убер стайл гайде из ссылок в гуглдоке
import (
	"os"
	"todoshnik/internal/app"
	"todoshnik/internal/cli"
)

// так никто не делает, не нужно указывать тип при инициализации
// тут бы больше подошла константа, потому что ты не ожидаешь что в рантайме значение будет переопределено
//
// еще лучше было бы получать путь в переменных окружения
//
// еще лучше следовать манифесту 12 факторного приложения и писать логи в stdout, stderr
var logFileName string = "/cli.log"

func main() {
	container := app.InitApp(logFileName)
	defer container.LogFile.Close()

	cli := cli.NewHandler(container.TaskService, container.TokenService)
	cli.Run(os.Args)
}
