# Filebox

Filebox is a simple solution to the problem of "how do I get this file from my
phone to my laptop?" The traditional solution is to email the file to yourself,
but then you're just using the email server as a file server, and requires
access to your email.

I wrote this as an alternative to my `scp <file> rpi:/home/share/serve/`
with a simple python http client serving static files out of that directory.

This means that there are a few core project requirements:

- High performance. Keep it simple, browsers are really fast if they don't
  have to parse megabytes of Javascript.
- Simple workflow. No complex OAuth2 login workflows, no navigating a million
  different menus, just click upload.

Obviously we still need *some* kind of authentication, but serving HTTP
basic passwords over an encrypted connection is good enough. Since the user
and hoster are one and the same, there's no point in complex E2E encryption
protocols, validation mechanisms, etc. If you willingly brick your own webapp
then that's your problem.

# Installation

Stub

# Development

Stub