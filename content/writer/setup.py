import os
import re

from setuptools import find_packages, setup


def read_version():
    regexp = re.compile(r"^__version__\W*=\W*'([\d.abrc]+)'")
    init_py = os.path.join(os.path.dirname(__file__),
                           'content_writer', '__init__.py')
    with open(init_py) as f:
        for line in f:
            match = regexp.match(line)
            if match is not None:
                return match.group(1)
        else:
            msg = 'Cannot find version in content_writer/__init__.py'
            raise RuntimeError(msg)


install_requires = ['aiohttp',
                    'aioamqp']


setup(name='content-writer',
      version=read_version(),
      description='CS 385 Final Project Instagram clone content service upload endpoint',
      platforms=['POSIX'],
      packages=find_packages(),
      install_requires=install_requires)
