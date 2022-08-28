const Eris = require("eris");

const {addCommand} = require("../lib/command");

function listChannels() {
  if (!comcord.state.currentGuild) {
    console.log("<not in a guild>");
    return;
  }

  let longest = 0;
  let longestTopic = 0;
  const guild = comcord.client.guilds.get(comcord.state.currentGuild);
  const channels = [...guild.channels.values()].filter((c) => c.type == 0);
  channels.sort((a, b) => a.position - b.position);

  for (const channel of channels) {
    const perms = channel.permissionsOf(comcord.client.user.id);
    const private = !perms.has(Eris.Constants.Permissions.readMessages);

    if (channel.name.length + (private ? 1 : 0) > longest)
      longest = Math.min(25, channel.name.length + (private ? 1 : 0));
    if (channel.topic != null && channel.topic.length > longestTopic)
      longestTopic = channel.topic.length;
  }

  console.log("");
  console.log(
    "  " +
      "channel-name".padStart(longest, " ") +
      "  " +
      "topic".padStart(Math.min(80 - (longest + 5), longestTopic), " ")
  );
  console.log("-".repeat(80));
  for (const channel of channels) {
    const topic =
      channel.topic != null ? channel.topic.replace(/\n/g, " ") : "";
    const perms = channel.permissionsOf(comcord.client.user.id);
    const private = !perms.has(Eris.Constants.Permissions.readMessages);

    const name = (private ? "*" : "") + channel.name;

    console.log(
      "  " +
        (name.length > 24 ? name.substring(0, 24) + "\u2026" : name).padStart(
          longest,
          " "
        ) +
        "  " +
        (topic.length > 80 - longest + 9
          ? topic.substring(0, 79 - (longest + 5)) + "\u2026"
          : topic.padStart(Math.min(80 - (longest + 5), longestTopic), " "))
    );
  }
  console.log("");
}

addCommand("l", "list channels", listChannels);

module.exports = {
  listChannels,
};
