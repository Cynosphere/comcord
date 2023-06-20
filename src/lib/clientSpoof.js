/*
 * This single file is **EXCLUDED** from the project license.
 *
 * (c) 2022 Cynthia Foxwell, all rights reserved.
 * Permission is hereby granted to redistribute this file ONLY with copies of comcord.
 * You may not reverse engineer, modify, copy, or redistribute this file for any other uses outside of comcord.
 */

const os = require("os");

async function fetchMainPage() {
  const res = await fetch("https://discord.com/channels/@me");
  return await res.text();
}

async function fetchAsset(assetPath) {
  return await fetch("https://discord.com/" + assetPath).then((res) =>
    res.text()
  );
}

const MATCH_SCRIPT = '<script src="(.+?)" integrity=".+?">';
const REGEX_SCRIPT = new RegExp(MATCH_SCRIPT);
const REGEX_SCRIPT_GLOBAL = new RegExp(MATCH_SCRIPT, "g");

async function extractScripts() {
  const mainPage = await fetchMainPage();

  return mainPage
    .match(REGEX_SCRIPT_GLOBAL)
    .map((script) => script.match(REGEX_SCRIPT)[1]);
}

const REGEX_BUILD_NUMBER = /Build Number: (\d+), Version Hash:/;
const REGEX_BUILD_NUMBER_SWC = /Build Number: "\).concat\("(\d+)"/;

async function getBuildNumber() {
  if (comcord.state.cachedBuildNumber) {
    return comcord.state.cachedBuildNumber;
  }

  const scripts = await extractScripts();
  const chunkWithBuildInfoAsset = scripts[3];
  const chunkWithBuildInfo = await fetchAsset(chunkWithBuildInfoAsset);

  const buildNumber =
    chunkWithBuildInfo.match(REGEX_BUILD_NUMBER_SWC)?.[1] ??
    chunkWithBuildInfo.match(REGEX_BUILD_NUMBER)?.[1];
  comcord.state.cachedBuildNumber = buildNumber;

  return buildNumber;
}

/*async function getClientVersion() {
  if (comcord.state.cachedClientVersion) {
    return comcord.state.cachedClientVersion;
  }

  const data = await fetch(
    "https://updates.discord.com/distributions/app/manifests/latest?channel=stable&platform=win&arch=x86"
  ).then((res) => res.json());
  const clientVersion = data.full.host_version.join(".");
  comcord.state.cachedClientVersion = clientVersion;

  return clientVersion;
}*/

async function getBrowserInfo() {
  let targetOS;
  switch (process.platform) {
    case "win32":
    default:
      targetOS = "windows";
      break;
    case "darwin":
      targetOS = "mac os";
      break;
    case "linux":
      targetOS = "linux";
      break;
  }

  const data = await fetch(
    `https://cdn.jsdelivr.net/gh/ray-lothian/UserAgent-Switcher/v2/firefox/data/popup/browsers/firefox-${encodeURIComponent(
      targetOS
    )}.json`
  ).then((res) => res.json());
  data.sort((a, b) => Number(b.browser.major) - Number(a.browser.major));
  const target = data[0];

  return {ua: target.ua, version: target.browser.version};
}

async function getSuperProperties() {
  const buildNumber = await getBuildNumber();
  // const clientVersion = await getClientVersion();
  const browserInfo = await getBrowserInfo();

  let _os;
  switch (process.platform) {
    case "win32":
    default:
      _os = "Windows";
      break;
    case "darwin":
      _os = "Mac OS X";
      break;
    case "linux":
      _os = "Linux";
      break;
  }

  const props = {
    browser: "Firefox",
    browser_user_agent: browserInfo.ua,
    browser_version: browserInfo.version,
    client_build_number: buildNumber,
    client_event_source: null,
    device: "",
    os: _os,
    os_version: os.release(),
    //os_arch: os.arch(),
    referrer: "",
    referrer_current: "",
    referring_domain: "",
    referring_domain_current: "",
    release_channel: "stable",
    system_locale: "en-US",
  };
  return props;
}

module.exports = {getSuperProperties};
