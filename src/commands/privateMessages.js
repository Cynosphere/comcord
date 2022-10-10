const chalk = require("chalk");

const {addCommand} = require("../lib/command");
const {startPrompt} = require("../lib/prompt");
const {listUsers} = require("./listUsers");

function startDM(user) {
  startPrompt(":msg> ", async function (input) {
    if (input == "") {
      console.log(`\n<message not sent to ${user.tag}>`);
    } else {
      try {
        const channel = await user.createDM();
        await channel.createMessage({content: input});
        console.log(chalk.bold.green(`\n<message sent to ${user.tag}>`));
      } catch (err) {
        console.log("\n<failed to send message: " + err.message + ">");
      }
    }
  });
}

addCommand("s", "send private", function () {
  console.log("Provide a RECIPIENT");
  startPrompt(":to> ", function (who) {
    let target;
    for (const user of comcord.client.users.values()) {
      if (user.tag == who) {
        target = user;
        break;
      }
    }

    if (target) {
      console.log("");
      startDM(target);
    } else {
      listUsers();
    }
  });
});

addCommand("a", "answer a send", function () {
  if (comcord.state.lastDM) {
    console.log(chalk.bold.green(`<answering ${comcord.state.lastDM.tag}>`));
    startDM(comcord.state.lastDM);
  } else {
    // FIXME: figure out the actual message in com
    console.log("<no one to answer>");
  }
});
