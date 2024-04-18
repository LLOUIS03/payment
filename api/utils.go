package api

import (
	echoLog "github.com/labstack/gommon/log"
)

func lvl(l string) echoLog.Lvl {
	switch l {
	case "DEBUG":
		return echoLog.DEBUG
	case "WARN":
		return echoLog.WARN
	case "ERROR":
		return echoLog.ERROR
	case "OFF":
		return echoLog.OFF
	default:
		return echoLog.INFO
	}
}
