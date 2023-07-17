package lib

import (
	"fmt"
	"reflect"
	"sync"
	"time"
	"unsafe"

	"github.com/Cynosphere/comcord/state"
	"github.com/bwmarrin/discordgo"
	"github.com/gorilla/websocket"
)

type ActivityMetadata struct {
  ButtonURLs []string `json:"button_urls,omitempty"`
}

type Activity struct {
  Name          string                 `json:"name"`
	Type          discordgo.ActivityType `json:"type"`
	URL           string                 `json:"url,omitempty"`
	CreatedAt     time.Time              `json:"created_at"`
	ApplicationID string                 `json:"application_id,omitempty"`
	State         string                 `json:"state,omitempty"`
	Details       string                 `json:"details,omitempty"`
	Timestamps    discordgo.TimeStamps   `json:"timestamps,omitempty"`
	Emoji         discordgo.Emoji        `json:"emoji,omitempty"`
	Party         discordgo.Party        `json:"party,omitempty"`
	Assets        discordgo.Assets       `json:"assets,omitempty"`
	Secrets       discordgo.Secrets      `json:"secrets,omitempty"`
	Instance      bool                   `json:"instance,omitempty"`
	Flags         int                    `json:"flags,omitempty"`
  Buttons       []string               `json:"buttons,omitempty"`
  Metadata      ActivityMetadata       `json:"metadata,omitempty"`
}

type GatewayPresenceUpdate struct {
  Since      int        `json:"since"`
  Activities []Activity `json:"activities,omitempty"`
  Status     string     `json:"status"`
  AFK        bool       `json:"afk"`
  Broadcast  string     `json:"broadcast,omitempty"`
}

type presenceOp struct {
  Op int                     `json:"op"`
  Data GatewayPresenceUpdate `json:"d"`
}

func getUnexportedField(field reflect.Value) interface{} {
  return reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem().Interface()
}

func UpdatePresence(session *discordgo.Session) {
  values := reflect.ValueOf(session)
  fieldWsConn := reflect.Indirect(values).FieldByName("wsConn")
  fieldWsMutex := reflect.Indirect(values).FieldByName("wsMutex")

  wsConn := getUnexportedField(fieldWsConn).(*websocket.Conn)
  wsMutex := getUnexportedField(fieldWsMutex).(sync.Mutex)

  afk := state.IsAFK()
  presence := GatewayPresenceUpdate{
    Since: 0,
    AFK: afk,
    Activities: make([]Activity, 0),
  }

  currentGuild := state.GetCurrentGuild()
  currentChannel := state.GetCurrentChannel()

  var activity Activity

  startTime := state.GetStartTime()

  if session.State.User.Bot {
    activity = Activity{
      Type: 0,
      Name: "comcord",
    }

    if currentGuild != "" && currentChannel != "" {
      guild, guildErr := session.State.Guild(currentGuild)
      channel, channelErr := session.State.Channel(currentChannel)

      if guildErr == nil && channelErr == nil {
        activity.Type = 3
        activity.Name = fmt.Sprintf("#%s in %s | comcord", channel.Name, guild.Name)
      }
    }

    if afk {
      activity.Name = activity.Name + " [AFK]"
    }
  } else {
    activity = Activity{
      Type: 0,
      ApplicationID: "1026163285877325874",
      Name: "comcord",
      Timestamps: discordgo.TimeStamps{
        StartTimestamp: startTime.Unix(),
      },
      Buttons: make([]string, 0),
      Metadata: ActivityMetadata{
        ButtonURLs: make([]string, 0),
      },
    }

    activity.Buttons = append(activity.Buttons, "comcord Repo")
    activity.Metadata.ButtonURLs = append(activity.Metadata.ButtonURLs, "https://gitdab.com/Cynosphere/comcord")

    if currentGuild != "" && currentChannel != "" {
      guild, guildErr := session.State.Guild(currentGuild)
      channel, channelErr := session.State.Channel(currentChannel)

      if guildErr == nil && channelErr == nil {
        activity.Details = fmt.Sprintf("#%s - %s", channel.Name, guild.Name)

        activity.Assets = discordgo.Assets{}
        activity.Assets.LargeText = guild.Name
        if guild.Icon != "" {
          activity.Assets.LargeImageID = fmt.Sprintf("mp:icons/%s/%s.png?size=1024", guild.ID, guild.Icon)
        }
      }
    }

    if afk {
      activity.State = "AFK"
    }
  }

  activity.CreatedAt = startTime

  presence.Activities = append(presence.Activities, activity)

  defaultStatus := state.GetConfigValue("defaultStatus")
  if defaultStatus != "" {
    presence.Status = defaultStatus
  } else {
    if afk {
      presence.Status = "idle"
    } else {
      presence.Status = "online"
    }
  }

  op := presenceOp{3, presence}
  wsMutex.Lock()
  wsConn.WriteJSON(op)
  wsMutex.Unlock()
}
