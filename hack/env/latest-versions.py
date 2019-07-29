#!/usr/bin/env python3

import collections
import sys
import json


_version = collections.namedtuple('_version', ['major', 'minor', 'micro'])


class Version(_version):

    @classmethod
    def from_string(cls, text):
        text = text.lstrip("v")
        text = text.split("-")[0]  # TODO: handle '-extra'?
        major, minor, micro = text.split('.')
        return cls(int(major), int(minor), int(micro))

    def __str__(self):
        return "v%d.%d.%d" % (self.major, self.minor, self.micro)


def versions(data):
    versions = [ Version.from_string(item['tag_name']) for item in data]
    buckets = {}
    for ver in versions:
        key = (ver.major, ver.minor)
        if key not in buckets:
            buckets[key] = ver
        else:
            cur = buckets[key]
            buckets[key] = ver if ver > cur else cur
    return list(reversed(sorted(buckets.values())))


if __name__ == "__main__":
    vers = versions(json.load(sys.stdin))
    print(str(vers[0]))
