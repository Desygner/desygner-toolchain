package extension

import (
	"context"
	"fmt"
	"github.com/desygner/desygner-toolchain/packages/slack/facades"
	"github.com/slack-go/slack"
	"strconv"
	"strings"
)

type Input struct {
	SlackChannel, SlackTS string

	shouldReplace bool
}

func (i *Input) SlackTimestamp() string {
	if i.SlackTS == "" {
		return ""
	}
	if strings.Contains(i.SlackTS, ".") {
		return i.SlackTS
	}
	raw, _ := strconv.Atoi(i.SlackTS)
	seconds := float64(raw) / 1_000_000
	return fmt.Sprintf("%.6f\n", seconds)
}
func (i *Input) SendMessage(ctx context.Context, msg string) error {
	if i == nil || i.SlackChannel == "" || i.SlackTimestamp() == "" {
		return nil
	}
	cli := facades.Slack().Client()
	if !i.shouldReplace {

		ch, ts, _, err := cli.SendMessageContext(ctx, i.SlackChannel,
			slack.MsgOptionText(msg, false),
			slack.MsgOptionTS(i.SlackTimestamp()),
		)
		if err != nil {
			return err
		}
		i.SlackTS = ts
		i.SlackChannel = ch
		i.shouldReplace = true
		return nil
	}
	_, _, _, err := cli.UpdateMessageContext(ctx, i.SlackChannel, i.SlackTS, slack.MsgOptionText(msg, false))
	if err != nil {
		return err
	}
	return nil
}
