package web

import (
	"log"
	"strconv"
	"time"

	"github.com/a-h/templ"
	"github.com/boundedinfinity/go-commoner/idiomatic/slicer"
	"github.com/boundedinfinity/go-commoner/idiomatic/stringer"
	"github.com/boundedinfinity/statementer/model"
	"github.com/boundedinfinity/statementer/runtime"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/google/uuid"
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
	this.fiber = fiber.New(fiber.Config{
		// https://docs.gofiber.io/#zero-allocation
		Immutable: true,
	})
	this.fiber.Use(logger.New())

	this.fiber.Get("/", func(c *fiber.Ctx) error {
		return Render(c, home(this.runtime.Config))
	})

	this.initFileRoutes()
	this.initLabelRoutes()
	this.initOtherRoutes()

	return nil
}

func (this *Web) initLabelRoutes() error {
	this.fiber.Get("/labels/select/:id", func(c *fiber.Ctx) error {
		if id, err := uuid.Parse(c.Params("id")); err != nil {
			log.Println(err.Error())
		} else {
			if label, ok := this.runtime.Labels.AddSelected(id); ok {
				this.setTrigger(c, "label-selected")
				return Render(c, labelView(label))
			}
		}

		return nil
	})

	this.fiber.Delete("/labels/select/:id", func(c *fiber.Ctx) error {
		if id, err := uuid.Parse(c.Params("id")); err != nil {
			log.Println(err.Error())
		} else {
			if label, ok := this.runtime.Labels.RemoveSelected(id); ok {
				this.setTrigger(c, "label-selected")
				return Render(c, labelView(label))
			}
		}

		return nil
	})

	this.fiber.Get("/labels/all", func(c *fiber.Ctx) error {
		return Render(c, labelList(this.runtime.Labels.All()))
	})

	this.fiber.Post("/labels/year/:year", func(c *fiber.Ctx) error {
		year := c.Params("year")
		var yearInt int
		var err error

		if year == "this" {
			yearInt = time.Now().Year()
		} else {
			yearInt, err = strconv.Atoi(year)
			if err != nil {
				log.Println(err.Error())
			}
		}

		if err = this.runtime.Labels.GenerateYear(yearInt); err != nil {
			log.Println(err.Error())
		}

		if err == nil {
			if err = this.runtime.SaveState(); err != nil {
				log.Println(err.Error())
			} else {
				this.setTrigger(c, "label-updated")
			}
		}

		return nil
	})

	this.fiber.Get("/labels/button", func(c *fiber.Ctx) error {
		return Render(c, labelFormButton())
	})

	this.fiber.Get("/labels/new", func(c *fiber.Ctx) error {
		return Render(c, labelNewForm())
	})

	this.fiber.Post("/labels/new", func(c *fiber.Ctx) error {
		name := c.FormValue("name")
		desc := c.FormValue("description")
		label := model.SimpleLabel{Name: name, Description: desc}

		if err := this.runtime.Labels.Add(&label); err != nil {
			log.Println(err.Error())
		}

		if err := this.runtime.SaveState(); err != nil {
			log.Println(err.Error())
		}

		this.setTrigger(c, "label-updated")
		return Render(c, labelFormButton())
	})

	return nil
}

func (this *Web) initFileRoutes() error {
	this.fiber.Get("/files/list", func(c *fiber.Ctx) error {
		return Render(c, filesList(this.runtime.FilesAllFiltered()))
	})

	this.fiber.Get("/files/duplicates", func(c *fiber.Ctx) error {
		return Render(c, filesDuplicates(this.runtime.FilesDuplicates()))
	})

	this.fiber.Get("/files/merge", func(c *fiber.Ctx) error {
		return nil
	})

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

		this.setTrigger(c, "file-updated")
		return Render(c, fileViewTitle(files[0]))
	})

	this.fiber.Get("/files/label/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		files := this.runtime.State.Files.ById(id)
		return Render(c, fileLabelView(files[0]))
	})

	this.fiber.Patch("/files/label/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		files := this.runtime.State.Files.ById(id)
		return Render(c, fileLabelEdit(
			files[0],
			labelSetChecked(this.runtime.State.Labels, files[0].Labels),
		))
	})

	this.fiber.Post("/files/label/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		labelIds := this.formValues(c, "label")
		files := this.runtime.State.Files.ById(id)
		labels := []*model.SimpleLabel{}

		for _, labelId := range labelIds {
			id, err := uuid.Parse(labelId)
			if err != nil {
				log.Println(err.Error())
			}

			if label, ok := this.runtime.Labels.ById(id); ok {
				labels = append(labels, label)
			}
		}

		files[0].Labels = labels

		if err := this.runtime.SaveState(); err != nil {
			log.Println(err.Error())
		}

		this.setTrigger(c, "file-updated")
		return Render(c, fileLabelView(files[0]))
	})

	return nil
}

func labelSetChecked(all, file []*model.SimpleLabel) []*model.SimpleLabel {
	copies := slicer.Map(func(_ int, label *model.SimpleLabel) *model.SimpleLabel {
		copy := model.SimpleLabelCopy(*label)
		return &copy
	}, all...)

	group := map[uuid.UUID]*model.SimpleLabel{}

	for _, label := range copies {
		group[label.Id] = label
	}

	for _, label := range file {
		group[label.Id].Checked = true
	}

	return copies
}

func (this *Web) initOtherRoutes() error {
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

// https://github.com/a-h/templ/blob/main/examples/integration-gofiber/main.go/

func Render(c *fiber.Ctx, component templ.Component) error {
	c.Set("Content-Type", "text/html")
	return component.Render(c.Context(), c.Response().BodyWriter())
}

func (this *Web) setTrigger(c *fiber.Ctx, triggers ...string) {
	c.Response().Header.DisableNormalizing()
	c.Response().Header.Add("HX-Trigger", stringer.Join(", ", triggers...))
}

func (this *Web) formValues(c *fiber.Ctx, name string) []string {
	body := string(c.Body())
	body = stringer.Replace(body, "", "?")
	params := stringer.Split(body, "&")
	kvs := map[string][]string{}

	for _, param := range params {
		kv := stringer.Split(param, "=")
		if _, ok := kvs[kv[0]]; !ok {
			kvs[kv[0]] = []string{}
		}

		kvs[kv[0]] = append(kvs[kv[0]], kv[1])
	}

	return kvs[name]
}
