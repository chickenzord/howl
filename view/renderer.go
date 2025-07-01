package howlview

import (
	"fmt"
	"io"
	"io/fs"

	"github.com/flosch/pongo2/v6"
	"github.com/labstack/echo/v4"
)

type Renderer struct {
	FS   fs.FS
	Data Data
}

func (r *Renderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	s := pongo2.NewSet("local", pongo2.NewFSLoader(r.FS))
	tpl, err := s.FromFile(name)
	if err != nil {
		return err
	}

	ctx := pongo2.Context{}
	ctx["url_for"] = c.Echo().Reverse
	for key, val := range r.Data {
		ctx[key] = val
	}

	switch v := data.(type) {
	case Data:
		for key, val := range v {
			ctx[key] = val
		}
	case pongo2.Context:
		for key, val := range v {
			ctx[key] = val
		}
	case map[string]interface{}:
		for key, val := range v {
			ctx[key] = val
		}
	default:
		return fmt.Errorf("invalid templating data type: %v", v)
	}

	return tpl.ExecuteWriter(ctx, w)
}
