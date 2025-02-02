import { IMGDDPlugin } from "./plugin";

window.IMGDD_PLUGINS = [];

function registerPlugin(plugin: IMGDDPlugin) {
  (window.IMGDD_PLUGINS ??= []).push(plugin);
}

window.registerPlugin = registerPlugin;
