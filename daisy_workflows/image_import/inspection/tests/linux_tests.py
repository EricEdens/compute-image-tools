#!/usr/bin/env python3
# Copyright 2020 Google Inc. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

import os
import unittest

import inspection
import model
import system.filesystems
import yaml


class TestYAMLEncoded(unittest.TestCase):
  """Runs inspections using real-world metadata files.

  Test cases are located in `test-data`.
  """

  def test_all_systems(self):
    for fname in os.listdir('test-data'):
      fpath = os.path.join('test-data', fname)
      with self.subTest(msg=fname):
        self.run_yaml_encoded_test(fpath)

  def run_yaml_encoded_test(self, fpath):
    with open(fpath) as stream:
      loaded_yaml = yaml.safe_load(stream)
      assert 'files' in loaded_yaml
      assert 'expected' in loaded_yaml
      distro = model.distro_for(loaded_yaml['expected']['distro'])
      assert distro is not None
      fs = system.filesystems.DictBackedFilesystem(loaded_yaml['files'])
      expected = model.OperatingSystem(
        distro,
        model.Version(loaded_yaml['expected']['major'],
                      loaded_yaml['expected']['minor']))

    inspector = inspection._linux_inspector(fs)
    actual = inspector.inspect()
    self.assertIsNotNone(actual)
    self.assertEqual(expected.distro, actual.distro,
                     'expected=%s, actual=%s' % (expected, actual))
    self.assertEqual(expected.version, actual.version,
                     'expected=%s, actual=%s' % (expected, actual))


if __name__ == '__main__':
  unittest.main()
