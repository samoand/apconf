""" Test dict operations."""

import unittest
from apconf.util import config_diff

class TestConfigDiff(unittest.TestCase):
    """
    Test cases for the config_diff function.
    """

    def test_no_difference(self):
        """
        Test case where there are no differences between the dictionaries.
        """
        config_new = {'a': 1, 'b': 2}
        config_old = {'a': 1, 'b': 2}
        expected_changed = {}
        expected_added = {}
        expected_removed = {}
        config_diff_result = config_diff(config_new, config_old)
        self.assertEqual(config_diff_result.diff_changed, expected_changed)
        self.assertEqual(config_diff_result.diff_added, expected_added)
        self.assertEqual(config_diff_result.diff_removed, expected_removed)

    def test_simple_difference(self):
        """
        Test case where there are simple key-value differences between the 
        dictionaries.
        """
        config_new = {'a': 1, 'b': 2}
        config_old = {'a': 1, 'b': 3}
        expected_changed = {'b': 2}
        expected_added = {}
        expected_removed = {}
        config_diff_result = config_diff(config_new, config_old)
        self.assertEqual(config_diff_result.diff_changed, expected_changed)
        self.assertEqual(config_diff_result.diff_added, expected_added)
        self.assertEqual(config_diff_result.diff_removed, expected_removed)

    def test_nested_difference(self):
        """
        Test case where there are nested dictionary differences.
        """
        config_new = {'a': 1, 'b': {'c': 2, 'd': 3}}
        config_old = {'a': 1, 'b': {'c': 2, 'd': 4}}
        expected_changed = {'b': {'d': 3}}
        expected_added = {}
        expected_removed = {}
        config_diff_result = config_diff(config_new, config_old)
        self.assertEqual(config_diff_result.diff_changed, expected_changed)
        self.assertEqual(config_diff_result.diff_added, expected_added)
        self.assertEqual(config_diff_result.diff_removed, expected_removed)

    def test_extra_keys_in_config_new(self):
        """
        Test case where config_new has extra keys not present in config_old.
        """
        config_new = {'a': 1, 'b': 2, 'c': 3}
        config_old = {'a': 1, 'b': 2}
        expected_changed = {}
        expected_added = {'c': 3}
        expected_removed = {}
        config_diff_result = config_diff(config_new, config_old)
        self.assertEqual(config_diff_result.diff_changed, expected_changed)
        self.assertEqual(config_diff_result.diff_added, expected_added)
        self.assertEqual(config_diff_result.diff_removed, expected_removed)

    def test_extra_keys_in_config_old(self):
        """
        Test case where config_old has extra keys not present in config_new.
        """
        config_new = {'a': 1}
        config_old = {'a': 1, 'b': 2}
        expected_changed = {}
        expected_added = {}
        expected_removed = {'b': 2}
        config_diff_result = config_diff(config_new, config_old)
        self.assertEqual(config_diff_result.diff_changed, expected_changed)
        self.assertEqual(config_diff_result.diff_added, expected_added)
        self.assertEqual(config_diff_result.diff_removed, expected_removed)

    def test_complex_nested_difference(self):
        """
        Test case with complex nested dictionary differences.
        """
        config_new = {'a': 1, 'b': {'c': 2, 'd': {'e': 3, 'f': 4}}}
        config_old = {'a': 1, 'b': {'c': 2, 'd': {'e': 3, 'f': 5}, 'g': 6}}
        expected_changed = {'b': {'d': {'f': 4}}}
        expected_added = {}
        expected_removed = {'b': {'g': 6}}
        config_diff_result = config_diff(config_new, config_old)
        self.assertEqual(config_diff_result.diff_changed, expected_changed)
        self.assertEqual(config_diff_result.diff_added, expected_added)
        self.assertEqual(config_diff_result.diff_removed, expected_removed)

if __name__ == '__main__':
    unittest.main()
