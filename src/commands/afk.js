const {addCommand} = require("../lib/command");
const {updatePresence} = require("../lib/presence");

addCommand("A", "toggles AFK mode", function () {
  if (comcord.state.afk == true) {
    comcord.client.shards.forEach((shard) => (shard.presence.afk = false));
    comcord.client.editStatus("online");
    console.log("<you have returned>");
  } else {
    comcord.state.afk = true;
    comcord.client.shards.forEach((shard) => (shard.presence.afk = true));
    comcord.client.editStatus("idle");
    console.log("<you go AFK>");
  }

  updatePresence();
});
