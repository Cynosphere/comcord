const {addCommand} = require("../lib/command");

addCommand("c", "clear", function () {
  console.clear();
});
