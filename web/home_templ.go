// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.778
package web

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import "github.com/boundedinfinity/statementer/model"

func home(config model.Config) templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<!doctype html><html lang=\"en\" data-theme=\"dark\"><head><meta charset=\"UTF-8\"><meta http-equiv=\"X-UA-Compatible\" content=\"IE=edge\"><meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\"><meta name=\"description\" content=\"Bounded Infinity Statement Management System\"><meta name=\"google\" content=\"notranslate\"><meta name=\"htmx-config\" content=\"{&#34;globalViewTransitions&#34;:&#34;true&#34;}\"><link rel=\"shortcut icon\" href=\"/img/gopher-svgrepo-com.svg\" type=\"image/svg+xml\"><link href=\"/css/daisyui.min.css\" rel=\"stylesheet\"><link href=\"/css/styles.css\" rel=\"stylesheet\" type=\"text/css\"><script src=\"/js/tailwind.min.js\"></script><script src=\"/js/htmx.min.js\"></script><title>Bounded Infinity : Statement Management System</title></head><body hx-boost=\"true\"><div class=\"min-h-screen flex flex-col h-screen gap-2 min-w-full\"><header hx-get=\"/labels/all\" hx-trigger=\"load,label-updated\" hx-target=\"#labels\" class=\"p-2 mt-2 border-2 border-r-2 border-slate-600\"><div id=\"labels\">Labels...</div><div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = labelFormButton().Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div></header><div class=\"flex-1 flex flex-row overflow-y-hidden gap-2\"><nav class=\"w-2/12 overflow-y-auto flex flex-col flex-wrap content-start gap-2 p-2\"><div>Config Files:</div><div class=\"ml-4\"><button hx-get=\"/open/config-file\" hx-target=\"#stdout\" class=\"btn btn-primary btn-outline\">Open</button></div><div>Repository Dir:</div><div class=\"ml-4\"><button hx-get=\"/open/repository-dir\" hx-target=\"#stdout\" class=\"btn btn-primary btn-outline\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var2 string
		templ_7745c5c3_Var2, templ_7745c5c3_Err = templ.JoinStringErrs(config.RepositoryDir)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `web/home.templ`, Line: 63, Col: 30}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var2))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</button></div><div>Source Dir:</div><div class=\"ml-4\"><button hx-get=\"/open/source-dir\" hx-target=\"#stdout\" class=\"btn btn-primary btn-outline\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var3 string
		templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinStringErrs(config.SourceDir)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `web/home.templ`, Line: 75, Col: 26}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var3))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</button></div><hr><button hx-get=\"/files/list\" hx-target=\"#results\" class=\"btn btn-primary btn-outline\">All Files</button> <button hx-get=\"/files/duplicates\" hx-target=\"#results\" class=\"btn btn-primary btn-outline\">Duplicate Files</button></nav><main class=\"overflow-y-auto border w-5/12\" hx-get=\"/files/list\" hx-trigger=\"load,file-updated\" hx-target=\"#results\"><div id=\"results\"></div></main><aside class=\"w-4/12 h-[1600px]\"><div id=\"details\" class=\"h-full p-2 border-2 border-r-2 border-slate-600\"></div></aside></div><!-- end main container --><footer class=\"bg-gray-100\"><div id=\"stdout\" class=\"mt-2 p-2 border-2 border-r-2 border-slate-600\">This is a test...</div></footer></div></body></html>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}

var _ = templruntime.GeneratedTemplate
