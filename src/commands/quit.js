const {addCommand} = require("../lib/command");

addCommand("q", "quit comcord", function () {
  comcord.client.disconnect(false);
  process.exit(0);
});
