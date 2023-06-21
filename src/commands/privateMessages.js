const chalk = require("chalk");

const {addCommand} = require("../lib/command");
const {startPrompt} = require("../lib/prompt");
const {listUsers} = require("./listUsers");

function startDM(channel) {
  startPrompt(":msg> ", async function (input) {
    if (input == "") {
      console.log(
        `\n<message not sent to ${channel.recipent?.username ?? "group DM"}>`
      );
    } else {
      try {
        await channel.createMessage({content: input});
        console.log(
          chalk.bold.green(
            `\n<message sent to ${channel.recipient?.username ?? "group DM"}>`
          )
        );
      } catch (err) {
        console.log(`\n<failed to send message: ${err.message}>`);
      }
    }
  });
}

addCommand("s", "send private", function () {
  console.log("Provide a RECIPIENT");
  startPrompt(":to> ", function (who) {
    let target;
    for (const user of comcord.client.users.values()) {
      if (user.username == who) {
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
    console.log(
      chalk.bold.green(
        `<answering ${comcord.state.lastDM.recipient?.username ?? "group DM"}>`
      )
    );
    startDM(comcord.state.lastDM);
  } else {
    // FIXME: figure out the actual message in com
    console.log("<no one to answer>");
  }
});
