console.log("Hello from Node.js!");
console.log("Environment:", process.env.NODE_ENV);
console.log("Workers:", process.env.GONODE_WORKERS);
setInterval(() => {
    console.log("Aplikasi sedang berjalan...");
}, 5000);