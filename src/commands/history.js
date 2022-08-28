const {addCommand} = require("../lib/command");
const {startPrompt} = require("../lib/prompt");
const {processMessage} = require("../lib/messages");

async function getHistory(limit = 20) {
  if (!comcord.state.currentChannel) {
    console.log("<not in a channel>");
    return;
  }

  const messages = await comcord.client.getMessages(
    comcord.state.currentChannel,
    {limit}
  );
  messages.reverse();

  console.log("--Beginning-Review".padEnd(72, "-"));

  for (const msg of messages) {
    if (msg.content.indexOf("\n") > -1) {
      const lines = msg.content.split("\n");
      for (const index in lines) {
        const line = lines[index];
        processMessage({
          name: msg.author.username,
          bot: msg.author.bot,
          content:
            line +
            (msg.editedTimestamp != null && index == lines.length - 1
              ? " (edited)"
              : ""),
          attachments: index == lines.length - 1 ? msg.attachments : null,
          reply: index == 0 ? msg.referencedMessage : null,
          noColor: true,
        });
      }
    } else {
      processMessage({
        name: msg.author.username,
        bot: msg.author.bot,
        content: msg.content + (msg.editedTimestamp != null ? " (edited)" : ""),
        attachments: msg.attachments,
        reply: msg.referencedMessage,
        noColor: true,
      });
    }
  }

  console.log("--Review-Complete".padEnd(73, "-"));
}

async function getExtendedHistory(input) {
  input = parseInt(input);
  if (isNaN(input)) {
    console.log("<not a number>");
    return;
  }

  try {
    await getHistory(input);
  } catch (err) {
    console.log("<failed to get history: " + err.message + ">");
  }
}

addCommand("r", "channel history", getHistory);
addCommand("R", "extended history", function () {
  startPrompt(":lines> ", async function (input) {
    process.stdout.write("\n");
    await getExtendedHistory(input);
  });
});
