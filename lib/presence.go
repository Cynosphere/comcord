package lib

import (
	"context"
	"fmt"

	"github.com/Cynosphere/comcord/state"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
)

func UpdatePresence() {
  client := state.GetClient()

  self, err := client.MeStore.Me()
  if err != nil {
    return
  }

  afk := state.IsAFK()
  presence := gateway.UpdatePresenceCommand{
    Since: 0,
    Activities: make([]discord.Activity, 0),
    AFK: false,
  }

  currentGuild := state.GetCurrentGuild()
  currentChannel := state.GetCurrentChannel()

  parsedGuildId, err := discord.ParseSnowflake(currentGuild)
  if err != nil {
    return
  }
  parsedChannelId, err := discord.ParseSnowflake(currentChannel)
  if err != nil {
    return
  }

  var activity discord.Activity

  startTime := state.GetStartTime()

  if self.Bot {
    activity = discord.Activity{
      Type: discord.GameActivity,
      Name: "comcord",
    }

    if currentGuild != "" && currentChannel != "" {
      guild, guildErr := client.GuildStore.Guild(discord.GuildID(parsedGuildId))
      channel, channelErr := client.ChannelStore.Channel(discord.ChannelID(parsedChannelId))

      if guildErr == nil && channelErr == nil {
        activity.Type = discord.WatchingActivity
        activity.Name = fmt.Sprintf("#%s in %s | comcord", channel.Name, guild.Name)
      }
    }

    if afk {
      activity.Name = activity.Name + " [AFK]"
    }
  } else {
    parsedAppId, err := discord.ParseSnowflake("1026163285877325874")
    if err != nil {
      return
    }

    activity = discord.Activity{
      Type: 0,
      AppID: discord.AppID(parsedAppId),
      Name: "comcord",
      Timestamps: &discord.ActivityTimestamps{
        Start: discord.UnixMsTimestamp(startTime.Unix()),
      },
      /*Buttons: make([]string, 0),
      Metadata: ActivityMetadata{
        ButtonURLs: make([]string, 0),
      },*/
    }

    //activity.Buttons = append(activity.Buttons, "comcord Repo")
    //activity.Metadata.ButtonURLs = append(activity.Metadata.ButtonURLs, "https://gitdab.com/Cynosphere/comcord")

    if currentGuild != "" && currentChannel != "" {
      guild, guildErr := client.GuildStore.Guild(discord.GuildID(parsedGuildId))
      channel, channelErr := client.ChannelStore.Channel(discord.ChannelID(parsedChannelId))

      if guildErr == nil && channelErr == nil {
        activity.Details = fmt.Sprintf("#%s - %s", channel.Name, guild.Name)

        activity.Assets = &discord.ActivityAssets{}
        activity.Assets.LargeText = guild.Name
        if guild.Icon != "" {
          activity.Assets.LargeImage = fmt.Sprintf("mp:icons/%s/%s.png?size=1024", guild.ID, guild.Icon)
        }
      }
    }

    if afk {
      activity.State = "AFK"
    }
  }

  activity.CreatedAt = discord.UnixTimestamp(startTime.Unix())

  presence.Activities = append(presence.Activities, activity)

  defaultStatus := state.GetConfigValue("defaultStatus")
  if defaultStatus != "" {
    presence.Status = discord.Status(defaultStatus)
  } else {
    if afk {
      presence.Status = discord.IdleStatus
    } else {
      presence.Status = discord.OnlineStatus
    }
  }

  client.Gateway().Send(context.Background(), &presence)
}
