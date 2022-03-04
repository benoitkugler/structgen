// Package inner is only used to test apigen.
package inner

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

const Url = "/const_url_from_inner_package/"

type Controller struct{}

func (Controller) HandleExt(c echo.Context) error {
	var in []int64
	t, v := c.QueryParam("query1"), c.QueryParam("query2")
	err := c.Bind(&in)
	_ = fmt.Errorf("%s%s%s", t, v, err)
	var out map[string][]int
	return c.JSON(200, out)
}
