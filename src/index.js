const {Client, Constants} = require("oceanic.js");
const DiscordRPC = require("discord-rpc");
const chalk = require("chalk");
const fs = require("fs");

const rcfile = require("./lib/rcfile");
const config = {};

if (fs.existsSync(rcfile.path)) {
  console.log("% Reading " + rcfile.path + " ...");
  rcfile.readFile(config);
}

const CLIENT_ID = "1026163285877325874";

const token = process.argv[2];
if (!config.token && token) {
  console.log("% Writing token to .comcordrc");
  config.token = token;
  rcfile.writeFile(config);
}

if (!config.token && !token) {
  console.log("No token provided.");
  process.exit(1);
}

process.title = "comcord";

global.comcord = {
  config,
  state: {
    rpcConnected: false,
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
const client = new Client({
  auth:
    (config.allowUserAccounts == "true" ? "" : "Bot ") +
    (token ?? config.token),
  defaultImageFormat: "png",
  defaultImageSize: 1024,
  gateway: {
    intents: ["ALL"],
    maxShards: 1,
    concurrency: 1,
    presence: {
      status: "online",
      activities: [
        {
          name: "comcord",
          type: 0,
          application_id: CLIENT_ID,
          timestamps: {
            start: comcord.state.startTime,
          },
        },
      ],
    },
  },
});
comcord.client = client;
const rpc = new DiscordRPC.Client({transport: "ipc"});
comcord.rpc = rpc;

const {finalizePrompt} = require("./lib/prompt");
const {processMessage, processQueue} = require("./lib/messages");
const {updatePresence} = require("./lib/presence");

require("./commands/quit");
require("./commands/clear");
require("./commands/help");
const {sendMode} = require("./commands/send");
require("./commands/emote");
const {listGuilds} = require("./commands/listGuilds");
const {switchGuild} = require("./commands/switchGuild"); // loads listChannels and listUsers
require("./commands/switchChannel"); //loads listUsers
require("./commands/history"); // includes extended history
require("./commands/afk");
require("./commands/privateMessages");

process.stdin.setRawMode(true);
process.stdin.resume();
process.stdin.setEncoding("utf8");

client.once("ready", function () {
  console.log(
    "Logged in as: " + chalk.yellow(`${client.user.tag} (${client.user.id})`)
  );
  comcord.state.nameLength = client.user.username.length + 2;

  listGuilds();

  if (config.defaultGuild) {
    const guild = client.guilds.get(config.defaultGuild);
    if (guild != null) {
      if (config.defaultChannel) {
        comcord.state.currentChannel = config.defaultChannel;
        comcord.state.lastChannel.set(
          config.defaultGuild,
          config.defaultChannel
        );
      }
      switchGuild(guild.name);
    } else {
      console.log("% This account is not in the defined default guild.");
    }
  } else {
    if (config.defaultChannel) {
      console.log("% Default channel defined without defining default guild.");
    }
  }

  if (client.user.bot) {
    rpc
      .login({
        clientId: CLIENT_ID,
      })
      .catch(function () {});
  }
});
client.on("error", function () {});

rpc.on("connected", function () {
  comcord.state.rpcConnected = true;
  updatePresence();
});
let retryingRPC = false;
rpc.once("ready", function () {
  rpc.transport.on("close", function () {
    comcord.state.rpcConnected = false;
    if (!retryingRPC) {
      retryingRPC = true;
      setTimeout(function () {
        rpc.transport
          .connect()
          .then(() => {
            retryingRPC = false;
          })
          .catch((err) => {
            retryingRPC = false;
            rpc.transport.emit("close");
          });
      }, 5000);
    }
  });
});
rpc.on("error", function () {});

client.on("messageCreate", async function (msg) {
  if (msg.author.id === client.user.id) return;

  if (msg.channelID && !msg.channel) {
    try {
      const dmChannel = await msg.author.createDM();
      if (dmChannel.id === msg.channelID) {
        msg.channel = dmChannel;
      }
    } catch {
      //
    }
  }

  if (
    (msg.channel ? msg.channel.id : msg.channelID) ==
      comcord.state.currentChannel ||
    msg.channel?.recipient != null
  ) {
    if (comcord.state.inPrompt) {
      comcord.state.messageQueue.push(msg);
    } else {
      processMessage(msg);
    }
  }

  if (msg.channel?.recipient != null) {
    comcord.state.lastDM = msg.author;
  }
});
client.on("messageUpdate", async function (msg, old) {
  if (msg.author.id === client.user.id) return;

  if (msg.channelID && !msg.channel) {
    try {
      const dmChannel = await msg.author.createDM();
      if (dmChannel.id === msg.channelID) {
        msg.channel = dmChannel;
      }
    } catch {
      //
    }
  }

  if (
    (msg.channel ? msg.channel.id : msg.channelID) ==
      comcord.state.currentChannel ||
    msg.channel?.recipient != null
  ) {
    if (old && msg.content == old.content) return;

    if (comcord.state.inPrompt) {
      comcord.state.messageQueue.push(msg);
    } else {
      processMessage(msg);
    }
  }

  if (msg.channel?.recipient != null) {
    comcord.state.lastDM = msg.author;
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

if (
  config.allowUserAccounts == "true" &&
  !(token ?? config.token).startsWith("Bot ")
) {
  if (fetch == null) {
    console.log("Node v18+ needed for user account support.");
    process.exit(1);
  }

  (async function () {
    comcord.clientSpoof = require("./lib/clientSpoof");
    const superProperties = await comcord.clientSpoof.getSuperProperties();

    console.log("% Allowing non-bot tokens to connect");
    const connectLines = client.connect.toString().split("\n");
    connectLines.splice(0, 4);
    connectLines.splice(-1, 1);

    const newConnect = new client.connect.constructor(connectLines.join("\n"));
    client.connect = newConnect.bind(client);

    // gross hack
    global.Constants_1 = Constants;
    try {
      global.Erlpack = require("erlpack");
    } catch {
      global.Erlpack = false;
    }

    console.log("% Injecting headers into request handler");
    client.rest.handler.options.userAgent = `Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) discord/${superProperties.client_version} Chrome/91.0.4472.164 Electron/13.6.6 Safari/537.36`;
    client.rest.handler._request = client.rest.handler.request.bind(
      client.rest.handler
    );
    client.rest.handler.request = async function (options) {
      options.headers = options.headers ?? {};
      options.headers["X-Super-Properties"] =
        await comcord.clientSpoof.getSuperPropertiesBase64();

      return await this._request.apply(this, [options]);
    }.bind(client.rest.handler);

    console.log("% Setting gateway connection properties");
    client.shards.options.connectionProperties = superProperties;

    console.log("% Injecting application into READY payload");
    client.shards._spawn = client.shards.spawn.bind(client.shards);
    client.shards.spawn = function (id) {
      const res = this._spawn.apply(this, [id]);
      const shard = this.get(id);
      if (shard) {
        shard._onDispatch = shard.onDispatch.bind(shard);
        shard.onDispatch = async function (packet) {
          if (packet.t == "READY") {
            packet.d.application = {id: CLIENT_ID, flags: 565248};
          }

          const ret = await this._onDispatch.apply(this, [packet]);

          if (packet.t == "READY") {
            for (const guild of packet.d.guilds) {
              await this._onDispatch.apply(this, [
                {
                  t: "GUILD_CREATE",
                  d: guild,
                },
              ]);
            }
          }

          return ret;
        }.bind(shard);
      }

      return res;
    }.bind(client.shards);

    console.log("% Connecting to gateway now");
    await client.connect();
  })();
} else {
  client.connect();
}

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
    comcord.state.nameLength = client.user.username.length + 2;
    sentTime = true;
  } else if (seconds > 2 && sentTime) {
    sentTime = false;
  }
}, 500);
