A simple web UI that will serve a random clip from a set of videos, and update the clip every day at a certain time.

To see a live demo of this, visit [clips.btschwartz.com](https://clips.btschwartz.com).

## Usage

### 1. Set up a media directory

For this to work, you need to have a directory with one or more video files. As of now, only `.mp4` is supported.

While the media directory is searched recursively, filenames must be unique. I recommend using a naming scheme that includes the season and episode number (e.g. `105.mp4` for season 1 episode 5), so that users can guess which episode the clip is from.

For example, here is a possible media directory structure:

```
$ tree /Volumes/SAUL
├── s1
│   ├── 101.mp4
│   ├── 102.mp4
│   ├── 103.mp4
│   ├── 104.mp4
│   ├── 105.mp4
│   ├── 106.mp4
│   ├── 107.mp4
│   ├── 108.mp4
│   ├── 109.mp4
│   └── 110.mp4
├── s2
│   ├── 201.mp4
│   ├── 202.mp4
│   ├── 203.mp4
│   ├── 204.mp4
...
```

### 2. Install `ffmpeg`

This project uses `ffmpeg` to generate clips. You can figure out how to install it. Just ensure you have the binaries `ffmpeg` and `ffprobe` in your `PATH`.

### 3. Set up your environment

However you want to do it, you need to have the following environment variables set:

```bash
CLIPS_MEDIA_DIR=/Volumes/SAUL # or wherever your media directory is
CLIPS_DURATION=30s            # the duration of each clip (must be parsable by time.ParseDuration)
CLIPS_TIME_OF_DAY=12:00:00    # the time of day to regenerate a new clip (EST)
```

### 4. Run the server

Although I would run this with Docker Compose through a Cloudflare Tunnel (see [compose.yml](compose.yml)), I'm just going to describe how to run it naively.

Make sure you have the above environment variables set, then run the following:

```bash
$ pwd
/path/to/clips
$ make clean && make app
$ ./app --port 8000
```

Now, you can visit `http://localhost:8000` in your browser to see a random clip!
