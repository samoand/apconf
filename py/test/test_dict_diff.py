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
        expected = {}
        result = config_diff(config_new, config_old)
        self.assertEqual(result, expected)

    def test_simple_difference(self):
        """
        Test case where there are simple key-value differences between the 
        dictionaries.
        """
        config_new = {'a': 1, 'b': 2}
        config_old = {'a': 1, 'b': 3}
        expected = {'b': 2}
        result = config_diff(config_new, config_old)
        self.assertEqual(result, expected)

    def test_nested_difference(self):
        """
        Test case where there are nested dictionary differences.
        """
        config_new = {'a': 1, 'b': {'c': 2, 'd': 3}}
        config_old = {'a': 1, 'b': {'c': 2, 'd': 4}}
        expected = {'b': {'d': 3}}
        result = config_diff(config_new, config_old)
        self.assertEqual(result, expected)

    def test_extra_keys_in_config_new(self):
        """
        Test case where config_new has extra keys not present in config_old.
        """
        config_new = {'a': 1, 'b': 2, 'c': 3}
        config_old = {'a': 1, 'b': 2}
        expected = {'c': 3}
        result = config_diff(config_new, config_old)
        self.assertEqual(result, expected)

    def test_extra_keys_in_config_old_without_preserve(self):
        """
        Test case where config_old has extra keys not present in config_new 
        without preserving old elements.
        """
        config_new = {'a': 1}
        config_old = {'a': 1, 'b': 2}
        expected = {}
        result = config_diff(config_new, config_old)
        self.assertEqual(result, expected)

    def test_extra_keys_in_config_old_with_preserve(self):
        """
        Test case where config_old has extra keys not present in config_new 
        with preserving old elements.
        """
        config_new = {'a': 1}
        config_old = {'a': 1, 'b': 2}
        expected = {'b': 2}
        result = config_diff(config_new, config_old, preserve_missing=True)
        self.assertEqual(result, expected)

    def test_complex_nested_difference_without_preserve(self):
        """
        Test case with complex nested dictionary differences without 
        preserving old elements.
        """
        config_new = {'a': 1, 'b': {'c': 2, 'd': {'e': 3, 'f': 4}}}
        config_old = {'a': 1, 'b': {'c': 2, 'd': {'e': 3, 'f': 5}}}
        expected = {'b': {'d': {'f': 4}}}
        result = config_diff(config_new, config_old)
        self.assertEqual(result, expected)

    def test_complex_nested_difference_with_preserve(self):
        """
        Test case with complex nested dictionary differences with preserving 
        old elements.
        """
        config_new = {'a': 1, 'b': {'c': 2, 'd': {'e': 3, 'f': 4}}}
        config_old = {'a': 1, 'b': {'c': 2, 'd': {'e': 3, 'f': 5}, 'g': 6}}
        expected = {'b': {'d': {'f': 4}, 'g': 6}}
        result = config_diff(config_new, config_old, preserve_missing=True)
        self.assertEqual(result, expected)

if __name__ == '__main__':
    unittest.main()
