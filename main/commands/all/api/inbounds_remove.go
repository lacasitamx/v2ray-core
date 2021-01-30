package api

import (
	"fmt"

	handlerService "github.com/v2fly/v2ray-core/v4/app/proxyman/command"
	"github.com/v2fly/v2ray-core/v4/common/cmdarg"
	"github.com/v2fly/v2ray-core/v4/infra/conf/serial"
	"github.com/v2fly/v2ray-core/v4/main/commands/base"
)

var cmdRemoveInbounds = &base.Command{
	CustomFlags: true,
	UsageLine:   "{{.Exec}} api rmi [--server=127.0.0.1:8080] <json_file|tag> [json_file] [tag]...",
	Short:       "Remove inbounds",
	Long: `
Remove inbounds from V2Ray.

Arguments:

	-s, -server 
		The API server address. Default 127.0.0.1:8080

	-t, -timeout
		Timeout seconds to call API. Default 3

Example:

    {{.Exec}} {{.LongName}} --server=127.0.0.1:8080 c1.json "tag name"
`,
	Run: executeRemoveInbounds,
}

func executeRemoveInbounds(cmd *base.Command, args []string) {
	setSharedFlags(cmd)
	cmd.Flag.Parse(args)
	unnamedArgs := cmd.Flag.Args()
	if len(unnamedArgs) == 0 {
		fmt.Println("reading from stdin:")
		unnamedArgs = []string{"stdin:"}
	}

	tags := make([]string, 0)
	for _, arg := range unnamedArgs {
		if r, err := cmdarg.LoadArg(arg); err == nil {
			conf, err := serial.DecodeJSONConfig(r)
			if err != nil {
				base.Fatalf("failed to decode %s: %s", arg, err)
			}
			ins := conf.InboundConfigs
			for _, i := range ins {
				tags = append(tags, i.Tag)
			}
		} else {
			// take request as tag
			tags = append(tags, arg)
		}
	}

	if len(tags) == 0 {
		base.Fatalf("no inbound to remove")
	}
	fmt.Println("removing inbounds:", tags)

	conn, ctx, close := dialAPIServer()
	defer close()

	client := handlerService.NewHandlerServiceClient(conn)
	for _, tag := range tags {
		fmt.Println("removing:", tag)
		r := &handlerService.RemoveInboundRequest{
			Tag: tag,
		}
		resp, err := client.RemoveInbound(ctx, r)
		if err != nil {
			base.Fatalf("failed to remove inbound: %s", err)
		}
		showResponese(resp)
	}
}