package main

import (
	"encoding/json"
	"fmt"
)
import "math"
import "time"
import "net/http"
import "io/ioutil"
import "server"

import ui "github.com/gizak/termui"

var svrData server.SvrStatus

const SPACE = " 																																							"

func test1() {
	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	ui.UseTheme("helloworld")

	par0 := ui.NewPar("\t\t性能探针服务端")
	par0.Height = 3
	par0.Width = ui.TermWidth()
	par0.Y = 1
	par0.HasBorder = false
	//	par0.Border.Label = "Dash board"

	par1 := ui.NewPar("")
	par1.Height = 3
	par1.Width = ui.TermWidth()
	par1.X = 0
	par1.Y = 3
	par1.Border.Label = "Client Count"

	ui.Render(par0, par1)

	redraw := make(chan bool)

	update := func() {
		for {
			time.Sleep(time.Second)
			str1 := time.Now().Format("2006-01-02 15:04:05")
			str2 := fmt.Sprintf("client count:%d (socket connent)%s %s", svrData.CurrentClientCount, SPACE, str1)
			par1.Text = str2

			redraw <- true
		}
	}

	go update()

	evt := ui.EventCh()

	for {
		select {
		case e := <-evt:
			if e.Type == ui.EventKey && e.Ch == 'q' {
				return
			}
			if e.Type == ui.EventResize {
				ui.Body.Width = ui.TermWidth()
				ui.Body.Align()
				// go func() { redraw <- true }()
			}

		case <-redraw:
			ui.Render(par0, par1)
		}
	}

	<-ui.EventCh()
}

func titleView() {

}

func readServerData() {
	response, _ := http.Get("http://127.0.0.1:6843/status")
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	//	fmt.Println(string(body))

	if err := json.Unmarshal(body, &svrData); err == nil {
		//		fmt.Println("haha", svrData.CurrentClientCount)
	}
}

func main() {

	go func() {
		for _ = range time.Tick(1 * time.Second) {
			readServerData()
		}
	}()

	test1()
	return
	fmt.Println("DashBoard")

	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	sinps := (func() []float64 {
		n := 400
		ps := make([]float64, n)
		for i := range ps {
			ps[i] = 1 + math.Sin(float64(i)/5)
		}
		return ps
	})()

	sinpsint := (func() []int {
		ps := make([]int, len(sinps))
		for i, v := range sinps {
			ps[i] = int(100*v + 10)
		}
		return ps
	})()

	ui.UseTheme("helloworld")

	spark := ui.Sparkline{}
	spark.Height = 8
	spdata := sinpsint
	spark.Data = spdata[:100]
	spark.LineColor = ui.ColorCyan
	spark.TitleColor = ui.ColorWhite

	sp := ui.NewSparklines(spark)
	sp.Height = 11
	sp.Border.Label = "Sparkline"

	lc := ui.NewLineChart()
	lc.Border.Label = "braille-mode Line Chart"
	lc.Data = sinps
	lc.Height = 11
	lc.AxesColor = ui.ColorWhite
	lc.LineColor = ui.ColorYellow | ui.AttrBold

	gs := make([]*ui.Gauge, 3)
	for i := range gs {
		gs[i] = ui.NewGauge()
		gs[i].Height = 2
		gs[i].HasBorder = false
		gs[i].Percent = i * 10
		gs[i].PaddingBottom = 1
		gs[i].BarColor = ui.ColorRed
	}

	ls := ui.NewList()
	ls.HasBorder = false
	ls.Items = []string{
		"[1] Downloading File 1",
		"", // == \newline
		"[2] Downloading File 2",
		"",
		"[3] Uploading File 3",
	}
	ls.Height = 5

	par := ui.NewPar("<> This row has 3 columns\n<- Widgets can be stacked up like left side\n<- Stacked widgets are treated as a single widget")
	par.Height = 5
	par.Border.Label = "Demonstration"

	// build layout
	ui.Body.AddRows(
		ui.NewRow(
			ui.NewCol(6, 0, sp),
			ui.NewCol(6, 0, lc)),
		ui.NewRow(
			ui.NewCol(3, 0, ls),
			ui.NewCol(3, 0, gs[0], gs[1], gs[2]),
			ui.NewCol(6, 0, par)))

	// calculate layout
	ui.Body.Align()

	done := make(chan bool)
	redraw := make(chan bool)

	update := func() {
		for i := 0; i < 103; i++ {
			for _, g := range gs {
				g.Percent = (g.Percent + 3) % 100
			}

			sp.Lines[0].Data = spdata[:100+i]
			lc.Data = sinps[2*i:]

			time.Sleep(time.Second / 2)
			redraw <- true
		}
		done <- true
	}

	evt := ui.EventCh()

	ui.Render(ui.Body)
	go update()

	for {
		select {
		case e := <-evt:
			if e.Type == ui.EventKey && e.Ch == 'q' {
				return
			}
			if e.Type == ui.EventResize {
				ui.Body.Width = ui.TermWidth()
				ui.Body.Align()
				go func() { redraw <- true }()
			}
		case <-done:
			return
		case <-redraw:
			ui.Render(ui.Body)
		}
	}
}
