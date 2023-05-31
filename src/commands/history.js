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
    processMessage(msg, {noColor: true, history: true});
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
    console.log(`<failed to get history: ${err.message}>`);
  }
}

addCommand("r", "channel history", getHistory);
addCommand("R", "extended history", function () {
  startPrompt(":lines> ", async function (input) {
    process.stdout.write("\n");
    await getExtendedHistory(input);
  });
});
