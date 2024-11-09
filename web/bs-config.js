module.exports = {
    proxy: "localhost:8080",
    files: [
        {
            match: ["public/css/styles.css"],
            fn: function (event, file) {
                this.reload("*.css");
            }
        },
        "templates/**/*.gohtml",
        "public/js/*.js"
    ],
    port: 3000,
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
        cwd: '.',
        usePolling: false
    }
}; 