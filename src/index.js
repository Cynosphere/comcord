const Eris = require("eris");
const chalk = require("chalk");

const token = process.argv[2];

process.title = "comcord";

global.comcord = {
  state: {
    startTime: Date.now(),
    currentGuild: null,
    currentChannel: null,
    nameLength: 2,
    inPrompt: false,
    messageQueue: [],
    lastChannel: new Map(),
    afk: false,
  },
  commands: {},
};
const client = new Eris("Bot " + token, {
  defaultImageFormat: "png",
  defaultImageSize: 1024,
  intents: Eris.Constants.Intents.all,
});
comcord.client = client;

const {finalizePrompt} = require("./lib/prompt");
const {processMessage, processQueue} = require("./lib/messages");

require("./commands/quit");
require("./commands/clear");
require("./commands/help");
const {sendMode} = require("./commands/send");
require("./commands/emote");
const {listGuilds} = require("./commands/listGuilds");
require("./commands/switchGuild"); // loads listChannels and listUsers
require("./commands/switchChannel"); //loads listUsers
require("./commands/history"); // includes extended history
require("./commands/afk");

process.stdin.setRawMode(true);
process.stdin.resume();
process.stdin.setEncoding("utf8");

client.once("ready", function () {
  console.log(
    "Logged in as: " +
      chalk.yellow(
        `${client.user.username}#${client.user.discriminator} (${client.user.id})`
      )
  );
  comcord.state.nameLength = client.user.username.length + 2;

  listGuilds();

  client.editStatus("online", [
    {
      application_id: "1026163285877325874",
      name: "comcord",
      timestamps: {
        start: comcord.state.startTime,
      },
    },
  ]);
});
client.on("error", function () {});

client.on("messageCreate", function (msg) {
  if (msg.author.id === client.user.id) return;

  if (msg.channel.id == comcord.state.currentChannel) {
    if (comcord.state.inPrompt) {
      comcord.state.messageQueue.push(msg);
    } else {
      processMessage(msg);
    }
  }
});
client.on("messageUpdate", function (msg, old) {
  if (msg.author.id === client.user.id) return;

  if (msg.channel.id == comcord.state.currentChannel) {
    if (msg.content == old.content) return;

    if (comcord.state.inPrompt) {
      comcord.state.messageQueue.push(msg);
    } else {
      processMessage(msg);
    }
  }
});

process.stdin.on("data", async function (key) {
  if (comcord.state.inPrompt) {
    if (key === "\r") {
      await finalizePrompt();
      processQueue();
    } else {
      if (key === "\b" || key === "\u007f") {
        if (comcord.state.promptInput.length > 0) {
          process.stdout.moveCursor(-1);
          process.stdout.write(" ");
          process.stdout.moveCursor(-1);
          comcord.state.promptInput = comcord.state.promptInput.substring(
            0,
            comcord.state.promptInput.length - 1
          );
        }
      } else {
        key = key.replace("\u001b", "");
        process.stdout.write(key);
        comcord.state.promptInput += key;
      }
    }
  } else {
    if (comcord.commands[key]) {
      comcord.commands[key].callback();
    } else {
      sendMode();
    }
  }
});

client.connect();

console.log("COMcord (c)left 2022");
console.log("Type 'h' for Commands");

const dateObj = new Date();
let sentTime = false;

const weekdays = ["Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"];
const months = [
  "Jan",
  "Feb",
  "Mar",
  "Apr",
  "May",
  "Jun",
  "Jul",
  "Aug",
  "Sep",
  "Oct",
  "Nov",
  "Dec",
];

setInterval(function () {
  dateObj.setTime(Date.now());

  const hour = dateObj.getUTCHours(),
    minutes = dateObj.getUTCMinutes(),
    seconds = dateObj.getUTCSeconds(),
    day = dateObj.getUTCDate(),
    month = dateObj.getUTCMonth(),
    year = dateObj.getUTCFullYear(),
    weekDay = dateObj.getUTCDay();

  const timeString = `[${weekdays[weekDay]} ${day
    .toString()
    .padStart(2, "0")}-${months[month]}-${year
    .toString()
    .substring(2, 4)} ${hour.toString().padStart(2, "0")}:${minutes
    .toString()
    .padStart(2, "0")}:${seconds.toString().padStart(2, "0")}]`;

  if (minutes % 15 == 0 && seconds < 2 && !sentTime) {
    if (comcord.state.inPrompt == true) {
      comcord.state.messageQueue.push({time: true, content: timeString});
    } else {
      console.log(timeString);
    }
    sentTime = true;
  } else if (seconds > 2 && sentTime) {
    sentTime = false;
  }
}, 500);
