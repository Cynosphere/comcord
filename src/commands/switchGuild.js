const {addCommand} = require("../lib/command");
const {startPrompt} = require("../lib/prompt");

const {listChannels} = require("./listChannels");
const {listUsers} = require("./listUsers");

function findTopChannel(guildId) {
  const guild = comcord.client.guilds.get(guildId);
  const channels = [...guild.channels.values()].filter((c) => c.type == 0);
  channels.sort((a, b) => a.position - b.position);

  return channels[0];
}

function switchGuild(input) {
  if (input == "") {
    listChannels();
    listUsers();
    return;
  }

  let target;

  for (const guild of comcord.client.guilds.values()) {
    if (guild.name.toLowerCase().indexOf(input.toLowerCase()) > -1) {
      target = guild.id;
      break;
    }
  }

  if (target == null) {
    console.log("<guild not found>");
  } else {
    comcord.state.currentGuild = target;
    // TODO: store last visited channel and switch to it if we've been to this guild before
    if (!comcord.state.lastChannel.has(target)) {
      const topChannel = findTopChannel(target);
      comcord.state.currentChannel = topChannel.id;
      comcord.state.lastChannel.set(target, topChannel.id);
    } else {
      comcord.state.currentChannel = comcord.state.lastChannel.get(target);
    }

    listChannels();
    listUsers();
  }
}

addCommand("G", "goto guild", function () {
  startPrompt(":guild> ", switchGuild);
});
