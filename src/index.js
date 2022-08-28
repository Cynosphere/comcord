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
      chalk.bold.cyan(
        `[${name}]` + " ".repeat(nameLength - (name.length + 2))
      ) + chalk.reset(" " + content)
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
function setupSendMode() {
  inSendMode = true;
  toSend = "";
  const name = `[${client.user.username}]`;
  stdout.write(
    chalk.bold.cyan(name) +
      " ".repeat(nameLength - name.length) +
      chalk.reset(" ")
  );
}
function sendMessage() {
  toSend = toSend.trim();
  if (toSend === "") {
    stdout.write("<no message sent>\n");
  } else {
    client.createMessage(currentChannel, toSend);
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
    const length = `  ${key} - ${desc}`.length;

    stdout.write(
      "  " +
        chalk.bold.yellow(key) +
        chalk.reset(" - " + desc) +
        " ".repeat(Math.abs(25 - length))
    );

    index++;
    if (index % 3 == 0) stdout.write("\n");
  }
  if (index % 3 != 0) stdout.write("\n");

  console.log("\nTo begin TALK MODE, press [SPACE]\n");
}

stdin.on("data", function (key) {
  if (inSendMode) {
    if (key === "\r") {
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
