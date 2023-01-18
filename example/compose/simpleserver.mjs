#!/usr/bin/env node

import * as http from "node:http";

const server = http.createServer((_, res) => {
  res.end("Hello, World!");
});
server.listen(8100);
