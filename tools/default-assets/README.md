# Default screenpack assets

`generate.py` makes the default screenpack's original images and sounds (`default/ASSETS.md` records each one's license). It needs Python 3 with Pillow, and ffmpeg with libvorbis for the sounds.

```sh
python3 tools/default-assets/generate.py /tmp/wav     # images into default/images/, WAV into /tmp/wav
for f in /tmp/wav/*.wav; do
  ffmpeg -loglevel error -y -i "$f" -c:a libvorbis -q:a 3 "default/sounds/$(basename "$f" .wav).ogg"
done
```

The fonts are Noto Sans, copied unmodified; the script reads `default/fonts/NotoSans-Bold.ttf` for the logo and the button letters.
