package web

import (
	"github.com/boundedinfinity/statementer/label"
	"github.com/gofiber/fiber/v2"
)

const (
	LABEL_EDIT_PATH   = "/label/edit"
	LABEL_EXPAND_PATH = "/label/expand"
)

func (this *Web) initLabelViewRoutes() error {
	return nil
}

func (this *Web) initLabelEditRoutes() error {
	this.fiber.Get(attrPath(LABEL_EXPAND_PATH, ":id"), func(c *fiber.Ctx) error {
		if found, ok := this.runtime.Labels.ByIdStr(c.Params("id")); ok {
			found.Expanded = !found.Expanded
			return Render(c, labelTaxonomy([]*label.LabelViewModel{found}))
		}

		return Render(c, detailsNull())
	})

	this.fiber.Patch(attrPath(LABEL_EDIT_PATH, ":id"), func(c *fiber.Ctx) error {
		if found, ok := this.runtime.Labels.ByIdStr(c.Params("id")); ok {
			filters := append(found.Children, found, found.Parent)
			labels := this.runtime.Labels.List(label.WithoutFilter(filters...))

			return Render(c, labelEditForm(LABEL_EDIT_PATH, found, labels))
		}

		return Render(c, detailsNull())
	})

	this.fiber.Post(LABEL_EDIT_PATH, func(c *fiber.Ctx) error {
		if found, ok := this.runtime.Labels.ByIdStr(c.FormValue("id")); ok {
			found.Name = c.FormValue("name")
			found.Description = c.FormValue("description")
			parentId := c.FormValue("parent")

			if parentId == "" || parentId == "0" {
				found.Parent = nil
			} else {
				if parent, ok := this.runtime.Labels.ByIdStr(parentId); ok {
					found.Parent = parent
				}
			}

			this.runtime.SaveState()

			return Render(c, labelEditForm(
				LABEL_EDIT_PATH,
				found,
				this.runtime.Labels.List(label.WithoutFilter(found.Children...))),
			)
		}

		return Render(c, detailsNull())
	})

	return nil
}
