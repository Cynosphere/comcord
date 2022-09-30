const chalk = require("chalk");

const REGEX_MENTION = /<@!?(\d+)>/g;
const REGEX_ROLE_MENTION = /<@&?(\d+)>/g;
const REGEX_CHANNEL = /<#(\d+)>/g;
const REGEX_EMOTE = /<(?:\u200b|&)?a?:(\w+):(\d+)>/g;
const REGEX_COMMAND = /<\/([^\s]+?):(\d+)>/g;

function replaceMentions(_, id) {
  const user = comcord.client.users.get(id);
  if (user) {
    return `@${user.username}`;
  } else {
    return "@Unknown User";
  }
}
function replaceRoles(_, id) {
  const role = comcord.client.guilds
    .get(comcord.state.currentGuild)
    .roles.get(id);
  if (role) {
    return `[@${role.name}]`;
  } else {
    return "[@Unknown Role]";
  }
}
function replaceChannels(_, id) {
  const guildForChannel = comcord.client.channelGuildMap[id];
  if (guildForChannel) {
    const channel = comcord.client.guilds.get(guildForChannel).channels.get(id);
    if (channel) {
      return `#${channel.name}`;
    } else {
      return "#unknown-channel";
    }
  } else {
    return "#unknown-channel";
  }
}
function replaceEmotes(_, name, id) {
  return `:${name}:`;
}
function replaceCommands(_, name, id) {
  return `/${name}`;
}

function processMessage({
  name,
  content,
  bot,
  attachments,
  stickers,
  reply,
  noColor = false,
}) {
  if (name.length + 2 > comcord.state.nameLength)
    comcord.state.nameLength = name.length + 2;

  if (reply) {
    const nameColor = reply.author.bot ? chalk.bold.yellow : chalk.bold.cyan;

    const headerLength = 5 + reply.author.username.length;
    const length = headerLength + reply.content.length;

    let replyContent = reply.content.replace(/\n/g, " ");
    replyContent = replyContent
      .replace(REGEX_MENTION, replaceMentions)
      .replace(REGEX_ROLE_MENTION, replaceRoles)
      .replace(REGEX_CHANNEL, replaceChannels)
      .replace(REGEX_EMOTE, replaceEmotes)
      .replace(REGEX_COMMAND, replaceCommands);

    if (noColor) {
      console.log(
        ` \u250d [${reply.author.username}] ${
          length > 79
            ? replyContent.substring(0, 79 - headerLength) + "\u2026"
            : replyContent
        }`
      );
    } else {
      console.log(
        chalk.bold.white(" \u250d ") +
          nameColor(`[${reply.author.username}] `) +
          chalk.reset(
            `${
              length > 79
                ? replyContent.substring(0, 79 - headerLength) + "\u2026"
                : replyContent
            }`
          )
      );
    }
  }

  content = content
    .replace(REGEX_MENTION, replaceMentions)
    .replace(REGEX_ROLE_MENTION, replaceRoles)
    .replace(REGEX_CHANNEL, replaceChannels)
    .replace(REGEX_EMOTE, replaceEmotes)
    .replace(REGEX_COMMAND, replaceCommands);

  if (
    (content.length > 1 && content.startsWith("*") && content.endsWith("*")) ||
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

  if (stickers) {
    for (const sticker of stickers) {
      if (noColor) {
        console.log(
          `<sticker: "${sticker.name}" https://media.discordapp.net/stickers/${sticker.id}.png >`
        );
      } else {
        console.log(
          chalk.bold.yellow(
            `<sticker: "${sticker.name}" https://media.discordapp.net/stickers/${sticker.id}.png >`
          )
        );
      }
    }
  }
}

function processQueue() {
  for (const msg of comcord.state.messageQueue) {
    if (msg.time) {
      console.log(msg.content);
    } else if (msg.content.indexOf("\n") > -1) {
      const lines = msg.content.split("\n");
      for (const index in lines) {
        const line = lines[index];
        processMessage({
          name: msg.author.username,
          bot: msg.author.bot,
          content: line,
          attachments: index == lines.length - 1 ? msg.attachments : [],
          stickers: index == lines.length - 1 ? msg.stickerItems : [],
          reply: index == 0 ? msg.referencedMessage : null,
        });
      }
    } else {
      processMessage({
        name: msg.author.username,
        bot: msg.author.bot,
        content: msg.content,
        attachments: msg.attachments,
        stickers: msg.stickerItems,
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
