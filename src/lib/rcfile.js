const fs = require("fs");
const path = require("path");
const os = require("os");

const RCPATH = path.resolve(os.homedir(), ".comcordrc");

function readFile(config) {
  const rc = fs.readFileSync(RCPATH, "utf8");
  const lines = rc.split("\n");
  for (const line of lines) {
    const [key, value] = line.split("=");
    config[key] = value;
  }
}

function writeFile(config) {
  if (fs.existsSync(RCPATH)) {
    readFile(config);
  }
  const newrc = [];

  for (const key in config) {
    const value = config[key];
    newrc.push(`${key}=${value}`);
  }

  fs.writeFileSync(RCPATH, newrc.join("\n"));
}

module.exports = {readFile, writeFile, path: RCPATH};
