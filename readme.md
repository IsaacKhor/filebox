# Filebox

Filebox is a simple solution to the problem of "how do I get this file
from my phone to my laptop?" The traditional solution is to email the file to
yourself, but then you're just using the email server as a file server, and
requires access to your email.

I wrote this as an alternative to my `scp <file> rpi:/home/share/serve/`
with a simple python http client serving static files out of that directory.

This means that there are a few core project requirements:

- Lightweight. Keep it simple, browsers are really fast if they don't
  have to parse megabytes of Javascript.
- Simple workflow. No complex OAuth2 login workflows, no navigating a million
  different menus, just click upload.

Obviously we still need *some* kind of authentication, but serving HTTP
basic passwords over an encrypted connection is good enough. Since the user
and hoster are one and the same, there's no point in complex E2E encryption
protocols, validation mechanisms, etc. If you willingly brick your own webapp
then that's your problem.

There's a small extension to core functionality that allows the generation
of download links to specific files. This is my solution to email file size
limits; instead of complicated GDrive/Dropbox workflows just click a button
and share that link instead.

## Why a self-hosted solution?

Security and privacy. I know that services like [filebin](https://filebin.net/)
exist, but for sensitive data there is no greater assurance of security than
something I host and wrote myself.

The codebase is intentionally kept small and compact to aid in auditing that
it says what it does on the label (not because I'm lazy I swear).

## Other users

In service of the "file too big for email" use case, I'm considering some 
kind of feature where a third party may upload to the same instance so I may 
receive big files from other people who do not host their own servers.

This would require the implementation of:

- Input validation
- Authentication for secondary users
- Auth control for the above (generation of temporary passwords, timeouts)
- Access control for secondary users

The feature set is a little too big at this time to implement, but the app
is designed with expansion for this in mind.

## Installation

```
git clone https://github.com/IsaacKhor/filebox.git
cd filebox
mkdir files
echo '[]' > files/filesdb.json
vim config.json # Set the variables up
go build
./filebox
```

## Configuration

```
DbPath: path to json file containing details about uploaded files
FilesPath: directory to store files
ProductionPort: port to host the server
Host: the HTTP host that incoming requests should have
```
