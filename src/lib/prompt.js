function startPrompt(display, callback) {
  comcord.state.inPrompt = true;
  comcord.state.promptText = display;
  comcord.state.promptInput = "";

  comcord.state.promptCallback = callback;

  process.stdout.write(display);
}

async function finalizePrompt() {
  comcord.state.inPrompt = false;
  comcord.state.promptText = null;

  const input = comcord.state.promptInput.trim();
  await comcord.state.promptCallback(input);

  comcord.state.promptInput = null;
  comcord.state.promptCallback = null;
}

module.exports = {
  startPrompt,
  finalizePrompt,
};
