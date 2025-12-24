package handlers

import (
	"github.com/dineshengg/matrimony/common/utils"
)

func SomeHandler() {
	utils.Logger.WithField("user_id", 123).Info("User handler called")
	utils.LogTracing("some-guid", "User handler tracing event")
}
