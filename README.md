## Purpose

RedGrittyBrick mbox is a software library written in the language "go" (golang)
which can be used to read the e-mail files in the "mbox" format that are
created by the e-mail program Thunderbird.

Thunderbird creates a local cache of e-mail massages, even for remote email
servers which use the IMAP protocol.

You can use this library to read those cached emails, or old or backed-up
copies.  This can be useful, for example, if your PC has died and you need to
recover files from the hard disk or a backup. You can then use this simpler
program to review messages without needing to reinstall Thunderbird and without
any risk of Thunderbird overwriting the files or performing any clean up that
might delete old messages.

## Usage

See `example/` folder

```
parser := mbox.New()
msgs, err := parser.IndexFile(fileName)       // index of messages, no contents

strt, end := msgs[n].Start(), msgs[n].End()
text, err := mbox.GetMsg(fileName, strt, end) // raw text
cooked, err := mbox.CookMessage(text)         // readable contents
outline, plain := mbox.Outline(text)          // multipart message structure
```
