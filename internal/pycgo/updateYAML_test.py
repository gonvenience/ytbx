# -*- coding: utf-8 -*-

# Copyright © 2018 Matthias Diester
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

import os
import unittest
import tempfile
import shutil
from ruamel.yaml import YAML
from updateYAML import getEntryIdxByIdentifier
from updateYAML import updateYAML


class TestStringMethods(unittest.TestCase):
    assets = os.path.dirname(os.path.realpath(__file__)) + "/../../assets"

    def test_empty_list_in_getEntryIdxByIdentifier(self):
        self.assertEqual(getEntryIdxByIdentifier([], "name", "foobar"), -1)

    def test_find_entry_using_test_empty_list_in_getEntryIdxByIdentifier(self):
        self.assertEqual(
            getEntryIdxByIdentifier([{
                "name": "zero"
            }, {
                "name": "one"
            }], "name", "one"), 1)

    def test_unable_to_find_entry_using_getEntryIdxByIdentifier(self):
        self.assertEqual(
            getEntryIdxByIdentifier([{
                "name": "zero"
            }, {
                "name": "one"
            }], "name", "two"), -1)

    def test_updateYAML_set_existing_map_path(self):
        with tempfile.NamedTemporaryFile() as temp:
            shutil.copy(self.assets + "/testbed/example.yml", temp.name)
            updateYAML(temp.name, "/yaml/structure/somekey", "Foo")
            with open(temp.name, "r") as file:
                data = YAML().load(file.read())
                self.assertEqual(data["yaml"]["structure"]["somekey"], "Foo")

    def test_updateYAML_set_existing_list_path(self):
        with tempfile.NamedTemporaryFile() as temp:
            shutil.copy(self.assets + "/testbed/example.yml", temp.name)
            updateYAML(temp.name, "/list/name=one/somekey", "Foo")
            with open(temp.name, "r") as file:
                data = YAML().load(file.read())
                self.assertEqual(data["list"][0]["somekey"], "Foo")

    def test_updateYAML_set_existing_simplelist_path(self):
        with tempfile.NamedTemporaryFile() as temp:
            shutil.copy(self.assets + "/testbed/example.yml", temp.name)
            updateYAML(temp.name, "/simpleList/1", "Foo")
            with open(temp.name, "r") as file:
                data = YAML().load(file.read())
                self.assertEqual(data["simpleList"][1], "Foo")

    def test_updateYAML_set_new_map_path(self):
        with tempfile.NamedTemporaryFile() as temp:
            shutil.copy(self.assets + "/testbed/example.yml", temp.name)
            updateYAML(temp.name, "/yaml/structure/newkey", "newval")
            with open(temp.name, "r") as file:
                data = YAML().load(file.read())
                self.assertEqual(data["yaml"]["structure"]["newkey"], "newval")

    def test_updateYAML_set_new_list_path(self):
        with tempfile.NamedTemporaryFile() as temp:
            shutil.copy(self.assets + "/testbed/example.yml", temp.name)
            updateYAML(temp.name, "/list/name=two/somekey", "Foo")
            with open(temp.name, "r") as file:
                data = YAML().load(file.read())
                self.assertEqual(data["list"][1]["somekey"], "Foo")

    def test_updateYAML_set_new_simplelist_path(self):
        with tempfile.NamedTemporaryFile() as temp:
            shutil.copy(self.assets + "/testbed/example.yml", temp.name)
            updateYAML(temp.name, "/simpleList/-1", "Foo")
            with open(temp.name, "r") as file:
                data = YAML().load(file.read())
                self.assertEqual(data["simpleList"][2], "Foo")


if __name__ == '__main__':
    unittest.main()
