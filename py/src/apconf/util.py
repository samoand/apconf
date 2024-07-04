"""Utility functions."""
import os
from typing import Callable, Optional

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

def config_diff(config_new, config_old, preserve_missing=False):
    """
    Return a dictionary with the elements from config_new that are
      different from config_old.
    Optionally preserves elements from config_old if they are
      missing in config_new.

    Args:
        config_new (dict): The first dictionary to compare.
        config_old (dict): The second dictionary to compare.
        preserve_old (bool): Whether to preserve elements
          from config_old that are missing in config_new.

    Returns:
        dict: A dictionary containing the elements from
                config_new that differ from config_old.
              If preserve_old is True, also includes elements
                from config_old that are missing in config_new.

    Example usage:
    config_one = {
        "a": 1,
        "b": {
            "c": 2,
            "d": 3
        },
        "e": 4
    }
    
    config_two = {
        "a": 1,
        "b": {
            "c": 2,
            "d": 4
        },
        "f": 5
    }
    
    diff = config_diff(config_one, config_two)
    print(diff)  # Output: {'b': {'d': 3}, 'e': 4}
    
    """
    diff_dict = {}

    all_keys = set(config_new.keys()).union(config_old.keys())

    for key in all_keys:
        if key in config_new and key in config_old:
            if isinstance(config_new[key], dict) and isinstance(
                    config_old[key], dict):
                nested_diff = config_diff(
                    config_new[key], config_old[key], preserve_missing)
                if nested_diff:  # Only add non-empty differences
                    diff_dict[key] = nested_diff
            elif config_new[key] != config_old[key]:
                diff_dict[key] = config_new[key]
        elif key in config_new:
            diff_dict[key] = config_new[key]
        elif key in config_old and preserve_missing:
            diff_dict[key] = config_old[key]
    return diff_dict
