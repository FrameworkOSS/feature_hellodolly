package hellodolly

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/FrameworkOSS/portal/features/commands/handler"
	"github.com/FrameworkOSS/portal/portal"
)

var (
	cmdHelloDolly = handler.NewCommand().
			SetID("hellodolly").
			SetName("Hello Dolly").
			SetAbout("Returns a random lyric from Louis Armstrong's \"Hello, Dolly.\"").
			SetUsage("Simply call the command without any arguments to receive a random lyric.").
			SetAliases("hd").
			SetRequiresPreprocessing(true).
			SetArgument(handler.NewCommandArg().
				SetID("lyric").
				SetName("lyric line").
				SetAbout("The specific lyric line number to retrieve, or a random line if not specified.").
				SetUsage(fmt.Sprintf("Provide a lyric line number in the range of 1 to %d.", len(lyrics))).
				SetAliases("line", "number").
				SetType(handler.CommandArgTypeNumber).
				SetRequiresValue(true),
		).
		SetSubcommand(cmdTest)
	cmdTest = handler.NewCommand().
		SetID("test").
		SetName("Test Command").
		SetAbout("A test subcommand for Hello Dolly").
		SetUsage("This is just a test subcommand.").
		SetAliases("t")

	lyrics = `Hello, Dolly
Well, hello, Dolly
It's so nice to have you back where you belong
You're lookin' swell, Dolly
I can tell, Dolly
You're still glowin', you're still crowin'
You're still goin' strong
We feel the room swayin'
While the band's playin'
One of your old favourite songs from way back when
So, take her wrap, fellas
Find her an empty lap, fellas
Dolly'll never go away again
Hello, Dolly
Well, hello, Dolly
It's so nice to have you back where you belong
You're lookin' swell, Dolly
I can tell, Dolly
You're still glowin', you're still crowin'
You're still goin' strong
We feel the room swayin'
While the band's playin'
One of your old favourite songs from way back when
Golly, gee, fellas
Find her a vacant knee, fellas
Dolly'll never go away
Dolly'll never go away
Dolly'll never go away again`
)

type HelloDolly struct {
	lockView  sync.Mutex
	lockResp  sync.Mutex
	views     map[int]int //line:views
	resps     []*portal.Event
	lastResp  time.Time
	processor *handler.EventCommandHandler
}

func NewHelloDolly() (hd *HelloDolly) {
	hd = new(HelloDolly)
	hd.views = make(map[int]int)
	hd.resps = make([]*portal.Event, 0)
	hd.processor = handler.NewEventCommandHandler()

	hd.processor.GetCommandHandler().
		Handle(hd.cmdHelloDolly, cmdHelloDolly)

	return
}

func (hd *HelloDolly) cmdHelloDolly(cmd *handler.Command, e *portal.Event) (err error) {
	var r *portal.Event

	switch cmd.GetID() {
	case cmdHelloDolly.GetID():
		line := 0
		if lyric := cmd.GetArgument("lyric"); lyric != nil {
			line = lyric.GetValueNumber()
		}

		isRandom := false
		lines := strings.Split(lyrics, "\n")
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

		r = portal.NewEventResponse(hd.ID(), []byte(lineFmt))
	case cmdHelloDolly.GetID() + " " + cmdTest.GetID():
		r = portal.NewEventResponse(hd.ID(), []byte(
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

func (hd *HelloDolly) storeResp(e *portal.Event) {
	now := time.Now()
	hd.lockResp.Lock()
	hd.lastResp = now
	hd.resps = append(hd.resps, e)
	hd.lockResp.Unlock()
}

func (hd *HelloDolly) readResp() (e *portal.Event) {
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
	return 0
}

func (hd *HelloDolly) ID() string {
	return "hellodolly"
}

func (hd *HelloDolly) Name() string {
	return "Hello Dolly"
}

func (hd *HelloDolly) Authors() []string {
	return []string{"JoshuaDoes"}
}

func (hd *HelloDolly) Description() string {
	return "Inspired by the WordPress sample plugin! Responds with a random lyric from Louis Armstrong's \"Hello, Dolly.\""
}

func (hd *HelloDolly) Version() string {
	return "v0.0.1"
}

func (hd *HelloDolly) Open() error {
	hd.storeResp(handler.NewEventCommandAdd(hd.ID(), cmdHelloDolly))
	hd.storeResp(portal.NewEventReady(hd.ID(), true))
	return nil
}

func (hd *HelloDolly) Close() (errs []error, retry bool) {
	return
}

func (hd *HelloDolly) Input(e *portal.Event) error {
	return hd.processor.Process(e)
}

func (hd *HelloDolly) Output() (*portal.Event, error) {
	return hd.readResp(), nil
}
