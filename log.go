package noob

import _logger "github.com/alfarih31/nb-go-logger"

var log = _logger.New("core")

var logR = log.NewChild("response")

func restartLogger() {
	log = _logger.New("core")

	logR = log.NewChild("response")
}
