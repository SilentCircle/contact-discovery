Contact discovery server
========================

This contact discovery server is a private... contact discovery server. Contact
discovery is the process through which services like WhatsApp figure out which
of your friends also have WhatsApp, so they can show you that you can talk to
them.

Usually, the process involves reading your whole address book, sending it off to
a server somewhere (probably to the NSA) for each contact to be looked up. The
contacts that turn out to actually have the app installed are then returned,
and the app knows that you can send those people a message.

At Silent Circle, we didn't want to know who your contacts are, for two reasons:

* It's none of our business.
* If we don't know who your friends are, we can't tell the NSA when they beat us
  with a rubber hose.

To that end, Phil Zimmermann developed a protocol (we call it the Silent Circle
Contact Discovery Protocol, or SCCDP for unpronounceable) for performing contact
discovery without sending the server all your contacts. This readme is both
a description of the protocol and a documentation of the prototype server
included here.

Before getting to the protocol, some information about the server: It's written
in Go, designed to have no dependencies (you just copy the binary to the target
machine) and can be run by independent third parties (for example, the EFF). The
security guarantees are such that the server *may be able* to discover the
identifiers of the users running on the service, but it cannot discover the
identifiers of the contacts on the device.


The protocol
------------

### Description

The protocol is pretty simple:

1. The client compiles a list of identifiers for all the people it wants to look
   up (usernames, emails, phone numbers, whatever the user identifiers for the
   service are).
2. It hashes each identifier with a pre-agreed-upon appropriate hashing function
   (SHA is probably a good choice).
3. The client then *truncates* each hash to its first N characters (selectable
   by the client, but usually 4+ bits) and sends the truncated hashes to the
   server.
4. The server replies with a list of slightly less truncated hashes (around 20
   characters) that it knows about and that begin with the characters the client
   sent.
5. The client then compares each less-truncated hash with the original, and,
   if all the characters match, it can be reasonably sure that the server knows
   the identifier of the user the client sent.


### Why this is private

The point of the protocol above is that the client doesn't send the server the
entire identifier (the space of all possible identifiers is so small that a hash
provides pretty much no privacy, it's just a convenient way to segment the
keyspace). After all the hashing, this is pretty much the equivalent of wanting
to see if the server knows the email address "stavros@silentcircle.com", and
sending "sta" to have the server reply with "stavros@silentcir", at which point
you're pretty sure the server knows about it.

However, since you only sent "sta", the server can't know whether you meant
"stavros@silentcircle.com" or "stakayama@softbank.jp". Thus, your privacy is
preserved, while the contact discovery process proceeds as usual.


### Considerations

There are a few considerations in the lengths that the client and server send.
The longer the hash that the client sends, the less privacy it has, but also the
less data it receives. The server will probably want to enforce a minimum length
to avoid DoS attacks by people asking for every hash starting with a single
character, but the client also probably does not want to specify a hash prefix
so general that it receives 100 hashes per contact.

Conversely, the server may want to weak its reply truncation length as well
(although it's less important, since the server doesn't care about hiding which
users are on the service as much). That will depend on how many users are on the
service and how comfortable the server is with sending long hashes.


The server
----------

The provided server implementation is written in Go, mainly for speed and ease
of deployment. It's pretty simple, it uses SQLite to store and query the list of
hashes and an HTTP API for communicating with the outside world. Here's how that
works:


### Installation and execution

To install and run the server, just do the following:

~~~
go get github.com/SilentCircle/contact-discovery

./contactdiscovery aiGh8ohLeewoo4iz
~~~

This will start the server on port 8080. The server has two endpoints, an
authenticated one for adding and deleting hashes, and an unauthenticated one to
look them up. Here they are:


### Adding a hash

Adding a hash can be done with a simple POST. You need to use the password you
started the server with as the HTTP Basic auth password:

~~~
http POST "http://:aiGh8ohLeewoo4iz@localhost:8080/hashes/5891b5b522d5df086d0ff0b110fbd9d21bb4fc7163af34d08286a2e846f6be03/"
{
    "result": "success"
}
~~~

### Deleting a hash

Exactly the same as above, except now it's a DELETE. Simple.

~~~
http DELETE "http://:aiGh8ohLeewoo4iz@localhost:8080/hashes/5891b5b522d5df086d0ff0b110fbd9d21bb4fc7163af34d08286a2e846f6be03/"
{
    "result": "success"
}
~~~

### Looking up hashes

To look up hashes, you don't need a password. you need a POST with some JSON in
the body. The JSON should look like this:

~~~.json
{
    "prefixes": [
        "5891"
    ]
}
~~~

Here's the full request:

~~~
16:54:15 $ http POST "http://localhost:8080/contacts/" prefixes:='["5891"]'
{
    "hashes": [
        "5891b5b522d5df086d0f"
    ],
    "result": "success"
}
~~~

That's pretty much the entire API.


Boring stuff
------------

The server is released under the Apache 2.0 license. You're free to use the
protocol however you like, but try not to use it for evil (that's a polite
request, not a stipulation).
