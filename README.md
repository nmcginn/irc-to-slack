# irc-to-slack
IRC -> Slack forwarding for Eve nerds

### Requirements

As a prerequisite for running this program, you must create an incoming slack webhook that sends messages to you privately.
To set this up, you can go to [https://my.slack.com/services/new/incoming-webhook/](https://my.slack.com/services/new/incoming-webhook/).
Make sure that the webhook is setup to send a private message to you (not a channel), and that you customize the emoji (for street cred).
Save the webhook URL, you'll need this in the next section.

### IRC -> Slack Forwarder Setup

Once the webhook has been created, you'll need to create some environment variables that the program will run on. These variables are:

- WEBHOOK_URL (from slack, in the form of a complete URL)
- IRC_SERVER (the server name including port. for example, irc.freenode.net:7000)
- IRC_NICK (your IRC username)
- IRC_PASSWORD (your IRC password, if applicable)

Once these environment variables have been set, you should be able to download an executable from [releases](https://github.com/odstderek/irc-to-slack/releases) and run.
