#!/usr/bin/env python

import collections
import json
import os
import re
import requests
import sys

# Typical naming convention: "Local language title (Original title) date"
MediaName = collections.namedtuple("MediaName", ["local", "original", "year"])
naming = re.compile(
    r"(?P<local>[^\(]+?)( \((?P<original>.*)\))?( (?P<year>\d{4}))?($|\.)"
)

# MovieDB access token
MOVIEDB_TOKEN = os.environ.get("MOVIEDB_TOKEN")
MOVIEDB_HEADERS = {"accept": "application/json", "Authorization": f"Bearer {MOVIEDB_TOKEN}"}
MOVIEDB_GENRES = "moviedb_genres.json"
MOVIEDB_DATA = "moviedb_data.json"

# Results
MAPPING_NAME = "mapping.json"

def parseName(moviePath):
    media = naming.match(os.path.basename(moviePath))
    if not media:
        return None
    return MediaName(**media.groupdict())


def requestMovieDBGenres(apiKey=None):
    return requests.get(
        "https://api.themoviedb.org/3/genre/movie/list?language=ru",
        headers=MOVIEDB_HEADERS,
    ).json()


def requestMovieDB(movieDef, apiKey=None):
    searchParams = {"query": movieDef.original if movieDef.original else movieDef.local}
    if movieDef.year:
        searchParams["year"] = movieDef.year
    return requests.get(
        "https://api.themoviedb.org/3/search/movie",
        searchParams,
        headers=MOVIEDB_HEADERS,
    ).json()

def readJSON(path):
    if os.path.isfile(path):
        with open(path, "r") as f:
            return json.loads(f.read())
    return {}

def writeJSON(path, data):
    with open(path, "w") as f:
        f.write(json.dumps(data, ensure_ascii=False))

def main():
    genreData = readJSON(MOVIEDB_GENRES)
    movieData = readJSON(MOVIEDB_DATA)

    if not genreData:
        print("Loading genre information")
        genreData = requestMovieDBGenres(MOVIEDB_TOKEN)
        writeJSON(MOVIEDB_GENRES, genreData)

    if not movieData:
        print("Loading moviedb information")
        for moviePath in os.listdir(sys.argv[1]):
            movieDef = parseName(moviePath)
            movieData[moviePath] = requestMovieDB(parseName(moviePath), MOVIEDB_TOKEN)
        writeJSON(MOVIEDB_DATA, movieData)

    print("Resolving genre names")
    genreById = {}
    for genre in genreData["genres"]:
        genreById[genre["id"]] = genre["name"]
    
    print("Building mapping")
    mapping = collections.defaultdict(dict)
    for moviePath, data in movieData.items():
        if not data["results"]:
            mapping[moviePath] = { "genres": ["other_unknown"] }
            continue
        genres = [
            genreById[id] for r in data["results"] for id in r.get("genre_ids", [])
        ]
        mapping[moviePath] = {
	        "genres": sorted(set(genres))
	    }
    writeJSON(MAPPING_NAME, mapping)


if __name__ == "__main__":
    main()
