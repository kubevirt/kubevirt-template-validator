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
    tags = [ item['tag_name'] for item in data ]
    versions = [ Version.from_string(tag) for tag in tags]
    buckets = {}
    for ver in versions:
        key = (ver.major, ver.minor)
        if key not in buckets:
            buckets[key] = ver
        else:
            cur = buckets[key]
            buckets[key] = ver if ver > cur else cur
    return list(reversed(sorted(buckets.values())))


def _find_ver_idx(vers, target):
    for idx, ver in enumerate(vers):
        if ver == target:
            return idx
    return None


def _main():
    builtin = Version.from_string(sys.argv[1])
    vers = versions(json.load(sys.stdin))
    out = {
            'last': vers[0],
            'builtin': builtin,
            'previous': builtin
    }

    idx = _find_ver_idx(vers, builtin)
    if idx is not None or idx < len(builtin):
        out['previous'] = vers[idx+1]

    print('last=%s\nbuiltin=%s\nprevious=%s' % (
          out['last'], out['builtin'], out['previous']))


if __name__ == "__main__":
    if len(sys.argv) != 2:
        sys.exit(2)
    _main()
