# pepper

Pepper is a small CLI tool written in Go to run Salt commands against the HTTP API. It's very thin and there's almost no logic to the script.

Both PAM and LDAP auth methods are supported.

## Installation

```bash
go get github.com/iamseth/pepper
go install github.com/iamseth/pepper
```

## Usage

```bash
NAME:
   pepper - pepper <target> <function> [ARGUMENTS ...]

USAGE:
   pepper [global options] command [command options] [arguments...]

VERSION:
   0.1.0

COMMANDS:
   help, h	Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --hostname, -H 	Salt API hostname. Should include http[s]//. [$SALT_HOST]
   --username, -u 	Salt API username. [$SALT_USER]
   --password, -p 	Salt API password. [$SALT_PASSWORD]
   --auth, -a "ldap"	Salt authentication method. [$SALT_AUTH]
   --help, -h		show help
   --version, -v	print the version
```

## Examples
```bash
pepper -H https://saltmaster -u salt_user -p supercoolpassword '*' test.ping
```
Or, if you set your environment variables
```bash
SALT_USER
SALT_HOST
SALT_AUTH
SALT_PASSWORD
```
you can just do the following:
```bash
# Ping all your hosts
pepper '*' test.ping

# Restart Apache on your web servers
pepper 'web*' cmd.run 'service httpd restart'

# Run a highstate on your Redis boxes
pepper 'redis*' state.highstate
```
