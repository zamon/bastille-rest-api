# Prerequisites:
- sudo
- bastille 0.14
- edit /usr/local/etc/sudoers, insert this line below:
[your username] ALL=(ALL) NOPASSWD:/usr/local/bin/bastille

# How to build from source:
- clone this repo
- cd to cloned directory
- mkdir build
- go mod tidy
- go build -o build/server (or any name you like)

# How to use:
- upload binary to your home in FreeBSD machine. then change .env.example to .env:
    - default app port is 8000
    - Whitelist is used for ignoring request from other machine, set IP address whitelist using your client machine IP and if you have multiple IP, separate with commas without space, for example IP_WHITELIST=1.1.1.1,2.2.2.2 etc. 

# Available end point
## Package list
URL: POST http://[your server IP]:8000/bastille-pkg-list

Body Request (json):
{
    "jail": "[your jail name]",
    "package": "[package to search]"
}

response: list all available package in jail

## List all jail
URL: GET http://[your server IP]:8000/bastille-list-all

response: list all jail in machine

## List all bootstrap
URL: GET http://[your server IP]:8000/bootstrap-list

response: all available bootstrap for jail

## Ping
URL: GET http://[your server IP]:8000

response: pong

## Create jail
URL: POST http://[your server IP]:8000/create-jail

Body Request (json):
{
    "jail_name": "[jail name]",
    "ip_address": "[jail IP address]",
    "release": "[release]"
}

response: json output

## Stop jail
URL: POST http://[your server IP]:8000/stop-jail

Body Request (json):
{
    "jail_name": "[jail name]"
}

response: json output

## Destroy jail
URL: POST http://[your server IP]:8000/destroy-jail

Body Request (json):
{
    "jail_name": "[jail name]"
}

response: json output