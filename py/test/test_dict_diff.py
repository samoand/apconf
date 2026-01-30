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

    def test_contains_changed(self):
        """Test __contains__ for changed keys."""
        config_new = {'a': 1, 'b': 2}
        config_old = {'a': 1, 'b': 3}
        result = config_diff(config_new, config_old)
        self.assertIn('b', result)
        self.assertNotIn('a', result)
        self.assertNotIn('c', result)

    def test_contains_added(self):
        """Test __contains__ for added keys."""
        config_new = {'a': 1, 'b': 2}
        config_old = {'a': 1}
        result = config_diff(config_new, config_old)
        self.assertIn('b', result)
        self.assertNotIn('a', result)

    def test_contains_removed(self):
        """Test __contains__ for removed keys."""
        config_new = {'a': 1}
        config_old = {'a': 1, 'b': 2}
        result = config_diff(config_new, config_old)
        self.assertIn('b', result)
        self.assertNotIn('a', result)

    def test_bool_true_when_changes(self):
        """Test __bool__ returns True when there are changes."""
        config_new = {'a': 1, 'b': 2}
        config_old = {'a': 1}
        result = config_diff(config_new, config_old)
        self.assertTrue(result)
        self.assertTrue(bool(result))

    def test_bool_false_when_no_changes(self):
        """Test __bool__ returns False when there are no changes."""
        config_new = {'a': 1, 'b': 2}
        config_old = {'a': 1, 'b': 2}
        result = config_diff(config_new, config_old)
        self.assertFalse(result)
        self.assertFalse(bool(result))

    def test_repr(self):
        """Test __repr__ for readable output."""
        config_new = {'a': 1, 'b': 2, 'c': 3}
        config_old = {'a': 1, 'd': 4}
        result = config_diff(config_new, config_old)
        repr_str = repr(result)
        self.assertIn('ConfigDiffResult', repr_str)
        self.assertIn('changed=', repr_str)
        self.assertIn('added=', repr_str)
        self.assertIn('removed=', repr_str)


if __name__ == '__main__':
    unittest.main()
