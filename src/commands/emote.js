const {addCommand} = require("../lib/command");
const {startPrompt} = require("../lib/prompt");

addCommand("e", "emote", function () {
  if (!comcord.state.currentChannel) {
    console.log("<not in a channel>");
    return;
  }

  startPrompt(":emote> ", async function (input) {
    if (input == "") {
      console.log("<no message sent>");
    } else {
      try {
        process.stdout.write("\n");
        await comcord.client.createMessage(
          comcord.state.currentChannel,
          `*${input}*`
        );
        console.log(`<${comcord.client.user.username} ${input}>`);
      } catch (err) {
        console.log("<failed to send message: " + err.message + ">");
      }
    }
  });
});
