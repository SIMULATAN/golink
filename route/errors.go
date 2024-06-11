package route

import (
	"errors"
	"github.com/labstack/echo/v4"
	"golink/service"
	"log"
	"net/http"
)

func HandleError(app *echo.Echo, err error, c echo.Context) {
	code := http.StatusInternalServerError
	message := "Internal Server Error!"
	if app.Debug {
		message = err.Error()
	}

	var herr *echo.HTTPError
	if errors.As(err, &herr) {
		code = herr.Code
		msg := herr.Message.(string)
		if msg != "" {
			message = msg
		} else {
			message = herr.Error()
		}
		err = c.JSON(code, service.ErrorResponse{
			Message: message,
		})
		return
	}

	var serr *service.ServiceError
	if errors.As(err, &serr) {
		code = serr.Status
		message = serr.Error()
	} else {
		log.Println("Unhandled error", err)
	}

	err = c.JSON(code, service.ErrorResponse{
		Message: message,
	})
}
