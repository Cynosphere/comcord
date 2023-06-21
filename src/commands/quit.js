const {addCommand} = require("../lib/command");

addCommand("q", "quit comcord", function () {
  comcord.client.disconnect({reconnect: false});
  process.exit(0);
});
