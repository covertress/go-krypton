#!/usr/bin/env python
import os
from distutils.core import setup, Extension
sources = [
    'src/python/core.c',
    'src/libkrash/io.c',
    'src/libkrash/internal.c',
    'src/libkrash/sha3.c']
if os.name == 'nt':
    sources += [
        'src/libkrash/util_win32.c',
        'src/libkrash/io_win32.c',
        'src/libkrash/mmap_win32.c',
    ]
else:
    sources += [
        'src/libkrash/io_posix.c'
    ]
depends = [
    'src/libkrash/krash.h',
    'src/libkrash/compiler.h',
    'src/libkrash/data_sizes.h',
    'src/libkrash/endian.h',
    'src/libkrash/krash.h',
    'src/libkrash/io.h',
    'src/libkrash/fnv.h',
    'src/libkrash/internal.h',
    'src/libkrash/sha3.h',
    'src/libkrash/util.h',
]
pykrash = Extension('pykrash',
                     sources=sources,
                     depends=depends,
                     extra_compile_args=["-Isrc/", "-std=gnu99", "-Wall"])

setup(
    name='pykrash',
    author="Matthew Wampler-Doty",
    author_email="matthew.wampler.doty@gmail.com",
    license='GPL',
    version='0.1.23',
    url='https://github.com/krypton/krash',
    download_url='https://github.com/krypton/krash/tarball/v23',
    description=('Python wrappers for krash, the krypton proof of work'
                 'hashing function'),
    ext_modules=[pykrash],
)
