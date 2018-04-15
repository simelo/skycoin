"""PySkycoin - Skycoin client library for Python.
"""

# Always prefer setuptools over distutils
from setuptools import setup, find_packages
# To use a consistent encoding
from codecs import open
from os import path

here = path.abspath(path.dirname(__file__))

# Get the long description from the README file
with open(path.join(here, 'README.md'), encoding='utf-8') as f:
    long_description = f.read()

setup(
    name='pyskycoin',
    version='0.23.0',
    description='PySkycoin - Skycoin client library for Python.',
    long_description=long_description,
    long_description_content_type='text/markdown',
    url='https://github.com/skycoin/skycoin/tree/master/lib/py',
    author='The Skycoin Team',
    author_email='skycoin.dev@gmail.com',
    maintainer='Olemis Lang',
    maintainer_email='olemis@simelo.tech',
    classifiers=[
        'Development Status :: 4 - Beta',
        'Intended Audience :: Developers',
        'Intended Audience :: Financial and Insurance Industry',
        'Operating System :: Microsoft :: Windows',
        'Operating System :: MacOS',
        'Operating System :: POSIX :: Linux',
        'Topic :: Security :: Cryptography',
        'Topic :: System :: Distributed Computing',
        'License :: OSI Approved :: MIT License',
        'Programming Language :: Python :: 3',
        'Programming Language :: Python :: 3.4',
        'Programming Language :: Python :: 3.5',
        'Programming Language :: Python :: 3.6',
    ],
    keywords='cryptocurrency fintech blockchain cryptography decentralizatio netneutrality',
    packages=find_packages(exclude=['contrib', 'docs', 'tests']),
    extras_require={
        'dev': ['pybindgen'],
        'test': ['coverage', 'pytest'],
    },
    project_urls={
        'Bug Reports': 'https://github.com/skycoin/skycoin/issues',
        'Source': 'https://github.com/skycoin/skycoin',
    },
)
