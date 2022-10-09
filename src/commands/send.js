const chalk = require("chalk");

const {startPrompt} = require("../lib/prompt");
const {updatePresence} = require("../lib/presence");

function sendMode() {
  if (!comcord.state.currentChannel) {
    console.log("<not in a channel>");
    return;
  }

  startPrompt(
    chalk.bold.cyan(`[${comcord.client.user.username}]`) +
      " ".repeat(
        comcord.state.nameLength - (comcord.client.user.username.length + 2)
      ) +
      chalk.reset(" "),
    async function (input) {
      if (input == "") {
        console.log("<no message sent>");
      } else {
        try {
          process.stdout.write("\n");
          await comcord.client.guilds
            .get(comcord.state.currentGuild)
            .channels.get(comcord.state.currentChannel)
            .createMessage({content: input});

          if (comcord.state.afk == true) {
            comcord.state.afk = false;
            comcord.client.editStatus("online");
            comcord.client.editAFK(false);
            console.log("<you have returned>");

            updatePresence();
          }
        } catch (err) {
          console.log("<failed to send message: " + err.message + ">");
        }
      }
    }
  );
}

module.exports = {sendMode};
