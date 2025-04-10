# vangogh

`vangogh` helps you create and manage a personal, DRM-free library of your GOG game collection. It works alongside `theo`(https://github.com/arelate/theo) to simplify downloading, installing and playing games on your devices.

| Dark theme                                       | Light theme                                        |
|--------------------------------------------------|----------------------------------------------------|
| ![List Dark](github_images/list-dark.jpeg)       | ![List Light](github_images/list-light.jpeg)       |
| ![Product Dark](github_images/product-dark.jpeg) | ![Product Light](github_images/product-light.jpeg) |


## Future vangogh ideas and aspirations

`vangogh` roadmap is available [here](ROADMAP.md).

## Installation

The recommended way to install `vangogh` is with [docker compose](https://docs.docker.com/compose/):

- create a `compose.yml` file (this minimal example omits common settings like network, restart, etc):

```yaml
version: '3'
services:
  vangogh:
    container_name: vangogh
    image: ghcr.io/arelate/vangogh:latest
    environment:
    # Set operating systems, languages, etc - customize to your needs
      - VANGOGH_OS=Windows # possible values: Windows, macOS, Linux
      - VANGOGH_LANG-CODE=en # see all possible values in https://github.com/arelate/southern_light/blob/main/gog_integration/languages.go  
      - VANGOGH_NO-PATCHES=true # this disables individual patches downloads
    # Uncomment (remove leading # in the line below) to create a log file on every sync 
    # - VANGOGH_SYNC_DEBUG=true
    # DANGER ZONE: Please read downloads layouts documentation below!
    # Uncomment (remove leading # in the line below) to use a different downloads layout
    # Note: `sharded` is the default behavior (downloads/s/slug) and `flat` is simpler layout (downloads/slug)
    # - VANGOGH_DOWNLOADS-LAYOUT=sharded
    volumes:
      # cold storage - less frequently accessed data, that can be stored on hibernating HDD.
      # hot storage - frequently accessed data, that can benefit from being stored on SSD.
      # backups (cold storage)
      - /docker/vangogh/backups:/var/lib/vangogh/backups
      # downloads (cold storage)
      - /docker/vangogh/downloads:/var/lib/vangogh/downloads
      # checksums (hot storage)
      - /docker/vangogh/checksums:/var/lib/vangogh/checksums
      # images (hot storage)
      - /docker/vangogh/images:/var/lib/vangogh/images
      # input (hot storage)
      - /docker/vangogh:/var/lib/vangogh/input
      # description_images (hot storage)
      - /docker/vangogh/description_images:/var/lib/vangogh/description_images
      # logs (cold storage)
      - /docker/vangogh/logs:/var/lib/vangogh/logs
      # metadata (hot storage)
      - /docker/vangogh/metadata:/var/lib/vangogh/metadata
      # output (hot storage)
      - /docker/vangogh:/var/lib/vangogh/output
      # recycle_bin (cold storage)
      - /docker/vangogh/recycle_bin:/var/lib/vangogh/recycle_bin
      # sharing timezone from the host
      - /etc/localtime:/etc/localtime:ro
      # certificates
      - /etc/ssl/certs/ca-certificates.crt:/etc/ssl/certs/ca-certificates.crt:ro
    ports:
      # https://en.wikipedia.org/wiki/Vincent_van_Gogh
      - "1853:1853"
```
- (move it to location of your choice, e.g. `/docker/vangogh` or remote server or anywhere else)
- while in the directory with that config - pull the image with `docker compose pull`
- start the service with `docker compose up -d`
- assuming everything was setup correctly - `vangogh` will now be available at this location address, port 1853 (e.g. http://localhost:1853 for a local server installation)

## Browser compatibility

`vangogh` has been tested in the most popular browsers (Chrome and other Chromium-based browsers, Safari, Firefox) and should work reasonably well in all of them with minor visual differences.

However, in order to make `vangogh` work in Firefox you'd need to enable [@scope](https://developer.mozilla.org/en-US/docs/Web/CSS/@scope) support if you're comfortable enabling preview Web Platform feature. Here's how to do that (copied from MDN): 

```text
From version 128: this feature is behind the layout.css.at-scope.enabled preference (needs to be set to true). To change preferences in Firefox, visit about:config.
```


## Getting started

After you've installed `vangogh`, you need to authenticate GOG.com access by providing browser cookie data. To do that you need to create `cookies.txt` in the `temporary data` folder (see [docker installation](#Installation)),
   then follow [instructions here](https://github.com/boggydigital/coost#copying-session-cookies-from-an-existing-browser-session) to copy `gog.com` cookies into that file. When you run `vangogh` for the first time, it'll import the cookie header value and split individual parameters.

Regardless of how you do it, the content of `cookies.txt` should look like this:

```text
gog.com
 cart_token=(some value)
 gog-al=(some value)
 gog_lc=(some value)
 gog_us=(some value)
```

If you want to de-authorize `vangogh` from accessing your GOG.com data - delete the `cookies.txt` file. After that you'll still be able to download anything that does not require account authorization. All the account specific data you'll have accumulated until that point will be preserved.

## Updating data

To update your data you can use `vangogh` with a CLI interface to get and maintain all the publicly available GOG.com data, including your account data (game installers):

- in the same folder with `compose.yml` config use `docker compose exec vangogh vangogh <command> <options>`
- most commonly you would run sync `docker compose exec vangogh vangogh sync -all` that gets all available data from GOG.com. 

Sync is optimized to get as little data as possible on each run - only the newly added images and updated installers. There is no great way to determine if metadata was updated remotely, so all of it is fetched on each sync - however upon doing that `vangogh` would know exactly what changed and use this information to optimize decisions. You can also create scheduled runs of synchronization (e.g. with cron) to keep your data fresh. We recommend running sync no more often than every 24 hours. 

### Disk space requirements

Please note that all data === a lot of data, so make sure you have space available or use specific options to get what you need (any combination of `-data -images -screenshots -thumbnails -downloads`)

Here are few estimates of how much space you'll need for each type of data:

- metadata (GOG, Steam, PCGamingWiki, HLTB, ProtonDB data for GOG products): 1Gb
- images, including screenshots: 68Gb
- description items: 19Gb
- product installers estimates (done with `size` command): 
  - 6.5Tb/1000 products for all operating systems (Windows, macOS, Linux) for the 10 most common languages installers, DLCs, extras
  - 4.5Tb/1000 products for all operating systems (Windows, macOS, Linux), English installers, DLCs, extras 
  - 3.3Tb/1000 products for Windows, English installers, DLCs, extras
  - 1.0Tb/1000 products for macOS, English installers, DLCs, extras
  - 0.6Tb/1000 products for Linux, English installers, DLCs, extras

`vangogh` directories allow you to use the best storage type for each data type. For example you might want to store installers, checksums and logs on slower HDDs as those are typically written and read sequentially. Metadata and images on the other hand would benefit from faster SSDs if you're using `vangogh` as a web frontend for your data.

## Downloads layouts

`vangogh` supports two ways product directories will be created in the downloads directory:
- `sharded`, the default behavior. If you don't specify downloads layout - you'll get this layout. Product directories will be placed in a first-letter parent directory under downloads directory. For example for a product with a slug "abc", you'll get "{downloads directory}/a/abc". This is done to minimize the impact of large number of child objects in a directory that might impact certain filesystems. 
- `flat`. In this layout product directories will be placed as-is in the downloads directory. For example for a product with a slug "abc", you'll get "{downloads directory}/abc". This is the default layout used by some other GOG downloads tools and might allow easier integrations. Please note that you might start experiencing performance issues with larger libraries, proceed with caution!

Downloads layouts also impact `local-manual-url` property that stores GOG.com ManualUrl mapping to local files (e.g. `/slug/en1installer0` -> `s/slug/installer.exe`).

In order to change download layouts please use the following steps:
- backup local data (run `backup` on your `vangogh` instance)
- run `relayout-downloads -from {from} -to {to}`, where {from}, {to} are either `flat` or `sharded` depending on your needs
- specify `VANGOGH_DOWNLOADS-LAYOUT={new-layout}` in the `compose.yml` (see the example [above](#Installation))
- recreate the container and restart `vangogh`

## Setting up authentication

`vangogh` requires configured usernames/password to access more sensitive data (e.g. installers downloads). Unless you set those - you won't be able to access them through the CLI. You also need to set those for `theo` to be able to download games.

Here's how to do that - you need to add that to `compose.yml` in the `environment:` section:

```yaml
  # add this under environment: in compose.yml
  - VANGOGH_SERVE_ADMIN-USERNAME=admin-user
  - VANGOGH_SERVE_ADMIN-PASSWORD=admin-password
  - VANGOGH_SERVE_SHARED-USERNAME=shared-user
  - VANGOGH_SERVE_SHARED-PASSWORD=shared-password
```

After setting those values - you'll need to restart `vangogh` service with `docker compose restart` (you'll need to be in the same directory `compose.yml` for `vangogh` is).

You can see up to date specification of what endpoints require authentication, as well as role requirements [here](https://github.com/arelate/vangogh/blob/main/rest/routing.go).

## Taking care of your data

Syncing data keeps it in sync with GOG.com. As part of this process some data is left on the system and `vangogh` provides few ways to clean up and vet the data:

- `cleanup` will take care of older downloads versions. `cleanup -all -test` will enumerate all installers that are not linked to most current data. `cleanup -all` (without `-test`) will move that stale data to `recycle_bin` under your state directory. `cleanup -all -delete` will delete stale data (instead of moving to `recycle_bin`).
- `vet` will take care of various data problems, such as local-only images, files in the `recycle_bin`, old logs, etc. `vet -all` will run series of tests of data and print out recommendations. `vet -all -fix` will also attempt to repair the problem.
- `validate` will test installers you've downloaded using validation files provided by GOG.com. `validate -all-not-valid` will do that for all installers that are not in valid state. Please note that GOG.com is missing validation files for some installers and this will not be considered a critical error - you can use `vet` to generate missing checksums (by computing them locally - won't validate data validity, but will help maintain data consistency over time). 

## Sharing games

`vangogh` assumes you follow GOG.com [games sharing guidelines](https://support.gog.com/hc/en-us/articles/212184489-Can-I-share-games-with-others-?product=gog). Just like GOG.com, we trust you that this will not be abused.
