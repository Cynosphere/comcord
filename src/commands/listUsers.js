const chalk = require("chalk");

const {addCommand} = require("../lib/command");

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
  if (!comcord.state.currentGuild) {
    console.log("<not in a guild>");
    return;
  }
  if (!comcord.state.currentChannel) {
    console.log("<not in a channel>");
    return;
  }

  const guild = comcord.client.guilds.get(comcord.state.currentGuild);
  const channel = guild.channels.get(comcord.state.currentChannel);

  console.log(
    `\n[you are in '${guild.name}' in '${channel.name}' among ${guild.memberCount}]\n`
  );

  const online = [...guild.members.values()].filter((m) => m.presence);
  online.sort((a, b) => a.name - b.name);

  let longest = 0;
  for (const member of online) {
    const name = member.user.tag;
    if (name.length + 3 > longest) longest = name.length + 3;
  }

  const columns = Math.ceil(process.stdout.columns / longest);

  let index = 0;
  for (const member of online) {
    const name = member.user.tag;
    const status = getStatus(member.presence.status);
    const nameAndStatus = chalk.reset(name) + status;

    index++;
    process.stdout.write(
      nameAndStatus +
        " ".repeat(
          index % columns == 0 ? 0 : Math.abs(longest - (name.length + 3))
        )
    );

    if (index % columns == 0) process.stdout.write("\n");
  }
  if (index % columns != 0) process.stdout.write("\n");
  console.log("");

  if (channel.topic != null) {
    console.log("--Topic".padEnd(80, "-"));
    console.log(channel.topic);
    console.log("-".repeat(80));
    console.log("");
  }
}

if (!comcord.commands.w) {
  addCommand("w", "who is in guild", listUsers);
}

module.exports = {
  listUsers,
};
