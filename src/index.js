const Eris = require("eris");
const chalk = require("chalk");

const token = process.argv[2];
const stdin = process.stdin;
const stdout = process.stdout;

stdin.setRawMode(true);
stdin.resume();
stdin.setEncoding("utf8");

let currentGuild,
  currentChannel,
  inSendMode = false,
  guildSwitch = false,
  channelSwitch = false,
  nameLength = 2;

const messageQueue = [];

const commands = {
  q: "quit comcord",
  e: "emote",
  g: "goto a channel",
  G: "goto a guild",
  l: "list channels",
  L: "list guilds",
  w: "who is in guild",
  f: "finger",
  r: "channel history",
  R: "extended channel history",
  h: "command help",
  c: "clear",
};

const client = new Eris("Bot " + token, {
  defaultImageFormat: "png",
  defaultImageSize: 1024,
  intents: Eris.Constants.Intents.all,
});

client.once("ready", function () {
  console.log(
    "Logged in as: " +
      chalk.yellow(
        `${client.user.username}#${client.user.discriminator} (${client.user.id})`
      )
  );
  nameLength = client.user.username.length + 2;

  listGuilds();
});

function processMessage({name, content, bot}) {
  if (name.length + 2 > nameLength) nameLength = name.length + 2;

  if (
    (content.startsWith("*") && content.endsWith("*")) ||
    (content.startsWith("_") && content.endsWith("_"))
  ) {
    console.log(chalk.bold.green(`<${name} ${content}>`));
  } else {
    // TODO: markdown
    console.log(
      chalk.bold.cyan(`[${name}]`).padEnd(nameLength, " ") +
        chalk.reset(" " + content)
    );
  }
}

function processQueue() {
  for (const msg of messageQueue) {
    if (msg.content.indexOf("\n") > -1) {
      const lines = msg.content.split("\n");
      for (const line of lines) {
        processMessage({
          name: msg.author.name,
          bot: msg.author.bot,
          content: line,
        });
      }
    } else {
      processMessage({
        name: msg.author.name,
        bot: msg.author.bot,
        content: msg.content,
      });
    }
  }
}

client.on("messageCreate", function (msg) {
  if (msg.channel.id == currentChannel) {
    if (inSendMode) {
      messageQueue.push(msg);
    } else {
      if (msg.content.indexOf("\n") > -1) {
        const lines = msg.content.split("\n");
        for (const line of lines) {
          processMessage({
            name: msg.author.name,
            bot: msg.author.bot,
            content: line,
          });
        }
      } else {
        processMessage({
          name: msg.author.name,
          bot: msg.author.bot,
          content: msg.content,
        });
      }
    }
  }
});

let toSend = "";
async function setupSendMode() {
  inSendMode = true;
  toSend = "";
  stdout.write(
    chalk.bold.cyan(`[${client.user.username}]`).padEnd(nameLength, " ") +
      chalk.reset(" ")
  );
  try {
    await client.guilds
      .get(currentGuild)
      .channels.get(currentChannel)
      .sendTyping();
  } catch (err) {
    //
  }
}
async function sendMessage() {
  toSend = toSend.trim();
  if (toSend === "") {
    stdout.write("<no message sent>\n");
  } else {
    try {
      await client.createMessage(currentChannel, toSend);
    } catch (err) {
      console.log("<failed to send message: " + err.message + ">");
    }
  }
  inSendMode = false;
  processQueue();
}

function showHelp() {
  console.log("\nCOMcord (c)left 2022\n");

  const keys = Object.keys(commands);
  keys.sort();

  let index = 0;
  for (const key of keys) {
    const desc = commands[key];

    stdout.write(
      ("  " + chalk.bold.yellow(key) + chalk.reset(" - " + desc)).padEnd(
        25,
        " "
      )
    );

    index++;
    if (index % 3 == 0) stdout.write("\n");
  }
  if (index % 3 != 0) stdout.write("\n");

  console.log("\nTo begin TALK MODE, press [SPACE]\n");
}

function listGuilds() {
  let longest = 0;
  const guilds = [];

  for (const guild of client.guilds.values()) {
    if (guild.name.length > longest) longest = guild.name.length;

    const online = [...guild.members.values()].filter((m) => m.status).length;
    guilds.push({name: guild.name, members: guild.memberCount, online});
  }

  console.log("");
  console.log("  " + "guild-name".padStart(longest, " ") + "  online  total");
  console.log("-".repeat(80));
  for (const guild of guilds) {
    console.log(
      "  " +
        guild.name.padStart(longest, " ") +
        "  " +
        guild.online.toString().padStart(6, " ") +
        "  " +
        guild.members.toString().padStart(5, " ")
    );
  }
  console.log("");
}

let targetGuild = "";
function gotoGuild() {
  targetGuild = "";
  guildSwitch = true;

  stdout.write(":guild> ");
}

function findTopChannel(guildId) {
  const guild = client.guilds.get(guildId);
  const channels = [...guild.channels.values()].filter((c) => c.type == 0);
  channels.sort((a, b) => a.position - b.position);

  return channels[0];
}

function getStatus(status) {
  let color;
  switch (status) {
    case "online":
      color = chalk.bold.green;
      break;
    case "idle":
      color = chalk.bold.yellow;
      break;
    case "dnd":
      color = chalk.bold.red;
      break;
    default:
      color = chalk.bold;
      break;
  }

  return color(" \u2022 ");
}

function listUsers() {
  const guild = client.guilds.get(currentGuild);
  const channel = guild.channels.get(currentChannel);

  console.log(
    `\n[you are in '${guild.name}' in '${channel.name}' among ${guild.memberCount}]\n`
  );

  const online = [...guild.members.values()].filter((m) => m.status);
  online.sort((a, b) => a.name - b.name);

  let longest = 0;
  for (const member of online) {
    const name = member.user.username + "#" + member.user.discriminator;
    if (name.length + 3 > longest) longest = name.length + 3;
  }

  const columns = Math.ceil(stdout.columns / longest);

  let index = 0;
  for (const member of online) {
    const name = member.user.username + "#" + member.user.discriminator;
    const status = getStatus(member.status);
    const nameAndStatus = chalk.reset(name) + status;

    index++;
    stdout.write(
      nameAndStatus +
        " ".repeat(
          index % columns == 0 ? 0 : Math.abs(longest - (name.length + 3))
        )
    );

    if (index % columns == 0) stdout.write("\n");
  }
  if (index % columns != 0) stdout.write("\n");
  console.log("");

  if (channel.topic != null) {
    console.log("--Topic".padEnd(80, "-"));
    console.log(channel.topic);
    console.log("-".repeat(80));
    console.log("");
  }
}

function switchGuild() {
  let target;

  for (const guild of client.guilds.values()) {
    if (guild.name.toLowerCase().indexOf(targetGuild.toLowerCase()) > -1) {
      target = guild.id;
      break;
    }
  }

  if (target == null) {
    console.log("<guild not found>");
  } else {
    currentGuild = target;
    // TODO: store last visited channel
    const topChannel = findTopChannel(target);
    currentChannel = topChannel.id;

    listUsers();
  }

  guildSwitch = false;
}

stdin.on("data", function (key) {
  if (guildSwitch) {
    if (key === "\r") {
      console.log("");
      switchGuild();
    } else {
      if (key === "\b") {
        if (targetGuild.length > 0) {
          stdout.moveCursor(-1);
          stdout.write(" ");
          stdout.moveCursor(-1);
          targetGuild = targetGuild.substring(0, targetGuild.length - 1);
        }
      } else {
        stdout.write(key);
        targetGuild += key;
      }
    }
  } else if (inSendMode) {
    if (key === "\r") {
      console.log("");
      sendMessage();
    } else {
      if (key === "\b") {
        if (toSend.length > 0) {
          stdout.moveCursor(-1);
          stdout.write(" ");
          stdout.moveCursor(-1);
          toSend = toSend.substring(0, toSend.length - 1);
        }
      } else {
        stdout.write(key);
        toSend += key;
      }
    }
  } else {
    switch (key) {
      case "\u0003":
      case "q": {
        client.disconnect(false);
        process.exit(0);
        break;
      }
      case "h": {
        showHelp();
        break;
      }
      case "g": {
        if (currentGuild == null) {
          console.log("<not in a guild>");
          break;
        }
        break;
      }
      case "G": {
        gotoGuild();
        break;
      }
      case "l": {
        if (currentGuild == null) {
          console.log("<not in a guild>");
          break;
        }
        break;
      }
      case "L": {
        listGuilds();
        break;
      }
      case "w": {
        if (currentGuild == null) {
          console.log("<not in a guild>");
          break;
        }
        listUsers();
        break;
      }
      case " ":
      case "\r":
      default: {
        if (currentChannel == null) {
          console.log("<not in a channel>");
          break;
        }
        setupSendMode();
        break;
      }
    }
  }
});

client.connect();

console.log("COMcord (c)left 2022");
console.log("Type 'h' for Commands");
