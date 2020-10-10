package serivcemiddleware

import "github.com/vrazdalovschi/url-shortener/internal/service"

type Middleware func(service.Service) service.Service
