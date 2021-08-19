package app

import "github.com/MrHanson/gin-blog/pkg/logging"

func MarkErrors(errors []error) {
	logging.Info(errors)
}
