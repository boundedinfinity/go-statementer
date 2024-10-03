package web

import (
	"github.com/a-h/templ"
	"github.com/boundedinfinity/statementer/runtime"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func New(runtime *runtime.Runtime) *Web {
	web := &Web{
		runtime: runtime,
	}

	return web
}

type Web struct {
	runtime *runtime.Runtime
	fiber   *fiber.App
}

func (this *Web) Listen() error {
	return this.fiber.Listen(":3000")
}

func (this *Web) Init() error {
	this.fiber = fiber.New()
	this.fiber.Use(logger.New())
	this.fiber.Static("/", "./assets")

	this.fiber.Get("/", this.handleHome())
	this.fiber.Get("/files/list", this.handleFilesList())
	this.fiber.Get("/files/duplicates", this.handleFilesDuplistList())
	this.fiber.Get("/files/merge", this.handleFilesDuplistList())

	return nil
}

func (this *Web) handler(component templ.Component, options ...func(*templ.ComponentHandler)) func(*fiber.Ctx) error {
	return adaptor.HTTPHandler(templ.Handler(component, options...))
}

func (this *Web) handleHome() func(*fiber.Ctx) error {
	return this.handler(home())
}

func (this *Web) handleFilesList() func(*fiber.Ctx) error {
	return this.handler(filesList(this.runtime.FilesAll()))
}

func (this *Web) handleFilesDuplistList() func(*fiber.Ctx) error {
	return this.handler(filesDuplicates(this.runtime.FilesDuplicates()))
}
