package noob

import _logger "github.com/alfarih31/nb-go-logger"

var Log = _logger.New("core")

var logR = Log.NewChild("response")
