const chalk = require("chalk");

const {addCommand} = require("../lib/command");

addCommand("h", "command help", function () {
  console.log("\nCOMcord (c)left 2022\n");

  const keys = Object.keys(comcord.commands);
  keys.sort((a, b) => a.localeCompare(b));

  let index = 0;
  for (const key of keys) {
    const desc = comcord.commands[key].name;
    const length = `  ${key} - ${desc}`.length;

    process.stdout.write(
      "  " +
        chalk.bold.yellow(key) +
        chalk.reset(" - " + desc) +
        " ".repeat(Math.abs(25 - length))
    );

    index++;
    if (index % 3 == 0) process.stdout.write("\n");
  }
  if (index % 3 != 0) process.stdout.write("\n");

  console.log("\nTo begin TALK MODE, press [SPACE]\n");
});
