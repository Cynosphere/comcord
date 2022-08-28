const {addCommand} = require("../lib/command");

function listGuilds() {
  let longest = 0;
  const guilds = [];

  for (const guild of comcord.client.guilds.values()) {
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

addCommand("L", "list guilds", listGuilds);

module.exports = {
  listGuilds,
};
