//go:build darwin || linux
// +build darwin linux

package main

// World is a basic gomobile app.
import (
	"errors"
	"flag"
	"image"
	"log"
	"os"

	"github.com/trimble-oss/tierceron-nute-core/mashupsdk"
	sdk "github.com/trimble-oss/tierceron-nute-core/mashupsdk"
	"github.com/trimble-oss/tierceron-nute/mashupsdk/guiboot"
	"github.com/trimble-oss/tierceron-nute/mashupsdk/server"
	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/event/touch"
	"golang.org/x/mobile/exp/sprite"
	"golang.org/x/mobile/exp/sprite/glsprite"
	"golang.org/x/mobile/gl"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	engine sprite.Engine
)

var worldCompleteChan chan bool

type worldApiHandler struct {
}

type worldClientInitHandler struct {
}

type WorldApp struct {
	wApiHandler *worldApiHandler
	mainWin     *app.App
	mainWinDims *image.Point
}

var worldApp WorldApp

func (w *WorldApp) InitMainWindow() {

	var mobileWinHandler interface{} = func(a app.App) {
		var glCtx gl.Context
		var szEvent size.Event

		for event := range a.Events() {
			switch filteredEvent := a.Filter(event).(type) {
			case lifecycle.Event:
				switch filteredEvent.Crosses(lifecycle.StageVisible) {
				case lifecycle.CrossOn:
					glCtx, _ = filteredEvent.DrawContext.(gl.Context)
					onStart(glCtx)
					a.Send(paint.Event{})
				case lifecycle.CrossOff:
				}
			case size.Event:
				// Capture screen sizing.
				szEvent = filteredEvent
			case paint.Event:
				// Do some painting.
				onPaint(glCtx, szEvent)
			case touch.Event:
				// Capture a screen touched event.
			case key.Event:
				// Capture general keyboard events.
			}
		}
	}

	guiboot.InitMainWindow(guiboot.Gomobile, nil, mobileWinHandler)
}

func (w *worldClientInitHandler) RegisterContext(context *mashupsdk.MashupContext) {
	// TODO: implement something meaningful.
}

func (w *worldApiHandler) OnDisplayChange(displayHint *sdk.MashupDisplayHint) {
	log.Printf("Received OnDisplayChange.")
	if worldApp.mainWin == nil {
		worldApp.InitMainWindow()
	} else {
		(*worldApp.mainWin).Send(size.Event{WidthPx: int(displayHint.Width), HeightPx: int(displayHint.Height)})
	}
}

func (c *worldApiHandler) ResetStates() {
	// Not implemented.
}

func (c *worldApiHandler) GetElements() (*mashupsdk.MashupDetailedElementBundle, error) {
	// Not implemented.
	return &mashupsdk.MashupDetailedElementBundle{}, errors.New("Could not get items.")
}

func (c *worldApiHandler) UpsertElements(detailedElementBundle *sdk.MashupDetailedElementBundle) (*sdk.MashupDetailedElementBundle, error) {
	// Not implemented.
	return nil, errors.New("Could not capture items.")
}

func (c *worldApiHandler) TweakStates(elementStateBundle *sdk.MashupElementStateBundle) (*sdk.MashupElementStateBundle, error) {
	// Not implemented.
	return nil, errors.New("Could not capture items.")
}

func (c *worldApiHandler) TweakStatesByMotiv(motivIn *mashupsdk.Motiv) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, errors.New("Could not capture items.")
}

func main() {
	callerCreds := flag.String("CREDS", "", "Credentials of caller")
	insecure := flag.Bool("tls-skip-validation", false, "Skip server validation")
	flag.Parse()
	worldLog, err := os.OpenFile("world.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf(err.Error(), err)
	}
	log.SetOutput(worldLog)

	worldApp = WorldApp{wApiHandler: &worldApiHandler{}}

	server.InitServer(*callerCreds, *insecure, 0, worldApp.wApiHandler, nil)

	<-worldCompleteChan
}

func onStart(glCtx gl.Context) {
	engine = glsprite.Engine(nil)
}

func onPaint(glCtx gl.Context, szEvent size.Event) {
	glCtx.ClearColor(1, 1, 1, 1)
	glCtx.Clear(gl.COLOR_BUFFER_BIT)
}
