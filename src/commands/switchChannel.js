const {addCommand} = require("../lib/command");
const {startPrompt} = require("../lib/prompt");

const {listUsers} = require("./listUsers");

function switchChannel(input) {
  if (input == "") {
    listUsers();
    comcord.state.channelSwitch = false;
    return;
  }
  let target;

  const guild = comcord.client.guilds.get(comcord.state.currentGuild);
  const channels = [...guild.channels.values()].filter((c) => c.type == 0);
  channels.sort((a, b) => a.position - b.position);

  for (const channel of channels) {
    if (channel.name.toLowerCase().indexOf(input.toLowerCase()) > -1) {
      target = channel.id;
      break;
    }
  }

  if (target == null) {
    console.log("<channel not found>");
  } else {
    comcord.state.currentChannel = target;
    comcord.state.lastChannel.set(comcord.state.currentGuild, target);

    listUsers();
  }
}

addCommand("g", "goto channel", function () {
  if (!comcord.state.currentGuild) {
    console.log("<not in a guild>");
    return;
  }
  startPrompt(":channel> ", switchChannel);
});
