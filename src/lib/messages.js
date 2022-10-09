const chalk = require("chalk");

const REGEX_CODEBLOCK = /```(?:([a-z0-9_+\-.]+?)\n)?\n*([^\n][^]*?)\n*```/i;
const REGEX_CODEBLOCK_GLOBAL =
  /```(?:([a-z0-9_+\-.]+?)\n)?\n*([^\n][^]*?)\n*```/gi;

const REGEX_MENTION = /<@!?(\d+)>/g;
const REGEX_ROLE_MENTION = /<@&?(\d+)>/g;
const REGEX_CHANNEL = /<#(\d+)>/g;
const REGEX_EMOTE = /<(?:\u200b|&)?a?:(\w+):(\d+)>/g;
const REGEX_COMMAND = /<\/([^\s]+?):(\d+)>/g;

function readableTime(time) {
  const seconds = time / 1000;
  const minutes = seconds / 60;
  const hours = minutes / 60;
  const days = hours / 24;
  const weeks = days / 7;
  const months = days / 30;
  const years = days / 365.25;

  if (years >= 1) {
    return `${years.toFixed(0)} year${years > 1 ? "s" : ""}`;
  } else if (weeks > 5 && months < 13) {
    return `${months.toFixed(0)} month${months > 1 ? "s" : ""}`;
  } else if (days > 7 && weeks < 5) {
    return `${weeks.toFixed(0)} week${weeks > 1 ? "s" : ""}`;
  } else if (hours > 24 && days < 7) {
    return `${days.toFixed(0)} day${days > 1 ? "s" : ""}`;
  } else if (minutes > 60 && hours < 24) {
    return `${hours.toFixed(0)} hour${hours > 1 ? "s" : ""}`;
  } else if (seconds > 60 && minutes < 60) {
    return `${minutes.toFixed(0)} minute${minutes > 1 ? "s" : ""}`;
  } else {
    return `${seconds.toFixed(0)} second${seconds > 1 ? "s" : ""}`;
  }
}

const MONTH_NAMES = [
  "January",
  "Feburary",
  "March",
  "April",
  "May",
  "June",
  "July",
  "August",
  "September",
  "October",
  "November",
  "December",
];
const DAY_NAMES = [
  "Sunday",
  "Monday",
  "Tuesday",
  "Wednesday",
  "Thursday",
  "Friday",
  "Saturday",
];
const TIME_FORMATS = {
  t: function (time) {
    const timeObj = new Date(time);
    return timeObj.getUTCHours() + 1 + ":" + timeObj.getUTCMinutes();
  },
  T: function (time) {
    const timeObj = new Date(time);
    return TIME_FORMATS.t(time) + ":" + timeObj.getUTCSeconds();
  },
  d: function (time) {
    const timeObj = new Date(time);
    return (
      timeObj.getUTCFullYear() +
      "/" +
      (timeObj.getUTCMonth() + 1).toString().padStart(2, "0") +
      "/" +
      timeObj.getUTCDate().toString().padStart(2, "0")
    );
  },
  D: function (time) {
    const timeObj = new Date(time);
    return (
      timeObj.getUTCDate() +
      " " +
      MONTH_NAMES[timeObj.getUTCMonth()] +
      " " +
      timeObj.getUTCFullYear()
    );
  },
  f: function (time) {
    return TIME_FORMATS.D(time) + " " + TIME_FORMATS.t(time);
  },
  F: function (time) {
    const timeObj = new Date(time);
    return DAY_NAMES[timeObj.getUTCDay()] + ", " + TIME_FORMATS.f(time);
  },
  R: function (time) {
    const now = Date.now();

    if (time > now) {
      const delta = time - now;
      return "in " + readableTime(delta);
    } else {
      const delta = now - time;
      return readableTime(delta) + " ago";
    }
  },
};
const REGEX_TIMESTAMP = new RegExp(
  `<t:(-?\\d{1,17})(?::(${Object.keys(TIME_FORMATS).join("|")}))?>`,
  "g"
);

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
function replaceTimestamps(_, time, format = "f") {
  return TIME_FORMATS[format](time * 1000);
}

function formatMessage({
  name,
  content,
  bot,
  attachments,
  stickers,
  reply,
  noColor = false,
  dump = false,
  history = false,
}) {
  if (name.length + 2 > comcord.state.nameLength)
    comcord.state.nameLength = name.length + 2;

  if (reply) {
    const nameColor = reply.author.bot ? chalk.bold.yellow : chalk.bold.cyan;

    const headerLength = 5 + reply.author.username.length;

    let replyContent = reply.content.replace(/\n/g, " ");
    replyContent = replyContent
      .replace(REGEX_MENTION, replaceMentions)
      .replace(REGEX_ROLE_MENTION, replaceRoles)
      .replace(REGEX_CHANNEL, replaceChannels)
      .replace(REGEX_EMOTE, replaceEmotes)
      .replace(REGEX_COMMAND, replaceCommands)
      .replace(REGEX_TIMESTAMP, replaceTimestamps);
    if (reply.attachments.length > 0) {
      replyContent += ` <${reply.attachments.length} attachment${
        reply.attachments.length > 1 ? "s" : ""
      }>`;
      replyContent = replyContent.trim();
    }

    const length = headerLength + replyContent.length;

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

  if (dump) {
    if (history) {
      const headerLength = 80 - (name.length + 5);
      console.log(`--- ${name} ${"-".repeat(headerLength)}`);
      console.log(content);
      console.log(`--- ${name} ${"-".repeat(headerLength)}`);
    } else {
      const wordCount = content.split(" ").length;
      const lineCount = content.split("\n").length;
      if (noColor) {
        console.log(
          `<${name} DUMPs in ${content.length} characters of ${wordCount} word${
            wordCount > 1 ? "s" : ""
          } in ${lineCount} line${lineCount > 1 ? "s" : ""}>`
        );
      } else {
        console.log(
          chalk.bold.yellow(
            `<${name} DUMPs in ${
              content.length
            } characters of ${wordCount} word${
              wordCount > 1 ? "s" : ""
            } in ${lineCount} line${lineCount > 1 ? "s" : ""}>`
          )
        );
      }
    }
  } else {
    content = content
      .replace(REGEX_MENTION, replaceMentions)
      .replace(REGEX_ROLE_MENTION, replaceRoles)
      .replace(REGEX_CHANNEL, replaceChannels)
      .replace(REGEX_EMOTE, replaceEmotes)
      .replace(REGEX_COMMAND, replaceCommands)
      .replace(REGEX_TIMESTAMP, replaceTimestamps);

    if (
      (content.length > 1 &&
        content.startsWith("*") &&
        content.endsWith("*")) ||
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
  }

  if (attachments) {
    for (const attachment of attachments.values()) {
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

function processMessage(msg, options) {
  if (msg.time) {
    console.log(msg.content);
  } else if (msg.content.indexOf("\n") > -1) {
    if (msg.content.match(REGEX_CODEBLOCK)) {
      formatMessage({
        name: msg.author.username,
        bot: msg.author.bot,
        content: msg.content.replace(
          REGEX_CODEBLOCK_GLOBAL,
          (_, lang, content) => content
        ),
        attachments: msg.attachments,
        stickers: msg.stickerItems,
        reply: msg.referencedMessage,
        dump: true,
        ...options,
      });
    } else {
      const lines = msg.content.split("\n");
      for (const index in lines) {
        const line = lines[index];
        formatMessage({
          name: msg.author.username,
          bot: msg.author.bot,
          content:
            line +
            (msg.editedTimestamp != null && index == lines.length - 1
              ? " (edited)"
              : ""),
          attachments: index == lines.length - 1 ? msg.attachments : [],
          stickers: index == lines.length - 1 ? msg.stickerItems : [],
          reply: index == 0 ? msg.referencedMessage : null,
          ...options,
        });
      }
    }
  } else {
    formatMessage({
      name: msg.author.username,
      bot: msg.author.bot,
      content: msg.content + (msg.editedTimestamp != null ? " (edited)" : ""),
      attachments: msg.attachments,
      stickers: msg.stickerItems,
      reply: msg.referencedMessage,
      ...options,
    });
  }
}

function processQueue() {
  for (const msg of comcord.state.messageQueue) {
    processMessage(msg);
  }

  comcord.state.messageQueue.splice(0, comcord.state.messageQueue.length);
}

module.exports = {
  processMessage,
  processQueue,
  formatMessage,
};
