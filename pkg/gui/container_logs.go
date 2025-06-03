package gui

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/fatih/color"
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazydocker/pkg/commands"
	"github.com/jesseduffield/lazydocker/pkg/tasks"
	"github.com/jesseduffield/lazydocker/pkg/utils"
)

func (gui *Gui) renderContainerLogsToMain(container *commands.Container) tasks.TaskFunc {
	return gui.NewTickerTask(TickerTaskOpts{
		Func: func(ctx context.Context, notifyStopped chan struct{}) {
			gui.renderContainerLogsToMainAux(container, ctx, notifyStopped)
		},
		Duration: time.Millisecond * 200,
		// TODO: see why this isn't working (when switching from Top tab to Logs tab in the services panel, the tops tab's content isn't removed)
		Before:     func(ctx context.Context) { gui.clearMainView() },
		Wrap:       gui.Config.UserConfig.Gui.WrapMainPanel,
		Autoscroll: true,
	})
}

func (gui *Gui) renderContainerLogsToMainAux(container *commands.Container, ctx context.Context, notifyStopped chan struct{}) {
	gui.clearMainView()
	defer func() {
		notifyStopped <- struct{}{}
	}()

	mainView := gui.Views.Main

	if err := gui.writeContainerLogs(container, ctx, mainView); err != nil {
		gui.Log.Error(err)
	}

	// if we are here because the task has been stopped, we should return
	// if we are here then the container must have exited, meaning we should wait until it's back again before
	ticker := time.NewTicker(time.Millisecond * 100)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			result, err := container.Inspect()
			if err != nil {
				// if we get an error, then the container has probably been removed so we'll get out of here
				gui.Log.Error(err)
				return
			}
			if result.State.Running {
				return
			}
		}
	}
}

func (gui *Gui) renderLogsToStdout(container *commands.Container) {
	stop := make(chan os.Signal, 1)
	defer signal.Stop(stop)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		signal.Notify(stop, os.Interrupt)
		<-stop
		cancel()
	}()

	if err := gui.g.Suspend(); err != nil {
		gui.Log.Error(err)
		return
	}

	defer func() {
		if err := gui.g.Resume(); err != nil {
			gui.Log.Error(err)
		}
	}()

	if err := gui.writeContainerLogs(container, ctx, os.Stdout); err != nil {
		gui.Log.Error(err)
		return
	}

	gui.promptToReturn()
}

func (gui *Gui) promptToReturn() {
	if !gui.Config.UserConfig.Gui.ReturnImmediately {
		fmt.Fprintf(os.Stdout, "\n\n%s", utils.ColoredString(gui.Tr.PressEnterToReturn, color.FgGreen))

		// wait for enter press
		if _, err := fmt.Scanln(); err != nil {
			gui.Log.Error(err)
		}
	}
}

type ContainerLogsView struct {
    lines     []string
    selected  int
    scroll    int
    maxHeight int
	writer    *io.Writer
}
func (gui *Gui) writeContainerLogs(ctr *commands.Container, ctx context.Context, writer io.Writer) error {
    gui.State.LogsView = &ContainerLogsView{
        lines:    make([]string, 0),
        selected: -1,
        scroll:   0,
		writer:   &writer,
    }
	logsView := gui.State.LogsView

    readCloser, err := gui.DockerCommand.Client.ContainerLogs(ctx, ctr.ID, container.LogsOptions{
        ShowStdout: true,
        ShowStderr: true,
        Timestamps: gui.Config.UserConfig.Logs.Timestamps,
        Since:      gui.Config.UserConfig.Logs.Since,
        Tail:       gui.Config.UserConfig.Logs.Tail,
        Follow:     true,
    })
    if err != nil {
        gui.Log.Error(err)
        return err
    }
    defer readCloser.Close()

    pr, pw := io.Pipe()
    go func() {
        defer pw.Close()
        if ctr.Details.Config.Tty {
            io.Copy(pw, readCloser)
        } else {
            stdcopy.StdCopy(pw, pw, readCloser)
        }
    }()

    scanner := bufio.NewScanner(pr)
    index := 0
    for scanner.Scan() {
        select {
        case <-ctx.Done():
            return nil
        default:
            line := scanner.Text()
            // Format each line with a number prefix and separator
            formattedLine := fmt.Sprintf("[%d] ┃ %s", 
                index,
                line,
            )
			logsView.lines = append(logsView.lines, formattedLine)
			logsView.selected = len(logsView.lines) - 1 // Select the last line by default
            gui.g.Update(func(g *gocui.Gui) error {
                return gui.renderLogsView()
            })
            index++
        }
    }

    return scanner.Err()
}

func (gui *Gui) renderLogsView() error {
    logsView := gui.State.LogsView
	v := logsView.writer
    for i, line := range logsView.lines {
        if i == logsView.selected {
            fmt.Fprintf(*v, "\x1b[7m%s\x1b[0m\n", line)  // Only add highlighting during render
        } else {
            fmt.Fprintln(*v, line)
        }
    }
    return nil
}

func (gui *Gui) scrollUpLogs() error {
	logsView := gui.State.LogsView
	if logsView.selected > 0 {
		logsView.selected--
	}
	return gui.renderLogsView()
}
func (gui *Gui) scrollDownLogs() error {
	logsView := gui.State.LogsView
	if logsView.selected < len(logsView.lines)-1 {
		logsView.selected++
	}
	return gui.renderLogsView()
}
