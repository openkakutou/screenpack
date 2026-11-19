"""Generate the default screenpack's original images and sound sources.

Run from the repository root:  python3 tools/default-assets/generate.py
Images are written to default/images/. Sounds are written as WAV to a
temporary folder given as the first argument, then encoded to OGG Vorbis
with ffmpeg (see tools/default-assets/README.md). Everything this script makes
is original and released under the repository's MIT license.
"""
import json
import math
import os
import struct
import sys
import wave

from PIL import Image, ImageDraw, ImageFont

ROOT = os.path.join(os.path.dirname(__file__), "..", "..", "default")
IMAGES = os.path.join(ROOT, "images")
BOLD = os.path.join(ROOT, "fonts", "NotoSans-Bold.ttf")

ACCENT = (255, 92, 56, 255)
INK = (240, 236, 228, 255)
NIGHT = (14, 15, 22, 255)


def background():
    w, h = 8, 512
    img = Image.new("RGBA", (w, h))
    top, bottom = (28, 26, 40), (8, 8, 12)
    for y in range(h):
        t = y / (h - 1)
        c = tuple(round(a + (b - a) * t) for a, b in zip(top, bottom))
        for x in range(w):
            img.putpixel((x, y), c + (255,))
    img.save(os.path.join(IMAGES, "background.png"))


def grid():
    s = 64
    img = Image.new("RGBA", (s, s), (0, 0, 0, 0))
    d = ImageDraw.Draw(img)
    d.line([(0, 0), (s - 1, 0)], fill=(255, 255, 255, 18))
    d.line([(0, 0), (0, s - 1)], fill=(255, 255, 255, 18))
    img.save(os.path.join(IMAGES, "grid.png"))


def logo():
    font = ImageFont.truetype(BOLD, 160)
    text = "OPENKAKUTOU"
    box = font.getbbox(text)
    w, h = box[2] - box[0] + 40, box[3] - box[1] + 60
    img = Image.new("RGBA", (w, h), (0, 0, 0, 0))
    d = ImageDraw.Draw(img)
    d.text((20 - box[0], 20 - box[1]), text, font=font, fill=INK, stroke_width=6, stroke_fill=NIGHT)
    d.rectangle([20, h - 22, w - 20, h - 12], fill=ACCENT)
    img.save(os.path.join(IMAGES, "logo.png"))


def sheet(name, frames, size, draw, animations=None):
    """A horizontal strip of frames, and its sprite-sheet description."""
    img = Image.new("RGBA", (size * len(frames), size), (0, 0, 0, 0))
    d = ImageDraw.Draw(img)
    desc = {"image": name + ".png", "frames": {}}
    for i, frame in enumerate(frames):
        draw(d, i * size, size, frame)
        desc["frames"][frame] = [i * size, 0, size, size]
    if animations:
        desc["animations"] = animations
    img.save(os.path.join(IMAGES, name + ".png"))
    with open(os.path.join(IMAGES, name + ".json"), "w") as f:
        json.dump(desc, f, indent=2)
        f.write("\n")


def cursor(d, x, s, frame):
    alpha = 255 if frame == "on" else 120
    w = 4
    for i in range(w):
        d.rectangle([x + i, i, x + s - 1 - i, s - 1 - i], outline=ACCENT[:3] + (alpha - i * 25,))


def arrow(d, x, s, frame):
    m = s // 4
    if frame == "left":
        pts = [(x + s - m, m), (x + m, s // 2), (x + s - m, s - m)]
    else:
        pts = [(x + m, m), (x + s - m, s // 2), (x + m, s - m)]
    d.polygon(pts, fill=INK)


BUTTON_LABELS = {"a": "A", "b": "B", "c": "C", "x": "X", "y": "Y", "z": "Z", "start": "S"}
DIRECTIONS = {"up": 270, "down": 90, "forward": 0, "back": 180}


def button(d, x, s, frame):
    d.ellipse([x + 2, 2, x + s - 3, s - 3], fill=(40, 40, 56, 255), outline=INK, width=2)
    c = s / 2
    if frame in DIRECTIONS:
        a = math.radians(DIRECTIONS[frame])
        r = s * 0.28
        tip = (x + c + r * math.cos(a), c + r * math.sin(a))
        left = (x + c + r * math.cos(a + 2.5), c + r * math.sin(a + 2.5))
        right = (x + c + r * math.cos(a - 2.5), c + r * math.sin(a - 2.5))
        d.polygon([tip, left, right], fill=ACCENT)
        return
    font = ImageFont.truetype(BOLD, int(s * 0.5))
    d.text((x + c, c), BUTTON_LABELS[frame], font=font, fill=INK, anchor="mm")


RATE = 44100


def tone(path, notes, volume=0.4):
    """notes: (frequency Hz, duration s) pairs, a short decaying square-ish blip each."""
    samples = []
    for freq, dur in notes:
        n = int(RATE * dur)
        for i in range(n):
            t = i / RATE
            env = math.exp(-6 * t / dur)
            v = math.sin(2 * math.pi * freq * t) + 0.3 * math.sin(2 * math.pi * 2 * freq * t)
            samples.append(volume * env * v / 1.3)
    write_wav(path, samples)


def music(path):
    """An 8-second loop: a slow pad over four chords."""
    chords = [(220.0, 261.63, 329.63), (174.61, 220.0, 261.63), (196.0, 246.94, 293.66), (164.81, 207.65, 246.94)]
    samples = []
    for chord in chords:
        n = RATE * 2
        for i in range(n):
            t = i / RATE
            env = min(1, t / 0.4) * min(1, (2 - t) / 0.4)
            v = sum(math.sin(2 * math.pi * f * t) + 0.2 * math.sin(2 * math.pi * f * 2.01 * t) for f in chord)
            samples.append(0.12 * env * v)
    write_wav(path, samples)


def write_wav(path, samples):
    with wave.open(path, "wb") as w:
        w.setnchannels(1)
        w.setsampwidth(2)
        w.setframerate(RATE)
        w.writeframes(b"".join(struct.pack("<h", int(max(-1, min(1, s)) * 32767)) for s in samples))


def main():
    os.makedirs(IMAGES, exist_ok=True)
    background()
    grid()
    logo()
    sheet("cursor", ["on", "off"], 32, cursor,
          {"pulse": {"frames": [{"frame": "on", "ms": 500}, {"frame": "off", "ms": 500}]}})
    sheet("arrows", ["left", "right"], 32, arrow)
    sheet("buttons", list(DIRECTIONS) + list(BUTTON_LABELS), 32, button)
    if len(sys.argv) > 1:
        out = sys.argv[1]
        tone(os.path.join(out, "move.wav"), [(880, 0.06)])
        tone(os.path.join(out, "confirm.wav"), [(660, 0.07), (990, 0.12)])
        tone(os.path.join(out, "cancel.wav"), [(660, 0.07), (440, 0.12)])
        tone(os.path.join(out, "error.wav"), [(180, 0.18)], volume=0.5)
        tone(os.path.join(out, "tick.wav"), [(1320, 0.05)], volume=0.3)
        tone(os.path.join(out, "gameover.wav"), [(392, 0.25), (330, 0.25), (262, 0.6)])
        music(os.path.join(out, "theme.wav"))


main()
