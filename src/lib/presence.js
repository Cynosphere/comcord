const CLIENT_ID = "1026163285877325874";

function updatePresence() {
  let guild, channel;
  if (comcord.state.currentGuild != null) {
    guild = comcord.client.guilds.get(comcord.state.currentGuild);
  }
  if (comcord.state.currentChannel != null && guild != null) {
    channel = guild.channels.get(comcord.state.currentChannel);
  }

  if (comcord.client.user.bot) {
    if (comcord.state.rpcConnected) {
      try {
        const activity = {
          startTimestamp: comcord.state.startTime,
          smallImageKey: `https://cdn.discordapp.com/avatars/${comcord.client.user.id}/${comcord.client.user.avatar}.png?size=1024`,
          smallImageText: comcord.client.user.username,
          buttons: [
            {
              label: "comcord Repo",
              url: "https://github.com/Cynosphere/comcord",
            },
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

    comcord.client.editStatus(
      comcord.state.afk ? "idle" : comcord.config.defaultStatus ?? "online",
      [
        {
          name: "comcord",
          type: 0,
          application_id: CLIENT_ID,
          timestamps: {
            start: comcord.state.startTime,
          },
        },
      ]
    );
  } else {
    const activity = {
      application_id: CLIENT_ID,
      name: "comcord",
      timestamps: {
        start: comcord.state.startTime,
      },
      assets: {},
      buttons: ["comcord Repo"],
      metadata: {
        button_urls: ["https://github.com/Cynosphere/comcord"],
      },
      type: 0,
    };

    if (guild != null) {
      activity.assets.large_image = `mp:icons/${guild.id}/${guild.icon}.png?size=1024`;
      activity.assets.large_text = guild.name;
      if (channel != null) {
        activity.details = `#${channel.name} - ${guild.name}`;
      }
    }
    if (comcord.state.afk == true) {
      activity.state = "AFK";
    }

    comcord.client.editStatus(
      comcord.state.afk ? "idle" : comcord.config.defaultStatus ?? "online",
      [activity]
    );
  }
}

module.exports = {updatePresence};
