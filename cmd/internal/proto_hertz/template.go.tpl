const (
{{range .Methods}}{{.MdMethodKey}} string = "/{{$.FullName}}/{{.HandlerName}}"
{{end}}
)

type {{$.InterfaceName}} interface {
{{range .MethodSet}}{{.Name}}(context.Context, *{{.Request}}) (*{{.Reply}}, error)
{{end}}
}

func Register{{$.InterfaceName}}(r route.IRouter, srv {{$.InterfaceName}}, respEncoders map[string]ResponseHandler, midlleware ...app.HandlerFunc) *{{.NameHttp}} {
	s := &{{$.NameHttp}} {
		router: r,
		server: srv,
		resp:   make(map[string]ResponseHandler),
	}
    for k, v := range respEncoders {
        s.resp[k] = v
    }
    _, exist := s.resp["*"]
    if len(s.resp) == 0 || !exist {
        s.resp["*"] = &defaultResponseImpl{}
    }
	s.initCustomMiddlerware()
	s.RegisterService(midlleware...)
	return s
}

type {{$.NameHttp}} struct {
	router      route.IRouter
	server      {{$.InterfaceName}}
	middlerware map[string][]app.HandlerFunc
    resp        map[string]ResponseHandler
}


func (s *{{$.NameHttp}}) initCustomMiddlerware() {
	s.middlerware = map[string][]app.HandlerFunc{
{{range .Methods}}{{.MdMethodKey}}: {},
{{end}}
	}
}

func (s *{{$.NameHttp}}) CustomMiddlerware(mid map[string][]app.HandlerFunc) {
	for k := range s.middlerware {
		if v, ok := mid[k]; ok {
			s.middlerware[k] = append(s.middlerware[k], v...)
		}
	}
}

func (s *{{.NameHttp}}) RegisterService(middlerware ...app.HandlerFunc) {
	r := s.router.Group("/", middlerware...)
{{range .Methods}}r.Handle("{{.Method}}", "{{.Path}}", append(s.middlerware[{{.MdMethodKey}}], s.{{.HandlerName}})...)
{{end}}
}

{{range .Methods}}
func (s *{{$.NameHttp}}) {{.HandlerName}}(ctx context.Context, c *app.RequestContext) {
	var req {{.Request}}
	err := c.BindAndValidate(&req)
	if err != nil {
        s.getRespEncoder(c).Encode(c, nil, err)
		return
	}
	resp, err := s.server.{{.Name}}(ctx, &req)
    s.getRespEncoder(c).Encode(c, resp, err)
    return
}
{{end}}

func (s *{{$.NameHttp}}) getRespEncoder(c *app.RequestContext) ResponseHandler {
    path := c.FullPath()
    encoder, ok := s.resp[path]
    if !ok {
        encoder = s.resp["*"]
    }
    return encoder
}

type ResponseHandler interface {
	Encode(ctx *app.RequestContext, data interface{}, err error)
}

type defaultResponseImpl struct{}

func (d *defaultResponseImpl) Encode(ctx *app.RequestContext, data interface{}, err error) {
	if err != nil {
		ctx.String(400, err.Error())
	} else {
		ctx.JSON(200, data)
	}
}

