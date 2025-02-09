/**
 * IMGDD - A simple image hosting program
 * Copyright (C) 2025 @ericls
 *
 * Licensed under the GNU Affero General Public License v3.0.
 * See https://www.gnu.org/licenses/agpl-3.0.txt for details.
 */
import { IMGDDPlugin } from "./plugin";

window.IMGDD_PLUGINS = [];

function registerPlugin(plugin: IMGDDPlugin) {
  (window.IMGDD_PLUGINS ??= []).push(plugin);
}

window.registerPlugin = registerPlugin;
