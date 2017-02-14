package log

import (
	"os"

	"github.com/Sirupsen/logrus"
)

var Logger *logrus.Logger

func main() {
	Logger = logrus.New()
	Logger.Out = os.Stdout

	// You could set this to any `io.Writer` such as a file
	// file, err := os.OpenFile("logrus.log", os.O_CREATE|os.O_WRONLY, 0666)
	// if err == nil {
	//  log.Out = file
	// } else {
	//  log.Info("Failed to log to file, using default stderr")
	// }
}
