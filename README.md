Impetus
===

Impetus finds pull requests which have been forgotten and reminds us to go review them. In a world where teams have more and more project repositories, its mission is to eliminate the forgotten pull request.

Building
---

##### Preferred Method

Assuming you have a working Go development environment,

```
go install github.com/vickleford/impetus/impetus/cmd/impetus
```

##### Another Way

Alternatively, a quickie build can be done like this:

```
go build -o impetuscli impetus/cmd/impetus/main.go
```

##### Without Go

```
make run
```

Will build the package, and then copy the cli file to the /artifacts folder.
You can then copy that impetuscli to /usr/local/bin or some other path and link it to systemd as explained below.

Configuration
---

Some things are needed to pull everything together:

- GHE personal access token (ro should be enough) exported as env var `IMPETUS_GIT_TOKEN`
- HipChat room notification token exported as env var `IMPETUS_HIPCHAT_ROOM_TOKEN`
- Hipchat room's API ID passed with -roomid ID flag
- GHE Base URL (eg https://github.company.com) as env var `GHE_URL`

Usage
---

To get it to run, export `IMPETUS_GIT_TOKEN` and `IMPETUS_HIPCHAT_ROOM_TOKEN` then hit the CLI:

```
impetus -roomid ROOMID -repo org/repo... -indefinite
```

Omitting -indefinite will cause it to run once and exit. 

Example:
```
impetus -roomid 0000000 -repo org2/project1 -repo org2/project2 -indefinite
```

To see more options, run `impetus` with the `-h` or `-help` flag.

Where Can I Run This?
---

Run it from your workstation if there's nowhere else to run it from. It's better than nothing.

##### systemd

Systemd runs services extraordinarily well. My systemd service unit looks like this:

```
# /etc/systemd/system/impetus.service
[Service]
Environment=http_proxy=http://proxy:8080
Environment=https_proxy=http://proxy:8080
Environment=GHE_URL=https://github.company.com
Environment=IMPETUS_GIT_TOKEN=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
Environment=IMPETUS_HIPCHAT_ROOM_TOKEN=yyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyy
ExecStart=/home/vwatkins/go/bin/impetus -indefinite -repo org/project -repo org/project2 -repo org/project3 -roomid 0000000
User=nobody
```

Put that file down, change permissions so that your tokens aren't exposed to anybody with a shell on that system, reload systemd with `systemctl daemon-reload` and run impetus with `systemctl start impetus`.

Spec
---

[x] It should be runnable as a daemon

[x] It should only worry about one github url per instance of impetus running

[x] It should take a list of multiple repositories in the format of org/repo,org/repo2,org2/repo...

[x] It should find pull requests that are open and have had no activity for 24 hours

[x] It should notify us of those pull requests over hipchat

[ ] It should only notify during business hours

[x] The number of hours should be configurable with a sane default

[x] It should be able to scan a specific repository (many teams having many repositories under one organization) or scan an entire organization (a team's repositories are all inside its own organization).

[x] It should handle paginated resposnes from GitHub ;(

[ ] The github URL should be configurable.
