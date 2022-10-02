const {addCommand} = require("../lib/command");

addCommand("A", "toggles AFK mode", function () {
  if (comcord.state.afk == true) {
    comcord.state.afk = false;
    comcord.client.editStatus("online");
    comcord.client.editAFK(false);
    console.log("<you have returned>");
  } else {
    comcord.state.afk = true;
    comcord.client.editStatus("idle");
    comcord.client.editAFK(true);
    console.log("<you go AFK>");
  }
});
