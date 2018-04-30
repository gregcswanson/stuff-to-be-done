console.log('copy-pwa-to-site.js running');

const fs = require('fs-extra');

const source = './apps/polymer2app/build/es5-bundled';
const dir = './web/public/app';
const sourceApp = './apps/polymer2app/build/es5-bundled/index.html';
const appTemplate = './web/templates/layouts/app.html';

fs.ensureDirSync(dir);
fs.emptyDirSync(dir);
fs.copySync(source, dir);
fs.copySync(sourceApp, appTemplate);

console.log('copy-pwa-to-site.js complete');