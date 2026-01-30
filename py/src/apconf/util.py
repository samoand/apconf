"""Utility functions."""
import os
from typing import Callable, Optional

class ConfigDiffResult:
    """
    Structure reflecting the result of config_diff operation.

    Attributes:
        diff_changed: Keys with changed values between old and new config
        diff_added: Keys only in new config (added)
        diff_removed: Keys only in old config (removed)

    Supports intuitive membership testing:
        >>> "logging_config" in diff_result  # Checks all three dicts
        True

    Supports boolean testing:
        >>> if diff_result:  # True if any changes exist
        ...     apply_changes()
    """

    def __init__(
            self,
            diff_changed: dict[str, object],
            diff_added: dict[str, object],
            diff_removed: dict[str, object]
    ) -> None:
        self.diff_changed = diff_changed
        self.diff_added = diff_added
        self.diff_removed = diff_removed

    def __contains__(self, key: str) -> bool:
        """
        Check if key is in any of the diff dicts (changed, added, or removed).

        This allows intuitive usage like:
            if "logging_config" in config_diff:
                apply_logging(config)
        """
        return (key in self.diff_changed or
                key in self.diff_added or
                key in self.diff_removed)

    def __bool__(self) -> bool:
        """
        Return True if there are any changes (changed, added, or removed).

        This allows intuitive usage like:
            if config_diff:
                deploy_changes()
        """
        return bool(self.diff_changed or self.diff_added or self.diff_removed)

    def __repr__(self) -> str:
        """Return a readable representation."""
        return (f"ConfigDiffResult(changed={list(self.diff_changed.keys())}, "
                f"added={list(self.diff_added.keys())}, "
                f"removed={list(self.diff_removed.keys())})")

def find_parent_dir(
        start_path,
        matchers:list[Callable[[str], bool]]) -> Optional[str]:
    """
    Find the root of the project by searching for the parent directory which matches the "matcher".
    
    Args:
    - start_path: A string representing the starting directory path.
    - matcher: a function which reads a dir and decides whether it's the right parent
    
    Returns:
    - The path to the root of the git project, or None if the proper directory is not found.
    """
    if start_path is None:
        start_path = os.path.realpath(__file__)
    current_path = os.path.dirname(start_path)
    while current_path != os.path.dirname(current_path):
        if any(matcher(current_path) for matcher in matchers):
            return current_path
        current_path = os.path.dirname(current_path)
    return None

def find_parent_dir_or_error(
        start_path,
        matchers:list[Callable[[str], bool]]) -> str:
    """ Find a parent which meets a criteria."""
    result = find_parent_dir(start_path,
                             matchers)
    if result is None:
        raise ValueError('Could not find project root.')
    return result

def find_git_root(start_path) -> Optional[str]:
    """
    Find the root of the git project by searching for the .git directory.
    
    Args:
    - start_path: A string representing the starting directory path. 
                  If None, starts from the directory containing the current script.
    
    Returns:
    - The path to the root of the git project, or None if the .git directory is not found.
    """
    return find_parent_dir(
        start_path,
        [lambda path: os.path.isdir(os.path.join(path, '.git'))])

def find_git_root_or_error(start_path) -> str:
    """Find the root of the git project, report error if not found."""
    return find_parent_dir_or_error(
        start_path,
        [lambda path: os.path.isdir(os.path.join(path, '.git'))])

def config_diff(
        config_new: dict[str, object],
        config_old: dict[str, object]) -> ConfigDiffResult:
    """
    Return three dictionaries: one with the elements from config_new that are
    different from config_old, one with elements that are in config_new but not
    in config_old (added), and one with elements that are in config_old but not
    in config_new (removed).

    Args:
        config_new (dict): The first dictionary to compare.
        config_old (dict): The second dictionary to compare.

    Returns:
        ConfigDiffResults: Three dictionaries:
            - changed: Elements that differ between config_new and config_old.
            - added: Elements that are in config_new but not in config_old.
            - removed: Elements that are in config_old but not in config_new.
    """
    changed = {}
    added = {}
    removed = {}

    all_keys = set(config_new.keys()).union(config_old.keys())

    for key in all_keys:
        if key in config_new and key in config_old:
            if isinstance(config_new[key], dict
                          ) and isinstance(config_old[key], dict):
                nested_diff = config_diff(
                    config_new[key], config_old[key])
                if nested_diff.diff_changed:  # Only add non-empty differences
                    changed[key] = nested_diff.diff_changed
                if nested_diff.diff_added:  # Only add non-empty additions
                    added[key] = nested_diff.diff_added
                if nested_diff.diff_removed:  # Only add non-empty removals
                    removed[key] = nested_diff.diff_removed
            elif config_new[key] != config_old[key]:
                changed[key] = config_new[key]
        elif key in config_new:
            added[key] = config_new[key]
        elif key in config_old:
            removed[key] = config_old[key]

    return ConfigDiffResult(changed, added, removed)
