package main

import (
	"fmt"
	"os"

	"github.com/tech-thinker/chatz/agents"
	"github.com/tech-thinker/chatz/utils"
	"github.com/urfave/cli/v2"
)

var (
    AppVersion = "v0.0.1"
    CommitHash = "unknown"
    BuildDate = "unknown"
)

func main() {
    
    var version bool
    var profile string
    var threadId string
    var output bool

    app := cli.NewApp()
    app.Name = "chatz"
    app.Description = "chatz is a versatile messaging app designed to send notifications to Google Chat, Slack, and Telegram."
    app.Flags = []cli.Flag {
        &cli.BoolFlag{
            Name: "output",
            Aliases: []string{"o"},
            Usage: "Print output to stdout",
            Destination: &output,
        },
        &cli.StringFlag{
            Name: "profile",
            Aliases: []string{"p"},
            Value: "default",
            Usage: "Profile from .chatz.ini",
            Destination: &profile,
        },
        &cli.StringFlag{
            Name: "thread-id",
            Aliases: []string{"t"},
            Value: "",
            Usage: "Thread ID for reply to a message",
            Destination: &threadId,
        },
        &cli.BoolFlag{
            Name: "version",
            Aliases: []string{"v"},
            Usage: "Print the version number",
            Destination: &version,
        },
    }
    app.Action = func(ctx *cli.Context) error {
        if version {
            fmt.Println("chatz version: ", AppVersion)
            fmt.Println("Commit Hash: ", CommitHash)
            fmt.Println("Build Date: ", BuildDate)
            return nil
        }


        var message string
        if ctx.Args().Len() == 0 {
            fmt.Println("Please provide a message.")
            return nil
        }
        for i, a := range ctx.Args().Slice() {
            if i == 0 {
                message = a
                continue
            }
            message = fmt.Sprintf("%s %s",message, a)
        }
        
        env, err := utils.LoadEnv(profile)
        if err!=nil {
            return nil
        }
        var agent agents.Agent
        switch env.Provider {
            case "slack":
                agent = agents.NewSlackAgent(env)
            case "google":
                agent = agents.NewGoogleAgent(env)
            case "telegram":
                agent = agents.NewTelegramAgent(env)
            default:
                fmt.Println("No valid provider. Please choose from [slack, google, telegram].")
                return nil
        }

        if len(threadId) > 0 {
            res, _ := agent.Reply(threadId, message)
            if output {
                fmt.Println(res)
            }
            return nil
        }
        res, _ := agent.Post(message)
        if output {
            fmt.Println(res)
        }

        return nil
    }

    if err := app.Run(os.Args); err != nil {
		panic(err)
	}

}

