package hellodolly

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/FrameworkOSS/event"
	"github.com/FrameworkOSS/feature_commands/handler"
	"github.com/FrameworkOSS/feature_hellodolly/metadata"
)

type HelloDolly struct {
	lockView  sync.Mutex
	lockResp  sync.Mutex
	views     map[int]int //line:views
	resps     []*event.Event
	lastResp  time.Time
	processor *handler.EventCommandHandler
}

func NewHelloDolly() (hd *HelloDolly) {
	hd = new(HelloDolly)
	hd.views = make(map[int]int)
	hd.resps = make([]*event.Event, 0)
	hd.processor = handler.NewEventCommandHandler()

	hd.processor.GetCommandHandler().
		Handle(hd.cmdHelloDolly, metadata.CmdHelloDolly)

	return
}

func (hd *HelloDolly) cmdHelloDolly(cmd *handler.Command, e *event.Event) (err error) {
	var r *event.Event

	switch cmd.GetID() {
	case metadata.CmdHelloDolly.GetID():
		line := 0
		if lyric := cmd.GetArgument("lyric"); lyric != nil {
			line = lyric.GetValueNumber()
		}

		isRandom := false
		lines := strings.Split(metadata.Lyrics, "\n")
		if line < 1 || line >= len(lines) {
			//Pick a random lyric!
			line = rand.Intn(len(lines))
			isRandom = true
		} else {
			line--
		}

		lineS := lines[line]
		views := hd.viewCount(line)
		plural := ""
		if views != 1 {
			plural = "s"
		}
		random := ""
		if isRandom {
			random = "random "
		}
		lineFmt := fmt.Sprintf("<code>%s</code>\n\nThis %slyric (line %d) has reached <u>%d</u> view%s!\nThe last response was generated <i>%s</i> ago.", lineS, random, line+1, views, plural, time.Since(hd.lastResp).String())

		r = event.NewEventResponse(hd.ID(), []byte(lineFmt))
	case metadata.CmdHelloDolly.GetID() + " " + metadata.CmdTest.GetID():
		r = event.NewEventResponse(hd.ID(), []byte(
			fmt.Sprintf("This is a test response from the <u>Hello Dolly</u> feature's test subcommand!\n\nThe last response was generated <i>%s</i> ago.", time.Since(hd.lastResp).String())),
		)
	default:
		err = fmt.Errorf("hellodolly: unknown command: %s", cmd.GetID())
	}

	if r != nil {
		if r.GetChannel() == "" {
			r.SetChannel(e.GetChannel())
		}
		if r.GetProducer() == "" {
			r.SetProducer(hd.ID())
		}
		r.AddParticipants(e.GetParticipants()...)
		r.AddParticipants(e.GetProducer())
		go hd.storeResp(r)
	}

	return
}

func (hd *HelloDolly) storeResp(e *event.Event) {
	now := time.Now()
	hd.lockResp.Lock()
	hd.lastResp = now
	hd.resps = append(hd.resps, e)
	hd.lockResp.Unlock()
}

func (hd *HelloDolly) readResp() (e *event.Event) {
	if len(hd.resps) > 0 {
		hd.lockResp.Lock()
		e = hd.resps[0]
		hd.resps = hd.resps[1:]
		hd.lockResp.Unlock()
	}
	return
}

func (hd *HelloDolly) viewCount(line int) int {
	hd.lockView.Lock()
	defer hd.lockView.Unlock()
	views := 1
	if count, exists := hd.views[line]; exists {
		views = count + 1
	}
	hd.views[line] = views
	return views
}

func (hd *HelloDolly) API() int {
	return metadata.API
}

func (hd *HelloDolly) ID() string {
	return metadata.ID
}

func (hd *HelloDolly) Name() string {
	return metadata.Name
}

func (hd *HelloDolly) Authors() []string {
	return strings.Split(metadata.Authors, ",")
}

func (hd *HelloDolly) Description() string {
	return metadata.Description
}

func (hd *HelloDolly) Version() string {
	return metadata.Version
}

func (hd *HelloDolly) Open() error {
	hd.storeResp(handler.NewEventCommandAdd(hd.ID(), metadata.CmdHelloDolly))
	hd.storeResp(event.NewEventReady(hd.ID(), true))
	return nil
}

func (hd *HelloDolly) Close() (errs []error, retry bool) {
	return
}

func (hd *HelloDolly) Input(e *event.Event) error {
	return hd.processor.Process(e)
}

func (hd *HelloDolly) Output() (*event.Event, error) {
	return hd.readResp(), nil
}
