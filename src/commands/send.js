const chalk = require("chalk");

const {startPrompt} = require("../lib/prompt");

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
          await comcord.client.createMessage(
            comcord.state.currentChannel,
            input
          );

          if (comcord.state.afk == true) {
            comcord.state.afk = false;
            comcord.client.editStatus("online");
            comcord.client.editAFK(false);
            console.log("<you have returned>");
          }
        } catch (err) {
          console.log("<failed to send message: " + err.message + ">");
        }
      }
    }
  );
}

module.exports = {sendMode};
