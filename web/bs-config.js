module.exports = {
    proxy: "localhost:8080",
    files: [
        "web/templates/**/*.gohtml",
        "web/public/css/**/*.css",
        "web/public/js/**/*.js"
    ],
    port: 3000,
    open: false,
    notify: false,
    reloadDelay: 50,
    reloadDebounce: 250
}; 