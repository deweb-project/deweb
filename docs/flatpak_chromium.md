# deweb with chromium installed in flatpak

If you have installed UngoogledChromium thru flatpak,

Create a file `/usr/bin/chromium`, and put following content inside:

```sh
#!/bin/sh
flatpak run com.github.Eloston.UngoogledChromium $@
exit $!
```

and run `sudo chmod +x /usr/bin/chromium`