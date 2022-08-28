const chalk = require("chalk");

function processMessage({
  name,
  content,
  bot,
  attachments,
  reply,
  noColor = false,
}) {
  if (name.length + 2 > comcord.state.nameLength)
    comcord.statenameLength = name.length + 2;

  if (reply) {
    const nameColor = reply.author.bot ? chalk.bold.yellow : chalk.bold.cyan;

    const headerLength = 5 + reply.author.username.length;
    const length = headerLength + reply.content.length;

    if (noColor) {
      console.log(
        ` \u250d [${reply.author.username}] ${
          length > 79
            ? reply.content.substring(0, length - headerLength) + "\u2026"
            : reply.content
        }`
      );
    } else {
      console.log(
        chalk.bold.white(" \u250d ") +
          nameColor(`[${reply.author.username}] `) +
          chalk.reset(
            `${
              length > 79
                ? reply.content.substring(0, length - headerLength) + "\u2026"
                : reply.content
            }`
          )
      );
    }
  }

  if (
    (content.startsWith("*") && content.endsWith("*")) ||
    (content.startsWith("_") && content.endsWith("_"))
  ) {
    if (noColor) {
      console.log(`<${name} ${content.substring(1, content.length - 1)}>`);
    } else {
      console.log(
        chalk.bold.green(
          `<${name} ${content.substring(1, content.length - 1)}>`
        )
      );
    }
  } else {
    if (noColor) {
      console.log(
        `[${name}]${" ".repeat(
          Math.abs(comcord.state.nameLength - (name.length + 2))
        )} ${content}`
      );
    } else {
      const nameColor = bot ? chalk.bold.yellow : chalk.bold.cyan;

      // TODO: markdown
      console.log(
        nameColor(`[${name}]`) +
          " ".repeat(Math.abs(comcord.state.nameLength - (name.length + 2))) +
          chalk.reset(" " + content)
      );
    }
  }

  if (attachments) {
    for (const attachment of attachments) {
      if (noColor) {
        console.log(`<attachment: ${attachment.url} >`);
      } else {
        console.log(chalk.bold.yellow(`<attachment: ${attachment.url} >`));
      }
    }
  }
}

function processQueue() {
  for (const msg of comcord.state.messageQueue) {
    if (msg.content.indexOf("\n") > -1) {
      const lines = msg.content.split("\n");
      for (const index in lines) {
        const line = lines[index];
        processMessage({
          name: msg.author.username,
          bot: msg.author.bot,
          content: line,
          attachments: index == lines.length - 1 ? msg.attachments : [],
          reply: index == 0 ? msg.referencedMessage : null,
        });
      }
    } else {
      processMessage({
        name: msg.author.username,
        bot: msg.author.bot,
        content: msg.content,
        attachments: msg.attachments,
        reply: msg.referencedMessage,
      });
    }
  }

  comcord.state.messageQueue.splice(0, comcord.state.messageQueue.length);
}

module.exports = {
  processMessage,
  processQueue,
};
