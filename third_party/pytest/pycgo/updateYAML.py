# -*- coding: utf-8 -*-

# Copyright © 2018 The Homeport Team
#
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in
# all copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
# THE SOFTWARE.

import re
import sys
import json
import urllib.parse
from ruamel.yaml import YAML


def getEntryIdxByIdentifier(list, identifier, name):
    for idx, item in enumerate(list):
        if item[identifier] == name:
            return idx

    return -1


def getEmptyStructure(pathElementTypes, idx):
    if idx >= len(pathElementTypes): return {}
    if pathElementTypes[idx] == "map": return {}
    if pathElementTypes[idx] == "list": return []
    return "undef"


def IsInteger(str):
    try:
        int(str)
        return True

    except ValueError:
        return False


def updateYAML(location, pathString, newValue):
    regex = r"([a-zA-Z0-9_-]+)=(.+)"

    yaml = YAML()
    yaml.default_flow_style = False
    yaml.explicit_start = True

    global data
    with open(location, "r") as file:
        data = yaml.load(file.read())

    pathElements = pathString.split("/")

    pathElementTypes = ["undef"]
    for pathIdx, pathElement in enumerate(pathElements):
        if pathIdx == 0: continue

        pathElement = urllib.parse.unquote(pathElement)
        if re.search(regex, pathElement):
            pathElementTypes.append("list")

        elif IsInteger(pathElement):
            pathElementTypes.append("list")

        else:
            pathElementTypes.append("map")

    pointer = data
    for pathIdx, pathElement in enumerate(pathElements):
        if pathIdx == 0: continue

        pathElement = urllib.parse.unquote(pathElement)
        lastOne = pathIdx == len(pathElements) - 1

        match = re.search(regex, pathElement)
        if match:
            identifier, name = match.group(1), match.group(2)
            idx = getEntryIdxByIdentifier(pointer, identifier, name)

            if idx < 0:
                pointer.append(yaml.load("%s: %s" % (identifier, name)))
                idx = len(pointer) - 1

            if not lastOne:
                pointer = pointer[idx]

            else:
                pointer[idx] = newValue

        elif IsInteger(pathElement):
            idx = int(pathElement)

            if idx < 0:
                pointer.append("undef")
                idx = len(pointer) - 1

            if not lastOne:
                pointer = pointer[idx]

            else:
                pointer[idx] = newValue

        else:
            if not pathElement in pointer:
                pointer.insert(
                    len(pointer), pathElement,
                    getEmptyStructure(pathElementTypes, pathIdx + 1))

            if not lastOne:
                pointer = pointer[pathElement]

            else:
                pointer[pathElement] = newValue

    with open(location, 'w') as f:
        yaml.dump(data, f)


if __name__ == '__main__':
    location = sys.argv[0]
    pathString = sys.argv[1]
    newValue = YAML().load(sys.argv[2])
    updateYAML(location, pathString, newValue)
