prerequisites:
- sudo
- bastille 0.14
- edit /usr/local/etc/sudoers, insert this line below
<your username> ALL=(ALL) NOPASSWD:/usr/local/bin/bastille

how to build from source:
- mkdir build
- go mod tidy
- go build -o build/server (or any name you like)

how to use:
- upload binary to your home in FreeBSD machine. then edit .env:
    - default app port is 8000
    - Whitelist is used for ignoring request from other IP, set ip whitelist using your client machine ip and if you have multiple IP, separate with commas without space, for example IP_WHITELIST=1.1.1.1,2.2.2.2 etc. 