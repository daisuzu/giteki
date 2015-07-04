#!/usr/bin/env python
# -*- coding: utf-8 -*-
""" read.py

    A script to read the downloaded xls.

    :copyright: Copyright (c) 2015 daisuzu
    :license: MIT
"""

import argparse
import datetime
import json

from logging import (
    getLogger,
    Formatter,
    StreamHandler,
    INFO,
)
logger = getLogger(__name__)

import xlrd


def read(src):
    result = []

    book = xlrd.open_workbook(src)
    sheet = book.sheet_by_index(0)
    for rx in range(sheet.nrows):
        tmp = []
        cells = sheet.row(rx)
        for cell in cells:
            if cell.ctype == 2:
                tmp.append(str(cell.value).replace('.0', ''))
            elif cell.ctype == 3:
                d = xlrd.xldate_as_tuple(cell.value, 0)
                if not d:
                    continue
                tmp.append(
                    datetime.datetime(d[0], d[1], d[2]).strftime('%Y-%m-%d')
                )
            else:
                tmp.append(cell.value)
        result.append(tmp)

    print(json.dumps({'Result': result}, ensure_ascii=False))


def main():
    # set logger
    formatter = Formatter(
        fmt='%(asctime)s %(levelname)s: %(message)s',
    )

    handler = StreamHandler()
    handler.setLevel(INFO)
    handler.formatter = formatter

    logger.setLevel(INFO)
    logger.addHandler(handler)

    # parse arguments
    parser = argparse.ArgumentParser(
        description=(
            'This script read the downloaded xls.'
        )
    )

    parser.add_argument(
        '-s', '--src',
        action='store',
        nargs=1,
        type=str,
        required=True,
        help='source path of the xls to load',
        metavar='PATH',
    )

    args = parser.parse_args()
    logger.debug(args)

    # read
    read(args.src[0])


if __name__ == '__main__':
    main()
