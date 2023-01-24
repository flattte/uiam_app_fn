const app = require("express")();
const cors = require("cors");
require("dotenv").HOST;
const https = require("https");
const fs = require("fs");
 
const server = https.createServer(
  {
    key: fs.readFileSync("../server/ssl/key.pem"),
    cert: fs.readFileSync("../server/ssl/cert.pem"),
  },
  app);

const io = require("socket.io")(server, {
  cors: {
    origin: "*",
    methods: ["GET", "POST"],
  },
});

app.use(cors());

const PORT = process.env.PORT || 5000;
const HOST = process.env.HOST;

app.get("/", (req, res) => {
  res.send("Server is running properly");
});

io.on("connection", (socket) => {
  socket.emit("me", socket.id);
  console.log("emit me ran");

  socket.on("disconnect", () => {
    socket.broadcast.emit("callended");
  });

  socket.on("calluser", ({ userToCall, signalData, from, name }) => {
    console.log("calluser ran");
    io.to(userToCall).emit("calluser", { signal: signalData, from, name });
    console.log("calluser emitted");
  });

  socket.on("answercall", (data) => {
    io.to(data.to).emit("callaccepted", { signal: data.signal });
    console.log("callaccepted emitted");
  });
});

server.listen(PORT, HOST,() => {
  console.log(`Server is listening on host ${HOST} port ${PORT}`);
});
