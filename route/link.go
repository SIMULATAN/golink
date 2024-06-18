package route

import (
	"errors"
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"golink/frontend/components/link"
	"golink/service"
	"golink/utils"
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

	if code == "" {
		code = utils.RandSeq(32)
	}

	createdLink, err := h.LinkService.CreateLink(target, code)
	if err != nil {
		if errors.Is(err, service.CodeExistsError) {
			return Render(ctx, link.CodeExists())
		}
		return err
	}

	fullUrl := h.BaseUrl + "/to/" + createdLink.Code
	safeUrl := templ.SafeURL(fullUrl)
	return Render(ctx, link.Created(fullUrl, safeUrl))
}

func (h *LinkHandler) RedirectLink(c echo.Context) error {
	code := c.Param("code")

	fetchedLink, err := h.LinkService.GetLink(code)
	if err != nil {
		return err
	}

	if fetchedLink == nil {
		return Render(c, link.NotFound())
	}

	return c.Redirect(http.StatusTemporaryRedirect, fetchedLink.Target)
}
