package web

import (
	"log"

	"github.com/a-h/templ"
	"github.com/boundedinfinity/statementer/runtime"
	"github.com/gofiber/fiber/v2"
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

var (
	_PREFIX_PROCESSED_DIR = "/processed-dir"
	_PREFIX_SOURCE_DIR    = "/source-dir"
)

func (this *Web) Init() error {
	this.fiber = fiber.New()
	this.fiber.Use(logger.New())

	this.fiber.Get("/", func(c *fiber.Ctx) error {
		return Render(c, home(this.runtime.Config))
	})

	this.fiber.Get("/files/list", func(c *fiber.Ctx) error {
		return Render(c, filesList(this.runtime.FilesAll()))
	})

	this.fiber.Get("/files/duplicates", func(c *fiber.Ctx) error {
		return Render(c, filesDuplicates(this.runtime.FilesDuplicates()))
	})

	this.fiber.Get("/files/merge", func(c *fiber.Ctx) error {
		return nil
	})

	this.fiber.Get("/labels/all", func(c *fiber.Ctx) error {
		return Render(c, simpleLabelsList(this.runtime.Labels.All()))
	})

	this.initFileDetails()
	this.initUtils()

	return nil
}

func (this *Web) initUtils() error {
	this.fiber.Get("/open/config-file", func(c *fiber.Ctx) error {
		text, err := this.runtime.OpenConfigFile()
		return Render(c, message(text, err))
	})

	this.fiber.Get("/open/source-dir", func(c *fiber.Ctx) error {
		text, err := this.runtime.OpenSourceDir()
		return Render(c, message(text, err))
	})

	this.fiber.Get("/open/repository-dir", func(c *fiber.Ctx) error {
		text, err := this.runtime.OpenRepositoryDir()
		return Render(c, message(text, err))
	})

	this.fiber.Get("/open/document/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		files := this.runtime.FileGet(id)
		if len(files) > 0 {
			if len(files[0].SourcePaths) > 0 {
				path := webSourcePath(this.runtime.Config, files[0].SourcePaths[0])
				return Render(c, documentViewer(path))
			}
		}

		return nil
	})

	this.fiber.Static("/", "./assets")
	this.fiber.Static(_PREFIX_PROCESSED_DIR, this.runtime.Config.RepositoryDir, fiber.Static{
		Browse: true,
	})
	this.fiber.Static(_PREFIX_SOURCE_DIR, this.runtime.Config.SourceDir, fiber.Static{
		Browse: true,
	})
	return nil
}

func (this *Web) initFileDetails() error {
	this.fiber.Get("/files/details/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		files := this.runtime.FileGet(id)

		if len(files) == 1 {
			return Render(c, filesDetails(files[0]))
		}

		return nil
	})

	this.fiber.Get("/files/title/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		files := this.runtime.State.Files.ById(id)
		return Render(c, fileViewTitle(files[0]))
	})

	this.fiber.Patch("/files/title/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		files := this.runtime.State.Files.ById(id)
		return Render(c, fileEditTitle(files[0]))
	})

	this.fiber.Post("/files/title/:id", func(c *fiber.Ctx) error {
		id := c.FormValue("id")
		title := c.FormValue("title")
		files := this.runtime.State.Files.ById(id)
		files[0].Title = title

		if err := this.runtime.SaveState(); err != nil {
			log.Println(err.Error())
		}

		c.Response().Header.Add("HX-Trigger", "file-updated")
		return Render(c, fileViewTitle(files[0]))
	})

	return nil
}

// https://github.com/a-h/templ/blob/main/examples/integration-gofiber/main.go/

func Render(c *fiber.Ctx, component templ.Component) error {
	c.Set("Content-Type", "text/html")
	return component.Render(c.Context(), c.Response().BodyWriter())
}
