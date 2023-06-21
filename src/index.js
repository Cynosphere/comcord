const {Client, Constants, Channel} = require("@projectdysnomia/dysnomia");
const DiscordRPC = require("discord-rpc");
const chalk = require("chalk");
const fs = require("fs");
const os = require("os");

const rcfile = require("./lib/rcfile");
const config = {};

if (fs.existsSync(rcfile.path)) {
  console.log(`% Reading ${rcfile.path.replace(os.homedir(), "~")} ...`);
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
    connected: true,
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
const client = new Client(
  (config.allowUserAccounts == "true" ? "" : "Bot ") + (token ?? config.token),
  {
    defaultImageFormat: "png",
    defaultImageSize: 1024,
    gateway: {
      intents: Object.values(Constants.Intents),
    },
    allowedMentions: {
      everyone: false,
    },
  }
);
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
    "Logged in as: " +
      chalk.yellow(`${client.user?.username} (${client.user?.id})`)
  );
  comcord.state.nameLength = (client.user?.username?.length ?? 0) + 2;

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

  if (client.user.bot && !config.disableRPC) {
    rpc
      .login({
        clientId: CLIENT_ID,
      })
      .catch(function () {});
  }
});
client.on("error", function () {});
client.on("ready", function () {
  if (comcord.state.connected === false) {
    console.log("% Reconnected");
  }
});
client.on("disconnect", function () {
  if (!comcord.state.quitting) {
    comcord.state.connected = false;
    console.log("% Disconnected, retrying...");
  }
});

rpc.on("connected", function () {
  comcord.state.rpcConnected = true;
  updatePresence();
});
let retryingRPC = false;
rpc.once("ready", function () {
  rpc.transport.on("error", function () {});
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
  if (
    (msg.mentions.find((user) => user.id == client.user.id) ||
      msg.mentionEveryone) &&
    msg.channel.id != comcord.state.currentChannel &&
    msg.channel.type !== Constants.ChannelTypes.DM &&
    msg.channel.type !== Constants.ChannelTypes.GROUP_DM
  ) {
    const data = {ping: true, channel: msg.channel, author: msg.author};
    if (comcord.state.inPrompt) {
      comcord.state.messageQueue.push(data);
    } else {
      processMessage(data);
    }
  }

  if (!msg.author) return;
  if (msg.author.id === client.user.id) return;

  if (
    !(msg.channel instanceof Channel) &&
    msg.author.id != client.user.id &&
    !msg.guildID
  ) {
    if (msg.channel.type === Constants.ChannelTypes.DM) {
      const newChannel = await client.getDMChannel(msg.author.id);
      if (msg.channel.id == newChannel.id) msg.channel = newChannel;
    } else if (msg.channel.type === Constants.ChannelTypes.GROUP_DM) {
      // TODO
    }
  }

  if (!(msg.channel instanceof Channel)) return;

  if (
    msg.channel.id == comcord.state.currentChannel ||
    msg.channel.type === Constants.ChannelTypes.DM ||
    msg.channel.type === Constants.ChannelTypes.GROUP_DM
  ) {
    if (comcord.state.inPrompt) {
      comcord.state.messageQueue.push(msg);
    } else {
      processMessage(msg);
    }
  }

  if (
    msg.channel.type === Constants.ChannelTypes.DM ||
    msg.channel.type === Constants.ChannelTypes.GROUP_DM
  ) {
    comcord.state.lastDM = msg.channel;
  }
});
client.on("messageUpdate", async function (msg, old) {
  if (!msg.author) return;
  if (msg.author.id === client.user.id) return;

  if (
    !(msg.channel instanceof Channel) &&
    msg.author.id != client.user.id &&
    !msg.guildID
  ) {
    if (msg.channel.type === Constants.ChannelTypes.DM) {
      const newChannel = await client.getDMChannel(msg.author.id);
      if (msg.channel.id == newChannel.id) msg.channel = newChannel;
    } else if (msg.channel.type === Constants.ChannelTypes.GROUP_DM) {
      // TODO
    }
  }

  if (!(msg.channel instanceof Channel)) return;

  if (
    msg.channel.id == comcord.state.currentChannel ||
    msg.channel.type === Constants.ChannelTypes.DM ||
    msg.channel.type === Constants.ChannelTypes.GROUP_DM
  ) {
    if (old && msg.content == old.content) return;

    if (comcord.state.inPrompt) {
      comcord.state.messageQueue.push(msg);
    } else {
      processMessage(msg);
    }
  }

  if (
    msg.channel.type === Constants.ChannelTypes.DM ||
    msg.channel.type === Constants.ChannelTypes.GROUP_DM
  ) {
    comcord.state.lastDM = msg.channel;
  }
});
client.on("messageReactionAdd", async function (msg, emoji, reactor) {
  if (msg.channel.id != comcord.state.currentChannel) return;
  const reply =
    msg.channel.messages.get(msg.id) ??
    (await msg.channel
      .getMessages({
        limit: 1,
        around: msg.id,
      })
      .then((msgs) => msgs[0]));

  const data = {
    channel: msg.channel,
    referencedMessage: reply,
    author: reactor?.user ?? client.users.get(reactor.id),
    timestamp: Date.now(),
    mentions: [],
    content: `*reacted with ${emoji.id ? `:${emoji.name}:` : emoji.name}*`,
  };

  if (comcord.state.inPrompt) {
    comcord.state.messageQueue.push(data);
  } else {
    processMessage(data);
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
    comcord.clientSpoof.superProperties = superProperties;
    comcord.clientSpoof.superPropertiesBase64 = Buffer.from(
      JSON.stringify(superProperties)
    ).toString("base64");

    // FIXME: is there a way we can string patch functions without having to
    //        dump locals into global
    global.MultipartData = require("@projectdysnomia/dysnomia/lib/util/MultipartData.js");
    global.SequentialBucket = require("@projectdysnomia/dysnomia/lib/util/SequentialBucket.js");
    global.DiscordHTTPError = require("@projectdysnomia/dysnomia/lib/errors/DiscordHTTPError.js");
    global.DiscordRESTError = require("@projectdysnomia/dysnomia/lib/errors/DiscordRESTError.js");
    global.Zlib = require("node:zlib");
    global.HTTPS = require("node:https");
    global.HTTP = require("node:http");
    global.GatewayOPCodes = Constants.GatewayOPCodes;
    global.GATEWAY_VERSION = Constants.GATEWAY_VERSION;

    client.getGateway = async function getGateway() {
      return {url: "wss://gateway.discord.gg"};
    };

    console.log("% Injecting headers into request handler");
    client.requestHandler.userAgent = superProperties.browser_user_agent;
    const requestFunction = client.requestHandler.request.toString();
    const newRequest = requestFunction
      .replace(
        "this.userAgent,",
        'this.userAgent,\n"X-Super-Properties":comcord.clientSpoof.superPropertiesBase64,'
      )
      .replace("._token", '._token.replace("Bot ","")');
    if (requestFunction === newRequest)
      throw new Error("Failed to patch request");
    client.requestHandler.request = new Function(
      "method",
      "url",
      "auth",
      "body",
      "file",
      "_route",
      "short",
      `return (function ${newRequest}).apply(this,arguments)`
    ).bind(client.requestHandler);

    console.log("% Injecting shard spawning");
    client.shards._spawn = client.shards.spawn.bind(client.shards);
    client.shards.spawn = function (id) {
      const res = this._spawn.apply(this, [id]);
      const shard = this.get(id);
      if (shard) {
        const identifyFunction = shard.identify.toString();
        const newIdentify = identifyFunction
          .replace(
            /properties: {\n\s+.+?\n\s+.+?\n\s+.+?\n\s+}\n/,
            "properties: comcord.clientSpoof.superProperties\n"
          )
          .replace(/\s+intents: this.client.shards.options.intents,/, "");
        if (identifyFunction === newIdentify)
          throw new Error("Failed to patch identify");
        shard.identify = new Function(
          `(function ${newIdentify}).apply(this, arguments)`
        );
        shard._wsEvent = shard.wsEvent;
        shard.wsEvent = function (packet) {
          if (packet.t == "READY") {
            packet.d.application = {id: CLIENT_ID, flags: 565248};
          }

          const ret = this._wsEvent.apply(this, [packet]);

          if (packet.t == "READY") {
            for (const guild of packet.d.guilds) {
              this._wsEvent.apply(this, [
                {
                  t: "GUILD_CREATE",
                  d: guild,
                },
              ]);
            }
          }

          return ret;
        };
      }

      return res;
    };

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
    comcord.state.nameLength = (client.user?.username?.length ?? 0) + 2;
    sentTime = true;
  } else if (seconds > 2 && sentTime) {
    sentTime = false;
  }
}, 500);
