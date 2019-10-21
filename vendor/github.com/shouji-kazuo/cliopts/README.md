
## Requirements

- gopkg.in/urfave/cli.v2

## Install

```
go get github.com/shouji-kazuo/cliopts
```

## Example

### Read missing option from stdin

```golang
import (
  "fmt"
  "os"

  "github.com/shouji-kazuo/cliopts"
  cli "gopkg.in/urfave/cli.v2"
)

func main() {
  app := &cli.App{
    // ...

    // Define flags
    Flags: []cli.Flag{
      &cli.StringFlag{
        Name:  "username",
        Aliases: []string{"u"},
        Usage:   "set user name to login ts3card.com",
      },
      &cli.StringFlag{
        Name:  "password",
        Aliases: []string{"p"},
        Usage:   "set password to to login ts3card.com",
      },
    },
    // ...
    Action: func(c *cli.Context) error {
      // Declate variable that will be passed by --username
      // If --username option is missing, then username variable will be "".
      username := c.String("username")
      // Set up the function that ensure username variable is set by user.
      optEnsure := cliopts.NewEnsure().
        With("username", cliopts.StdInteract("Enter username: ").After(func(s string) error {
          username = s
          return nil
        }))
      // If --username option is missing, then prompt "Enter username: " show in stdout,
      // read string line from stdin, and set line to username variable.
      if err := optEnsure.Do(c); err != nil {
          return err
      }
      // ...
      return nil
    },
  }
  // ...
}
```
```
go build 
./main
Enter username: // do input something
...
```