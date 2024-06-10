package route

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"golink/frontend/components/link"
	"golink/service"
	"net/http"
)

type LinkHandler struct {
	LinkService service.LinkService
	BaseUrl     string
}

func (h *LinkHandler) CreateLink(ctx echo.Context) error {
	formParams, err := ctx.FormParams()
	if err != nil {
		return err
	}

	code := formParams.Get("code")
	target := formParams.Get("target")

	createdLink, err := h.LinkService.CreateLink(target, &code)
	if err != nil {
		return err
	}
	shortUrl := templ.SafeURL(h.BaseUrl + "/to/" + createdLink.Code)
	return Render(ctx, link.Created(*createdLink, shortUrl))
}

func (h *LinkHandler) RedirectLink(c echo.Context) error {
	code := c.Param("code")

	fetchedLink := h.LinkService.GetLink(code)

	if fetchedLink == nil {
		return Render(c, link.NotFound())
	}

	return c.Redirect(http.StatusTemporaryRedirect, fetchedLink.Target)
}
