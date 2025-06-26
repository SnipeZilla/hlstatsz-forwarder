## hlstatsz-forwarder

**Multi-server CS2 log forwarder** that converts HTTP log streams from CS2 servers into HLstatsX/Z-compatible UDP packets.<br>
Supports multiple endpoint paths, hot-reloadable config, and customizable log levels.

---

## 🚀 Features

- 🔁 **Path-based routing** for multiple CS2 servers
- 📦 **HLstatsX/Z-compatible UDP output** with `PROXY` spoofing
- 🔧 **Live config reloads** via [`viper`](https://github.com/spf13/viper)
- 🔍 Adjustable **log verbosity**: `debug`, `info`, `silent`
- ⚙️ Built in Go: portable, fast, and zero dependencies beyond `viper`

---

## 🛠 Configuration

Create a `config.yaml` in the same directory as the binary:

```yaml
log_level: "debug"                   # options: debug, info, silent
http_port: 29999                     # port the forwarder listens on
udp_address: "127.0.0.1:27500"       # HLstatsX/Z daemon destination
proxy_key: "yourkey"                 # HLstatsX/Z ProxyKey setting

servers:
  - path: "/server1"                 # URL path for incoming logs
    ip: "64.74.97.164"               # Spoofed IP for HLstatsX/Z
    port: "27015"
  - path: "/server2"
    ip: "64.74.97.164"
    port: "27016"
```
## 💡 How It Works
Each CS2 server sends logs to a unique HTTP endpoint:
<br>
<br>
server.cfg:
 <br>
  `logaddress_add_http "http://hlstats-server-ip:29999/server1"`
<br>
<br>
/server1 is a header identifier<br>
CS2 doesn't include its own IP or game port in the request headers or body by default<br>


The forwarder receives the log, wraps it like this:<br>
L  06/25/2025 - 18:40:00: "CATS<1><BOT><CT>" say "I love SnipeZilla"<br>
...and sends it via UDP to your HLstatsX/Z daemon.<br>


## 📦 Building & Running
windows:<br>
go build -o hlstatsz-forwarder.exe<br>
hlstatsz-forwarder.exe<br>
<br>
linux:<br>
go build -o hlstatsz-forwarder<br>
./hlstatsz-forwarder<br>


## 🔐 HLstatsX/Z Setup
Ensure these match:

The spoofed IP:Port in config.yaml

The HLstatsX/Z game server list (in DB or web admin panel)

ProxyKey in HLstatsX/Z's hlstats.conf

Rejected log lines will show up like:<br>
E997: NOT ALLOWED SERVER

## 🐛 Debugging Tips
Logs not showing up in HLstatsX/Z? Run with log_level: "debug" and look for:

Firewall inbound rule for the TCP `http_port` xxxxx (if external address)

Incoming POST requests

Correct path routing

UDP send confirmations

## 📄 License
MIT — feel free to fork, contribute, or use this in your own stack.

## 🌠 Authored
<img src="https://snipezilla.com/hlstatsz/styles/css/images/Z.png" style="width:14px"> SnipeZilla Sniping Servers [<a href="https://snipezilla.com">https://snipezilla.com</a>]
