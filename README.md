# Sshwifty Web SSH & Telnet Client

**Sshwifty is a SSH and Telnet client made for the Web,** allow you to access
SSH and Telnet services right from your web browser.

![Screenshot](Screenshot.png)

(The glass glare effect above is only included in the Executive Golden Premium
Plus+ Platinum Ultimate AD-free version, which can be obtained after joining the
cult. Though, science has proven that the normal AD-free version is sufficient
for most people. In fact, the majority of people surveyed are annoyed by the
glare, while the rest showed a little interest)

![Build Status](https://github.com/nirui/sshwifty/workflows/Sshwifty-CI/badge.svg)

## Install

### Prebuilt Executables

We offer prebuilt ready-to-run programs (called executables) for a handful of
popular platforms in the [Releases] section. You might find one built for your
platform there.

Please aware that these executables are generated by an unsupervised automatic
procedure, and we (as the authors and contributors) cannot guarantee to test
them. If you've encountered unusual failure caused by those executables, feel
free to open an issue, so we can take a look.

[releases]: https://github.com/nirui/sshwifty/releases

### Docker Image (recommended)

Deploying Sshwifty as a [Docker] container allows you to isolate Sshwifty from
the rest of the system for better system organization and security.

We also offer prebuilt Docker Images for few selected platforms generated in the
same way as we generate the prebuilt Executables described above. To deploy one
on your Docker host, run:

```shell
$ docker run --detach \
  --restart unless-stopped \
  --publish 8182:8182 \
  --name sshwifty \
  niruix/sshwifty:latest
```

(Note: it's `niruix/sshwifty` with an `x`)

This will open port `8182` on the Docker host to accept traffic from all
clients (including remote ones), and serve them with the Sshwifty instance
just created.

To expose Sshwifty locally only to the Docker host, change the line
`--publish 8182:8182` from the command above to `--publish 127.0.0.1:8182:8182`.
This is useful for scenarios where you want to force remote clients to access
Sshwifty only though a reverse proxy, or just simply want to prevent remote
access.

When TLS is desired and you don't want to setup Docker Volumes, you can use
`SSHWIFTY_DOCKER_TLSCERT` and `SSHWIFTY_DOCKER_TLSCERTKEY` environment variables
to import certificate files to the container and automatically apply them:

```shell
$ openssl req \
  -newkey rsa:4096 -nodes -keyout domain.key -x509 -days 90 -out domain.crt
$ docker run --detach \
  --restart always \
  --publish 8182:8182 \
  --env SSHWIFTY_DOCKER_TLSCERT="$(cat domain.crt)" \
  --env SSHWIFTY_DOCKER_TLSCERTKEY="$(cat domain.key)" \
  --name sshwifty \
  niruix/sshwifty:latest
```

The `domain.crt` and `domain.key` in the command above is the location of valid
X509 certificate and key file.

Though, in most situations where a reverse proxy and/or load balancer (for
example, [Nginx] or [Traefik]) is used in front of Sshwifty instances, TLS
should usually terminate on the proxy, not on the individual Sshwifty instances.

[Docker]: https://www.docker.com
[Nginx]: https://github.com/nginx/nginx
[Traefik]: https://github.com/traefik/traefik

### Compile from source code (recommended for developers)

The following tools are required in order to compile the software from source
code:

- `git` to download the source code
- `node` and `npm` to build front-end application
- `go` to build back-end application

To start the build process, run:

```shell
$ git clone https://github.com/nirui/sshwifty
$ cd sshwifty
$ npm install
$ npm run build
```

When done, you can find the newly generated `sshwifty` binary under current
working directory.

Notice: `Dockerfile` contains the entire build procedure of this software.
Please refer to it when you encounter any compile/build related issue.

### Third-party Homebrew Formulae from [@unbeatable-101]

If you're a macOS user, [@unbeatable-101] is kindly hosting a Homebrew
Formulae that allows you to install his custom Sshwifty builds for macOS via
`homebrew`. You can hop over to [unbeatable-101/homebrew-sshwifty] for detailed
instruction and contribute to his work.

Please note that, due to the third-party nature of the work, the author(s) of
Sshwifty are unable to provide any audit, warranty or support for it. If you
have any question or request regarding to the Formulae, please contact
[@unbeatable-101] directly through appreciate channels.

Thank [@unbeatable-101] for his work.

[@unbeatable-101]: https://github.com/unbeatable-101
[unbeatable-101/homebrew-sshwifty]: https://github.com/unbeatable-101/homebrew-sshwifty

## Configuration

Sshwifty can be configured through either file or environment variables. By
default, the configuration loader will try to load file from default paths
first, when all failed, environment variables will be used.

You can also specify your own configuration file with `SSHWIFTY_CONFIG`
environment variable. For example:

```shell
$ SSHWIFTY_CONFIG=./sshwifty.conf.json ./sshwifty
```

This tells Sshwifty to only load configuration from file `./sshwifty.conf.json`.

### Configuration file option and descriptions

Here is all the options of the configuration file:

```jsonc
{
  // HTTP Host. Keep it empty to accept request from all hosts, otherwise, only
  // specified host is allowed to access
  "HostName": "localhost",

  // Web interface access password. Set to empty to allow public access to the
  // web interface (By pass the Authenticate page)
  "SharedKey": "WEB_ACCESS_PASSWORD",

  // Remote dial timeout. This limits how long of time the backend can spend
  // to connect to a remote host. The max timeout will be determined by
  // server configuration (ReadTimeout).
  // (In Seconds)
  "DialTimeout": 10,

  // Socks5 proxy. When set, Sshwifty backend will try to connect remote through
  // the given proxy
  "Socks5": "localhost:1080",

  // Username of the Socks5 server. Please set when needed
  "Socks5User": "",

  // Password of the Socks5 server. Please set when needed
  "Socks5Password": "",

  // Server side hooks, allowing operator to launch external processes on the
  // server side to influence server behaver
  //
  // The operation of a Hook must be completed within the time limit defined
  // by `HookTimeout` set below. Otherwise it will be terminated, and results
  // a failure for the execution
  //
  // To determine how much time is still left for the execution, a Hook can
  // fetch the deadline information from the `SSHWIFTY_HOOK_DEADLINE`
  // environment variable which is a RFC3339 formatted date string indicating
  // after what time the termination will occur
  //
  // Warning: the process will be launched within the same context and system
  // permission which Sshwifty is running under, thus is it crucial that the
  // Hook process is designed and operated in a secure manner, otherwise
  // SECURITY VULNERABILITY (commandline injection, for example) maybe created
  // as result
  //
  // Warning: all inputs passed by Sshwifty to the hook process must be
  // considered unsanitized, and must be sanitized by each hook themselves
  "Hooks": {
    // before_connecting is called before Sshwifty starts to connect to a remote
    // endpoint. If any of the Hook process exited with a non-zero return code,
    // the connection request is aborted
    //
    // This Hook offers two parameters:
    // - SSHWIFTY_HOOK_REMOTE_TYPE: Type of the connection (i.e. SSH or Telnet)
    // - SSHWIFTY_HOOK_REMOTE_ADDRESS: Address of the remote host
    "before_connecting": [
      // Following example command launches a `/bin/sh` to execute a for loop
      // that prints to Stdout as well as to Stderr
      //
      // Prints to Stdout will be sent to the client side visible to the user,
      // and prints to Stderr will be captured as server side logs and it is
      // invisible to the user (as server logs usually are)
      //
      // The command must be specified in Json array format. Each array element
      // is mapped to a command fragment separated by space. For example:
      // ["command", "-i", "Hello World"] will be mapped to `command -i "Hello
      // World"` before it is executed
      [
        "/bin/sh",
        "-c",
        "for n in $(seq 1 5); do sleep 1 && echo Stdout $SSHWIFTY_HOOK_REMOTE_TYPE $n && echo Stderr $SSHWIFTY_HOOK_REMOTE_TYPE $n 1>&2; done"
      ],
      // You can add multiple hooks, they're executed in sequence even when the
      // previous one fails
      [
        "/bin/sh",
        "-c",
        "/etc/sshwifty/before_connecting.sh"
      ],
      [
        "/bin/another-command",
        "...",
        "..."
      ]
    ]
  },

  // The maximum execution time of each hook, in seconds. If this timeout is 
  // exceeded, the hook will be terminated, and thus cause a failure
  "HookTimeout": 30,

  // Sshwifty HTTP server, you can set multiple ones to serve on different
  // ports
  "Servers": [
    {
      // Which local network interface this server will be listening
      "ListenInterface": "0.0.0.0",

      // Which local network port this server will be listening
      "ListenPort": 8182,

      // Timeout of initial request. HTTP handshake must be finished within
      // this time
      // (In Seconds)
      "InitialTimeout": 3,

      // How long do the connection can stay in idle before the backend server
      // disconnects the client
      // (In Seconds)
      "ReadTimeout": 60,

      // How long the server will wait until the client connection is ready to
      // recieve new data. If this timeout is exceed, the connection will be
      // closed.
      // (In Seconds)
      "WriteTimeout": 60,

      // The interval between internal echo requests
      // (In Seconds)
      "HeartbeatTimeout": 20,

      // Forced delay between each request
      // (In Milliseconds)
      "ReadDelay": 10,

      // Forced delay between each write
      // (In Milliseconds)
      "WriteDelay": 10,

      // Path to TLS certificate file. Set empty to use HTTP
      "TLSCertificateFile": "",

      // Path to TLS certificate key file. Set empty to use HTTP
      "TLSCertificateKeyFile": "",
      
      // Display a short text message on the Home page. Link is supported 
      // through `[Title text](https://link.example.com)` format
      "ServerMessage": ""
    },
    {
      "ListenInterface": "0.0.0.0",
      "ListenPort": 8182,
      "InitialTimeout": 3,
      .....
    }
  ],

  // Remote Presets, the operater can define few presets for user so the user
  // won't have to manually fill-in all the form fields
  //
  // Presets will be displayed in the "Known remotes" tab on the Connector
  // window
  //
  // Notice: You can use the same JSON value for `SSHWIFTY_PRESETS` if you are
  //         configuring your Sshwifty through enviroment variables.
  //
  // Warning: Presets Data will be sent to user client WITHOUT any protection.
  //          DO NOT add any secret information into Preset.
  //
  "Presets": [
    {
      // Title of the preset
      "Title": "SDF.org Unix Shell",

      // Preset Types, i.e. Telnet, and SSH
      "Type": "SSH",

      // Target address and port
      "Host": "sdf.org:22",

      // Define the tab and background color of the console in RGB hex format
      // for better visual identification
      //
      // For example: 110000 will give you a dark red background, 001100 is
      // dark green and 000011 is dark blue
      //
      // The color must not be too bright, as it will make the foreground text
      // hard to read
      "TabColor": "112233",

      // Form fields and values, you have to manually validate the correctness
      // of the field value
      //
      // Defining a Meta field will prevent user from changing it on their
      // Connector Wizard. If you want to allow users to use their own settings,
      // leave the field unsetted
      //
      // Values in Meta are scheme enabled, and supports following scheme
      // prefixes:
      // - "literal://": Text literal (Default)
      //                 Example: literal://Data value
      //                          (The final value will be "Data value")
      //                 Example: literal://file:///tmp/afile
      //                          (The final value will be "file:///tmp/afile")
      // - "file://": Load Meta value from given file.
      //              Example: file:///home/user/.ssh/private_key
      //                       (The file path is /home/user/.ssh/private_key)
      // - "environment://": Load Meta value from an Environment Variable.
      //                    Example: environment://PRIVATE_KEY_DATA
      //                    (The name of the target environment variable is
      //                    PRIVATE_KEY_DATA)
      //
      // All data in Meta is loaded during start up, and will not be updated
      // even the source already been modified.
      //
      "Meta": {
        // Data for predefined User field
        "User": "pre-defined-username",

        // Data for predefined Encoding field. Valid data is those displayed on
        // the page
        "Encoding": "pre-defined-encoding",

        // Data for predefined Password field
        "Password": "pre-defined-password",

        // Data for predefined Private Key field, should contains the content
        // of a Key file
        "Private Key": "file:///home/user/.ssh/private_key",

        // Data for predefined Authentication field. Valid values is what
        // displayed on the page (Password, Private Key, None)
        "Authentication": "Password",

        // Data for server public key fingerprint. You can acquire the value of
        // the fingerprint by manually connect to a new SSH host with Sshwifty,
        // the fingerprint will be displayed on the Fingerprint comformation
        // page.
        "Fingerprint": "SHA256:bgO...."
      }
    },
    {
      "Title": "Endpoint Telnet",
      "Type": "Telnet",
      "Host": "endpoint.nirui.org:23",
      "Meta": {
        // Data for predefined Encoding field. Valid data is those displayed on
        // the page
        "Encoding": "utf-8"
        ....
      }
    },
    ....
  ],

  // Allow the Preset Remotes only, and refuse to connect to any other remote
  // host
  //
  // NOTICE: You can only configure OnlyAllowPresetRemotes through a config
  //         file. This option is not supported when you are configuring with
  //         environment variables
  "OnlyAllowPresetRemotes": false
}
```

`sshwifty.conf.example.json` is an example of a valid configuration file, you
can make your own customization base on it.

### Environment variables

Valid environment variables are:

```
SSHWIFTY_HOSTNAME
SSHWIFTY_SHAREDKEY
SSHWIFTY_DIALTIMEOUT
SSHWIFTY_SOCKS5
SSHWIFTY_SOCKS5_USER
SSHWIFTY_SOCKS5_PASSWORD
SSHWIFTY_HOOK_BEFORE_CONNECTING
SSHWIFTY_HOOKTIMEOUT
SSHWIFTY_LISTENPORT
SSHWIFTY_INITIALTIMEOUT
SSHWIFTY_READTIMEOUT
SSHWIFTY_WRITETIMEOUT
SSHWIFTY_HEARTBEATTIMEOUT
SSHWIFTY_READDELAY
SSHWIFTY_WRITEELAY
SSHWIFTY_LISTENINTERFACE
SSHWIFTY_TLSCERTIFICATEFILE
SSHWIFTY_TLSCERTIFICATEKEYFILE
SSHWIFTY_SERVERMESSAGE
SSHWIFTY_PRESETS
SSHWIFTY_ONLYALLOWPRESETREMOTES
```

These options are correspond to their counterparts in the configuration file.

Notice: When you're using environment variables to configure Sshwifty, only one
Sshwifty HTTP server is then allowed. There is no way to setup multiple servers
under this method of configuration. If you need to serve on multiple ports, use
the configuration file instead.

Be aware: An invalid value inside following environment variables will cause
the value to be silently reset to default during configuration parsing phase
without warning:

```
SSHWIFTY_DIALTIMEOUT
SSHWIFTY_INITIALTIMEOUT
SSHWIFTY_READTIMEOUT
SSHWIFTY_WRITETIMEOUT
SSHWIFTY_HEARTBEATTIMEOUT
SSHWIFTY_READDELAY
SSHWIFTY_WRITEELAY
```

Please verify the value of these options before start the instance.

## FAQ

### Why the software says "The datetime difference ... is beyond tolerance"?

Currently, Sshwifty implemented a layer of obscuration in it's wire protocol.
The protocol utilize the common time of both endpoint to generate cryptographic
keys to be used to encrypt and decrypt traffic between the client and the
server. Thus both end must have accurate datetime relative to the time of the
world.

Please make sure the datetime on both the client and the server are correct by
resync them with a NTP server, and then simply reload the page. The problem
should be gone afterwards.

### Why I got error "TypeError: Cannot read property 'importKey' of undefined"

Sshwifty's wire protocol requires few cryptographic features which is only
available with WebCrypt API. Some old web browsers maybe not support WebCrypt
API, or the newer ones may disable WebCrypt API when the web page is not
accessed under [Secure contexts], which led to the error message.

[Secure Contexts]: https://developer.mozilla.org/en-US/docs/Web/Security/Secure_Contexts

Usually, setup HTTPS either on your reverse proxy or on Sshwifty itself should
allow you to access Sshwifty via HTTPS (a `https://` URL), and thus solving this
issue. Also, enabling HTTPS is mostly a good idea anyways, even in a LAN.

### Can I serve Sshwifty under a subpath such as `https://my.domain/ssh`?

The short story is No. Sshwifty was designed based on an assumption that it will
run as the only service under a given hostname, allowing web browsers to better
enforce their data isolation rules. This is very important because Sshwifty
saves user data locally, right on the web browser.

However, if you really want to put Sshwifty into a subpath, you can do so by
taking advantage of the fact that Sshwifty backend interface and assets are
always located under an URL prefix `/sshwifty`. You can thus redirect or proxy
those requests to their new location.

Keep in mind, doing so is really hacky, and it's not recommended by the author
thus no support will be provided if you decide to do so.

### Why I can't add my own key combinations to the Console tool bar?

The pre-defined key combinations are there mainly to make mobile operation
possible as well as to resolve some hotkey conflicts. However, if efficiency is
your first goal, please consider to use a software/on screen keyboard which is
specially designed for terminal.

And if that's not enough, connect a physical keyboard through Bluetooth or OTA
could be a better alternative. This way you can type as if you're using a
computer console.

## Credits

- Thanks to [Ryan Fortner](https://github.com/ryanfortner) for the grammar fix
- Thanks to [Tweak](https://github.com/Tweak4141) for the grammar fix
- Thanks to [CJendantix](https://github.com/CJendantix) for the grammar and typo 
  fix
- Thanks to [ZStrikeGit](https://github.com/ZStrikeGit) for the grammar and 
  formatting fixes

## License

Code of this project is licensed under AGPL, see [LICENSE.md] for detail.

Third-party components used by this project are licensed under their respective
licenses. See [DEPENDENCIES.md] to learn more about dependencies used by this
project and read their copyright statements.

[LICENSE.md]: LICENSE.md
[DEPENDENCIES.md]: DEPENDENCIES.md

## Contribute

This is a hobbyist project, meaning I don't have that much time to put into it, 
sorry.

Upon release (Which is then you're able to read this file), this project will
enter the _maintaining_ state, which includes doing some bug fixes and security 
updates. _Adding new features however, is not a part of the state_.

Please do not send any pull requests. If you need new feature, fork it, add it 
by yourself, then maintain it like one of your own project. It's not that hard 
with some Github features.

(Notice: Typos, grammar errors or invalid use of language in the source code and
documentation should be categorized as a bug, please report them if you find 
any. Thank you!)

Appreciate your help, enjoy!
