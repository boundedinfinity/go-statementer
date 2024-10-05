package web

import (
	"log"

	"github.com/a-h/templ"
	"github.com/boundedinfinity/go-commoner/idiomatic/stringer"
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

	this.fiber.Get("/", this.handleHome())

	this.fiber.Get("/files/list", this.handleFilesList())
	this.fiber.Get("/files/duplicates", this.handleFilesDuplistList())
	this.fiber.Get("/files/details/:id", this.handleFilesDetails())
	this.fiber.Patch("/files/details/:id", this.handleFilesUpdate())
	this.fiber.Get("/files/merge", this.handleFilesDuplistList())

	this.fiber.Get("/labels/all", this.handleLabelsAll())

	this.fiber.Get("/open/processed-dir", this.handleOpenProcessedDir())
	this.fiber.Get("/open/source-dir", this.handleOpenSourceDir())
	this.fiber.Get("/open/document/:id", this.handleOpenDocument())

	this.fiber.Static("/", "./assets")
	this.fiber.Static(_PREFIX_PROCESSED_DIR, this.runtime.Config.ProcessedDir, fiber.Static{
		Browse: true,
	})
	this.fiber.Static(_PREFIX_SOURCE_DIR, this.runtime.Config.SourceDir, fiber.Static{
		Browse: true,
	})

	return nil
}

// https://github.com/a-h/templ/blob/main/examples/integration-gofiber/main.go/

func Render(c *fiber.Ctx, component templ.Component) error {
	c.Set("Content-Type", "text/html")
	return component.Render(c.Context(), c.Response().BodyWriter())
}

func (this *Web) handleHome() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return Render(c, home())
	}
}

func (this *Web) handleOpenDocument() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		files := this.runtime.FileGet(id)
		if len(files) > 0 {
			if len(files[0].SourcePaths) > 0 {
				path := files[0].SourcePaths[0]
				path = stringer.Replace(path, _PREFIX_PROCESSED_DIR, this.runtime.Config.ProcessedDir)
				path = stringer.Replace(path, _PREFIX_SOURCE_DIR, this.runtime.Config.SourceDir)
				return Render(c, openDocument(path))
			}
		}

		return nil
	}
}

func (this *Web) handleOpenProcessedDir() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		text, err := this.runtime.OpenProcessedDir()
		return Render(c, message(text, err))
	}
}

func (this *Web) handleOpenSourceDir() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		text, err := this.runtime.OpenSourceDir()
		return Render(c, message(text, err))
	}
}

func (this *Web) handleFilesList() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return Render(c, filesList(this.runtime.FilesAll()))
	}
}

func (this *Web) handleFilesDuplistList() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return Render(c, filesDuplicates(this.runtime.FilesDuplicates()))
	}
}

func (this *Web) handleFilesDetails() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		files := this.runtime.FileGet(id)

		if len(files) == 1 {
			return Render(c, filesDetails(files[0]))
		}

		return nil
	}
}

func (this *Web) handleFilesUpdate() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		title := c.FormValue("title")
		files := this.runtime.State.Files.ById(id)

		if len(files) == 1 {
			files[0].Title = title
		}

		if err := this.runtime.SaveState(); err != nil {
			log.Println(err.Error())
		}

		c.Response().Header.Add("HX-Trigger", "file-updated")
		c.Response().Header.Add("HX-Trigger", attrId("file-updated", files[0].Id.String()))
		return nil
	}
}

func (this *Web) handleLabelsAll() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return Render(c, labelsList(this.runtime.Labels.All()))
	}
}
