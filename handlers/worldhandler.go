package handler

import (
	"encoding/json"
	"fmt"
	gopherlife "gopherlife/world"
	"html/template"
	"log"
	"net/http"
)

type SelectReturn struct {
	WorldRender string
	Gopher      *gopherlife.Gopher
}

type PageVariables struct {
	Data template.HTML
}

func SetUpPage() {

	var world = gopherlife.CreateWorld()
	renderer := gopherlife.NewRenderer()

	/*	for {
		world.ProcessWorld()
	}*/

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", worldToHTML(&world))
	http.HandleFunc("/PollWorld", ajaxProcessWorld(&world, &renderer))
	http.HandleFunc("/ShiftWorldView", ajaxHandleWorldInput(&world, &renderer))
	http.HandleFunc("/SelectGopher", ajaxSelectGopher(&world, &renderer))

	fmt.Println("Listening...")
	http.ListenAndServe(":8080", nil)

}

func worldToHTML(world *gopherlife.World) func(w http.ResponseWriter, r *http.Request) {

	renderer := gopherlife.NewRenderer()

	return func(w http.ResponseWriter, r *http.Request) {

		pageVariables := PageVariables{
			Data: template.HTML(renderer.RenderWorld(world).WorldRender),
		}

		t, err := template.ParseFiles("static/index.html")

		if err != nil {
			log.Print("Template parsing error: ", err)
		}

		err = t.Execute(w, pageVariables)

		if err != nil {
			log.Printf("Template executing error: ", err)
		}

	}

}

func ajaxProcessWorld(world *gopherlife.World, renderer *gopherlife.Renderer) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		//

		world.ProcessWorld()

		if true {
			jsonData, _ := json.Marshal(renderer.RenderWorld(world))
			w.Write(jsonData)
		} else {
			w.WriteHeader(404)
		}

	}
}

func ajaxSelectGopher(world *gopherlife.World, renderer *gopherlife.Renderer) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		r.ParseForm()

		position := r.FormValue("position")

		if _, ok := world.SelectEntity(position); ok {
			w.Header().Set("Content-Type", "application/json")
			jsonData, _ := json.Marshal(renderer.RenderWorld(world))
			w.Write(jsonData)
		} else {
			w.WriteHeader(404)
		}

	}
}

func ajaxHandleWorldInput(world *gopherlife.World, renderer *gopherlife.Renderer) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		r.ParseForm()

		var leftArrow = "37"
		var rightArrow = "39"
		var upArrow = "38"
		var downArrow = "40"

		var pKey = "80"

		var qKey = "81"
		var wKey = "87"

		keydown := r.FormValue("keydown")

		switch keydown {
		case wKey:
			world.UnSelectGopher()
		case qKey:
			world.SelectRandomGopher()
		case pKey:
			world.TogglePause()
		case leftArrow:
			renderer.ShiftRenderer(-1, 0)
		case rightArrow:
			renderer.ShiftRenderer(1, 0)
		case upArrow:
			renderer.ShiftRenderer(0, -1)
		case downArrow:
			renderer.ShiftRenderer(0, 1)
		}

	}
}
