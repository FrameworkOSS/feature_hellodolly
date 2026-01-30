package metadata

import (
	"fmt"

	"github.com/FrameworkOSS/feature_commands/handler"
)

const (
	API         = 0
	ID          = "hellodolly"
	Name        = "Hello Dolly"
	Authors     = "JoshuaDoes"
	Description = "Inspired by the WordPress sample plugin! Responds with a random lyric from Louis Armstrong's \"Hello, Dolly.\""
	Version     = "v0.0.1"
)

var (
	CmdHelloDolly = handler.NewCommand().
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
				SetUsage(fmt.Sprintf("Provide a lyric line number in the range of 1 to %d.", len(Lyrics))).
				SetAliases("line", "number").
				SetType(handler.CommandArgTypeNumber).
				SetRequiresValue(true),
		).
		SetSubcommand(CmdTest)
	CmdTest = handler.NewCommand().
		SetID("test").
		SetName("Test Command").
		SetAbout("A test subcommand for Hello Dolly").
		SetUsage("This is just a test subcommand.").
		SetAliases("t")

	Lyrics = `Hello, Dolly
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
