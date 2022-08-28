function addCommand(key, name, callback) {
  if (comcord.commands[key]) {
    console.error(
      `Registering duplicate key for "${key}": "${name}" wants to overwrite "${comcord.commands[key].name}"!`
    );
    return;
  }

  comcord.commands[key] = {
    name,
    callback,
  };
}

module.exports = {
  addCommand,
};
