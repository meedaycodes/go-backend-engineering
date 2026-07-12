# Linux Commands Cheatsheet

## Filesystem

| Command | Description |
|---------|-------------|
| `ls /` | List root filesystem directories |
| `ls -la` | List files with permissions, owner, size, date |
| `cat /etc/hostname` | Print the machine's hostname |
| `cat /etc/hosts` | Show hostname-to-IP mappings |
| `cat /etc/passwd` | List all system users |
| `cat /etc/resolv.conf` | Show DNS server configuration |
| `cat /etc/environment` | Show system-wide environment variables |

## Permissions

| Command | Description |
|---------|-------------|
| `chmod 600 file` | Owner read/write only (`rw-------`) |
| `chmod 500 file` | Owner read/execute only (`r-x------`) |
| `chmod 400 file` | Owner read only (`r--------`) |
| `chmod 700 dir` | Owner full access, others none (`rwx------`) |
| `chown user:group file` | Change file owner and group |
| `ls -la` | View permissions in `rwxr-xr-x` format |

### Permission Bits
```
r = 4, w = 2, x = 1
rwx = 7, rw- = 6, r-x = 5, r-- = 4, --- = 0
```

### Permission String Format
```
d  rwx  r-x  r-x
|   |    |    |
|   |    |    └── others
|   |    └─────── group
|   └──────────── owner
└──────────────── type: d=directory, -=file, l=symlink
```

## Users & Groups

| Command | Description |
|---------|-------------|
| `useradd --no-create-home --shell /usr/sbin/nologin appuser` | Create a service account |
| `id username` | Show user's UID, GID, and groups |
| `whoami` | Show current user |
| `groups` | Show groups the current user belongs to |

## Process Management

| Command | Description |
|---------|-------------|
| `ps aux` | List all running processes |
| `kill <pid>` | Send SIGTERM (graceful shutdown) to a process |
| `kill -9 <pid>` | Send SIGKILL (force kill) to a process |
| `sleep 10 &` | Run a command in the background |

### Key Signals
- `SIGTERM` — asks the process to shut down gracefully, can be caught
- `SIGKILL` — forces immediate termination, cannot be caught or ignored

## systemd Service Management

| Command | Description |
|---------|-------------|
| `sudo systemctl start <service>` | Start a service |
| `sudo systemctl stop <service>` | Stop a service |
| `sudo systemctl restart <service>` | Restart a service |
| `sudo systemctl enable <service>` | Start service automatically on boot |
| `sudo systemctl disable <service>` | Remove from boot startup |
| `sudo systemctl status <service>` | Show service status and recent logs |
| `sudo systemctl daemon-reload` | Reload unit files after editing |
| `sudo journalctl -u <service> -n 50` | View last 50 log lines for a service |
| `sudo journalctl -u <service> -f` | Follow live logs for a service |

## Networking

| Command | Description |
|---------|-------------|
| `curl http://localhost:8080/health` | Send HTTP GET request |
| `curl -X POST url -H "Content-Type: application/json" -d '{}'` | Send HTTP POST with JSON body |
| `ss -tlnp` | Show listening TCP ports and which process owns them |
| `ping <host>` | Test network connectivity to a host |

## SSH & SCP

| Command | Description |
|---------|-------------|
| `ssh -i key.pem user@ip` | SSH into a remote server |
| `scp -i key.pem file user@ip:/remote/path` | Copy file to remote server |
| `scp -i key.pem user@ip:/remote/file ./local` | Copy file from remote server |

## Cross-Compilation (Go)

```bash
# Build for Linux x86-64 from any OS
GOOS=linux GOARCH=amd64 go build -o bin/server-linux ./cmd/server

# Verify the binary target platform
file bin/server-linux
```

## systemd Unit File Structure

```ini
[Unit]
Description=My Go Service
After=network.target          # start only after network is up

[Service]
Type=simple
User=goapp                    # run as non-root service account
ExecStart=/usr/local/bin/server
Restart=on-failure            # restart automatically if it crashes
RestartSec=5                  # wait 5 seconds before restarting
Environment=PORT=8080         # inject environment variables

[Install]
WantedBy=multi-user.target    # start in normal multi-user mode
```
