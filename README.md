Jellfin-cli
-----------

Access jellfin through CLI. This is useful in case when browser
is unable to decode particular stream and server is not
configured for transcoding. The cli allows navigation through
the media stored on the server and can open the stream
directly in the vlc media player

Usage
-----
A config.json is expected in `~/.config/jellycli.conf` of binary with below structure

    {
        "AuthKey":"<authkey>",  // can be obtained from API keys section in dashboard
        "Host":"<server-url>",  // without the trailing slash
        "UserId":"<userid>" // can be obtained from browser API requests
    }

Screenshots
-----

![image](https://github.com/user-attachments/assets/76ee7394-3136-4562-a120-fb24efb92f10)

![image](https://github.com/user-attachments/assets/1df818b9-87ac-4a7e-8152-acef7430ab1e)



