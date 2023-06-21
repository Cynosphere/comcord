const {addCommand} = require("../lib/command");
const {startPrompt} = require("../lib/prompt");
const {processMessage} = require("../lib/messages");
const {listChannels} = require("./listChannels");

async function getHistory(limit = 20, channel = null) {
  if (!channel && !comcord.state.currentChannel) {
    console.log("<not in a channel>");
    return;
  }

  const messages = await comcord.client.getMessages(
    channel ?? comcord.state.currentChannel
  );
  messages.reverse();

  console.log("--Beginning-Review".padEnd(72, "-"));

  const lines = [];
  for (const msg of messages) {
    const processedLines = processMessage(msg, {noColor: true, history: true});
    if (processedLines) lines.push(...processedLines);
  }
  console.log(lines.slice(-limit).join("\n"));

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
addCommand("p", "peek at channel", function () {
  if (!comcord.state.currentGuild) {
    console.log("<not in a guild>");
    return;
  }

  listChannels();
  startPrompt(":peek> ", async function (input) {
    console.log("");
    if (input == "") {
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
      await getHistory(20, target);
    }
  });
});
