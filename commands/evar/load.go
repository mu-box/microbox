package evar

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/spf13/cobra"

	"github.com/mu-box/microbox/commands/steps"
	"github.com/mu-box/microbox/helpers"
	"github.com/mu-box/microbox/models"
	app_evar "github.com/mu-box/microbox/processors/app/evar"
	production_evar "github.com/mu-box/microbox/processors/evar"
	"github.com/mu-box/microbox/util"
	"github.com/mu-box/microbox/util/config"
	"github.com/mu-box/microbox/util/display"
)

// LoadCmd loads variables from a file.
var LoadCmd = &cobra.Command{
	Use:   "load [local|dry-run|remote-alias] filename",
	Short: "Loads environment variable(s) from a file",
	Long: `Loads environment variable(s) from a file.

The alias must be used when loading variables for a production app.
If you would like to load variables for a different app, please add
it as a remote as follows: 'microbox remote add <APPNAME> <ALIAS>'.
You may then perform the 'load' again, specifying that alias.`,
	Run: loadFn,
}

// loadFn parses a specified file and adds the contained variables to microbox.
// Read in the file, strip out 'export ', parse, add resulting vars
func loadFn(ccmd *cobra.Command, args []string) {
	// parse the evars excluding the context
	env, _ := models.FindEnvByID(config.EnvID())
	args, location, name := helpers.Endpoint(env, args, 0)
	vars, err := loadVars(args, fileGetter{})
	if err != nil {
		display.CommandErr(util.Err{
			Message: err.Error(),
			Code:    "USER",
			Stack:   []string{"failed to load evars from file"},
		})
		return
	}

	evars := parseEvars(vars)

	switch location {
	case "local":
		app, _ := models.FindAppBySlug(config.EnvID(), name)
		display.CommandErr(app_evar.Add(env, app, evars))
	case "production":
		steps.Run("login")(ccmd, args)

		production_evar.Add(env, name, evars)
	}
}

// loadVars loads variables from filenames passed in
func loadVars(args []string, getter contentGetter) ([]string, error) {
	vars := []string{}

	for _, filename := range args {
		contents, err := getter.getContents(filename)
		if err != nil {
			return nil, fmt.Errorf("Failed to get var contents - %s", err.Error())
		}

		// normalize file `key=val`
		newthings := strings.Replace(contents, "export ", "", -1)

		// strip commented (#) lines
		newthings = regexp.MustCompilePOSIX(`^#.*$\n*`).ReplaceAllString(newthings, "")

		// strip out blank lines
		newthings = regexp.MustCompilePOSIX(`\n\n+`).ReplaceAllString(newthings, "\n")

		// strip trailing newline
		newthings = regexp.MustCompilePOSIX(`\n$`).ReplaceAllString(newthings, "")

		// get index of variable start (change regex to allow more than alphanumeric chars as keys)
		indexes := regexp.MustCompilePOSIX(`(^[a-zA-Z0-9_]*?)=(\"\n|.)`).FindAllStringIndex(newthings, -1)

		start := 0
		for i := range indexes {
			end := indexes[i][0]
			if end <= 2 {
				continue
			}
			// end-1 leaves off the newline after the variable declaration
			vars = append(vars, newthings[start:end-1])
			start = end
		}
		// the newline after this variable declaration would have been previously stripped off
		vars = append(vars, newthings[start:])
	}

	return vars, nil
}

// contentGetter is an interface to allow us to test loading/parsing of variables.
type contentGetter interface {
	getContents(filename string) (string, error)
}

type fileGetter struct{}

func (f fileGetter) getContents(filename string) (string, error) {
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("Failed to read file - %s", err.Error())
	}
	return string(contents), nil
}
