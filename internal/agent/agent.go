package agent

import (
	"context"
	"fmt"
	"net/http"
	"time"

	socketio "github.com/zishang520/socket.io/clients/socket/v3"
	"github.com/zishang520/socket.io/v3/pkg/types"
)

type Agent struct {
	config Config

	socket *socketio.Socket

	httpClient *http.Client
}

func New(config Config) *Agent {
	return &Agent{
		config: config,

		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (a *Agent) Run(ctx context.Context) error {
	fmt.Println()
	fmt.Println("Tunnexo Agent")
	fmt.Println("================================")

	if a.config.Token == "" {
		fmt.Println("Connecting as guest...")
		token, err := requestGuestToken(ctx, a.config.ServerURL, &http.Client{})
		if err != nil {
			return err
		}
		a.config.Token = token
	}

	options := socketio.DefaultOptions()

	options.SetTransports(
		types.NewSet(
			socketio.WebSocket,
		),
	)

	options.SetAuth(
		map[string]any{
			"token":     a.config.Token,
			"agentName": a.config.AgentName,
			"agentIp":   a.config.AgentIP,
		},
	)

	socket, err := socketio.Connect(
		a.config.ServerURL,
		options,
	)

	if err != nil {
		return fmt.Errorf(
			"failed to connect to Tunnexo server",
		)
	}

	a.socket = socket

	a.registerEvents()

	<-ctx.Done()

	fmt.Println()
	fmt.Println("Stopping Tunnexo Agent...")

	a.socket.Disconnect()

	fmt.Println("Tunnexo Agent stopped.")

	return nil
}

func (a *Agent) registerEvents() {

	_ = a.socket.On(
		"connect",
		func(args ...any) {

			fmt.Println("✓ Connected to Tunnexo Server")

			fmt.Println("Connection ready")
		},
	)

	_ = a.socket.On(
		"authenticated",
		func(args ...any) {

			if len(args) == 0 {
				fmt.Println("✗ Missing authentication response")
				return
			}
			data, ok := args[0].(map[string]any)
			if !ok || data["ok"] != true {
				fmt.Println("✗ Agent authentication failed")
				return
			}
			fmt.Println("✓ Agent authenticated")

			a.createTunnel()
		},
	)

	_ = a.socket.On(
		"connect_error",
		func(args ...any) {

			fmt.Println("✗ Connection failed")
		},
	)

	_ = a.socket.On(
		"error",
		func(args ...any) {

			fmt.Println("✗ Server error")
		},
	)

	_ = a.socket.On(
		"disconnect",
		func(args ...any) {

			fmt.Println("⚠ Disconnected")
		},
	)

	_ = a.socket.On(
		"http:request",
		func(args ...any) {

			if len(args) == 0 {
				fmt.Println(
					"Invalid http:request: empty payload",
				)

				return
			}

			payload, ok :=
				args[0].(map[string]any)

			if !ok {

				fmt.Println("Invalid http:request payload")

				return
			}

			go a.handleHTTPRequest(payload)
		},
	)
}

func (a *Agent) createTunnel() {

	payload := TunnelCreateRequest{
		Name: a.config.TunnelName,
	}

	fmt.Println()
	fmt.Println("Creating tunnel...")

	a.socket.
		Timeout(10*time.Second).
		EmitWithAck(
			"tunnel:create",
			payload,
		)(
		func(args []any, err error) {

			if err != nil {

				fmt.Println("✗ Tunnel creation failed")

				return
			}

			if len(args) == 0 {

				fmt.Println(
					"✗ Tunnel server returned empty response",
				)

				return
			}

			response,
				ok := args[0].(map[string]any)

			if !ok {

				fmt.Println("✗ Invalid tunnel response")

				return
			}

			success,
				_ := response["ok"].(bool)

			if !success {

				fmt.Println("✗ Tunnel creation failed")

				return
			}

			publicURL, reason := validatedPublicURL(response["url"], a.config.Token, a.config.ServerURL, a.config.Environment)
			if reason != "" {
				fmt.Println("✗ Public URL rejected:", reason)
				return
			}
			fmt.Println()
			fmt.Println("✓ Tunnel Established!")
			fmt.Println("================================")
			fmt.Println("Public URL:", publicURL)

			fmt.Println("================================")
			fmt.Println()
		},
	)
}
