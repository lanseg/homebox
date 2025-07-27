#!/usr/bin/env python
# mount -o size=16G -t tmpfs none /mnt/tmpfs

import collections
import json
import os

MAPPING_FILE = "mapping.json"
MOVIES_DIR = "/mnt/storage/video/Кино"
BASE_DIR = "rams"


def main():
    mapping = {}
    if os.path.isfile(MAPPING_FILE):
        with open(MAPPING_FILE, "r") as f:
            mapping = json.loads(f.read())
    moviesByGenre = collections.defaultdict(list)
    for name, params in mapping.items():
        if not "genres" in params:
            continue
        for genre in params["genres"]:
            moviesByGenre[genre].append(name)

    if not os.path.isdir(BASE_DIR):
        os.mkdir(BASE_DIR)

    for genre, movies in moviesByGenre.items():
        genredir = os.path.join(BASE_DIR, genre)
        if not os.path.isdir(genredir):
            os.mkdir(genredir)
        for movie in sorted(movies):
            target = os.path.join(genredir, movie)
            if os.path.islink(target):
                print(f"Removing existing symlink {target}")
                os.remove(target)
            print(f"Creating symlink {movie} -> {target}")
            os.symlink(os.path.join(MOVIES_DIR, movie), target)


if __name__ == "__main__":
    main()
