package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"

	"github.com/eternal-flame-AD/gotify-netlify/api"

	"github.com/gotify/plugin-api"

	"github.com/gin-gonic/gin"

	"gopkg.in/square/go-jose.v2/jwt"
)

// GetGotifyPluginInfo returns gotify plugin info.
func GetGotifyPluginInfo() plugin.Info {
	return plugin.Info{
		ModulePath: "github.com/eternal-flame-AD/gotify-netlify",
		Name:       "Netlify Webhook Plugin",
	}
}

// Conf is configuration
type Conf struct {
	SecretToken string
}

// Plugin is plugin instance
type Plugin struct {
	userCtx    plugin.UserContext
	conf       *Conf
	msgHandler plugin.MessageHandler
	basePath   string

	enabled bool
}

// Enable implements plugin.Plugin
func (c *Plugin) Enable() error {
	c.enabled = true
	urlBase, _ := url.Parse(c.basePath)
	c.notify("Netlify webhook plugin started", "Listening on "+urlBase.ResolveReference(&url.URL{Path: "hook"}).String(), 0)
	return nil
}

// Disable implements plugin.Plugin
func (c *Plugin) Disable() error {
	c.enabled = false
	return nil
}

func (c *Plugin) notify(title string, msg string, priority int) error {
	return c.msgHandler.SendMessage(plugin.Message{
		Title:    title,
		Message:  msg,
		Priority: priority,
	})
}

func (c *Plugin) notifyMsg(m *api.WebhookMsg) error {
	priority := 0
	title := "Netlify deploy "
	switch m.State {
	case api.StateReady:
		title += "succeeded"
		priority = 1
	case api.StateError:
		title += "errored"
		priority = 1
	case api.StateBuilding:
		title += "started"
	default:
		title += string(m.State)
	}
	body := bytes.NewBuffer([]byte{})
	body.WriteString(m.Title)
	body.WriteRune('\n')
	switch m.Context {
	case api.ContextBranchDeploy:
		body.WriteString("branch deploy on ")
		body.WriteString(m.Branch)
	case api.ContextProduction:
		body.WriteString("production deploy")
	}
	body.WriteString(fmt.Sprintf(" for %s\n", m.SiteName))
	if m.DeployTime != nil && *m.DeployTime != 0 {
		body.WriteString(fmt.Sprintf("Deployed in %d seconds\n", *m.DeployTime))
	}
	body.WriteString(fmt.Sprintf("Commit URL: %s\n", m.CommitURL))
	body.WriteString(fmt.Sprintf("Site updated at %s\n", m.UpdatedAt.Format("Jan 2 15:04:05")))
	body.WriteString(fmt.Sprintf("Site alive at %s\n", m.URL))
	body.WriteString(m.Summary.String())
	return c.notify(title, body.String(), priority)
}

// GetDisplay implements plugin.Displayer
func (c *Plugin) GetDisplay(baseURL *url.URL) string {
	doc := bytes.NewBuffer([]byte{})
	doc.WriteString("# Netlify Webhook Plugin\n\n")

	baseURL.Path = c.basePath
	hookURL := &url.URL{
		Path: "hook",
	}
	hookURL = baseURL.ResolveReference(hookURL)
	doc.WriteString(fmt.Sprintf("Webhook URI: %s\n\n", hookURL))

	if !c.enabled {
		doc.WriteString("**Warning**: plugin disabled, not relaying messages\n\n")
	}
	if c.conf.SecretToken == "" {
		doc.WriteString("**Warning**: secret token not set, webhooks are not verified.\n\n")
	}
	return doc.String()
}

// SetMessageHandler implements plugin.Messenger
func (c *Plugin) SetMessageHandler(h plugin.MessageHandler) {
	c.msgHandler = h
}

// RegisterWebhook is called to register handlers for the plugin
func (c *Plugin) RegisterWebhook(baseURL string, r *gin.RouterGroup) {
	c.basePath = baseURL
	r.POST("/hook", func(ctx *gin.Context) {

		jsonData, err := ioutil.ReadAll(ctx.Request.Body)
		if err != nil {
			log.Println(err)
			ctx.AbortWithStatusJSON(500, err.Error())
			return
		}

		// Check JWT
		if c.conf.SecretToken != "" {
			jws := ctx.Request.Header.Get("X-Webhook-Signature")
			signature, err := jwt.ParseSigned(jws)
			if err != nil {
				log.Println(err)
				ctx.AbortWithStatusJSON(401, err.Error())
				return
			}
			sigData := struct {
				ISS    string `json:"iss"`
				SHA256 string `json:"sha256"`
			}{}
			if err := signature.Claims([]byte(c.conf.SecretToken), &sigData); err != nil {
				log.Println(err)
				ctx.AbortWithStatusJSON(401, err.Error())
				return
			}
			if sigData.ISS != "netlify" {
				ctx.AbortWithStatusJSON(401, "iss is not netlify")
				return
			}
			if sha256Str(jsonData) != sigData.SHA256 {
				ctx.AbortWithStatusJSON(401, "sha256 does not match")
				return
			}
		}

		// Unmarshal JSON
		res := new(api.WebhookMsg)
		if err := json.Unmarshal(jsonData, res); err != nil {
			log.Println(err)
			ctx.AbortWithStatusJSON(500, err.Error())
			return
		}

		// Make internal request to create gotify message
		if err := c.notifyMsg(res); err != nil {
			log.Println(err)
			ctx.AbortWithStatusJSON(500, err.Error())
			return
		}
		ctx.Status(200)
	})
}

// NewGotifyPluginInstance creates a plugin instance for a user context.
func NewGotifyPluginInstance(ctx plugin.UserContext) plugin.Plugin {
	return &Plugin{
		userCtx: ctx,
	}
}

func main() {
	panic("this should be built as a go plugin")
}
