function updatePresence() {
  let guild, channel;
  if (comcord.state.currentGuild != null) {
    guild = comcord.client.guilds.get(comcord.state.currentGuild);
  }
  if (comcord.state.currentChannel != null && guild != null) {
    channel = guild.channels.get(comcord.state.currentChannel);
  }

  try {
    const activity = {
      startTimestamp: comcord.state.startTime,
      smallImageKey: `https://cdn.discordapp.com/avatars/${comcord.client.user.id}/${comcord.client.user.avatar}.png?size=1024`,
      smallImageText: `${comcord.client.user.username}#${comcord.client.user.discriminator}`,
      buttons: [
        {label: "comcord Repo", url: "https://github.com/Cynosphere/comcord"},
      ],
    };

    if (guild != null) {
      activity.largeImageKey = `https://cdn.discordapp.com/icons/${guild.id}/${guild.icon}.png?size=1024`;
      activity.largeImageText = guild.name;
      if (channel != null) {
        activity.details = `#${channel.name} - ${guild.name}`;
      }
    }
    if (comcord.state.afk == true) {
      activity.state = "AFK";
    }
    comcord.rpc.setActivity(activity);
  } catch (err) {
    //
  }
}

module.exports = {updatePresence};
