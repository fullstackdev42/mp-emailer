require('dotenv').config({ path: '../.env' });

const serverPort = parseInt(process.env.APP_PORT || '8080', 10);
const serverHost = process.env.APP_HOST || 'localhost';
console.log(`Proxying to ${serverHost}:${serverPort}`);

module.exports = {
    proxy: `${serverHost}:${serverPort}`,
    files: [
        {
            match: ["public/css/styles.css"],
            fn: function (event, file) {
                this.reload("*.css");
            }
        },
        {
            match: ["templates/**/*.gohtml"],
            fn: function (event, file) {
                console.log(`Template changed: ${file}`);
                this.reload();
            }
        },
        "public/js/*.js"
    ],
    port: 3000,
    host: serverHost,
    open: false,
    notify: false,
    reloadDelay: 100,
    reloadDebounce: 250,
    logLevel: "debug",
    logPrefix: "BS",
    watchEvents: ["change", "add", "unlink", "addDir", "unlinkDir"],
    watchOptions: {
        ignored: ['node_modules', 'tmp'],
        ignoreInitial: true,
        cwd: 'web',
        usePolling: false,
        awaitWriteFinish: {
            stabilityThreshold: 100,
            pollInterval: 100
        }
    }
}; 