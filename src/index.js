const Eris = require("eris");
const token = process.argv[2];

let currentGuild,
  currentChannel,
  inSendMode = false;

const messageQueue = [];

const client = new Eris("Bot " + token, {
  defaultImageFormat: "png",
  defaultImageSize: 1024,
  intents: Eris.Constants.Intents.all,
});

client.once("ready", function () {
  console.log(
    `Logged in as: ${client.user.username}#${client.user.discriminator} (${client.user.id})`
  );
});

client.on("messageCreate", function (msg) {
  if (msg.channel.id == currentChannel) {
    if (inSendMode) {
      messageQueue.push(msg);
    } else {
    }
  }
});

client.connect();
