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

	var serr *service.ServiceError
	if errors.As(err, &serr) {
		code = serr.Status
		message = serr.Error()
	} else {
		log.Println(err)
	}

	err = c.JSON(code, service.ErrorResponse{
		Message: message,
	})
}
