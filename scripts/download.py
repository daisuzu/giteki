#!/usr/bin/env python
# -*- coding: utf-8 -*-
""" download.py

    A script to download the xls of Japan certified radio equipment list.

    :copyright: Copyright (c) 2015 daisuzu
    :license: MIT
"""

import argparse
import csv
import datetime
import os
import re

from logging import (
    getLogger,
    Formatter,
    StreamHandler,
    INFO,
)
logger = getLogger(__name__)

from bs4 import BeautifulSoup
import mechanicalsoup
import requests

BASE_URL = 'http://www.tele.soumu.go.jp'
TECH_URL = 'http://www.tele.soumu.go.jp/j/sys/equ/tech/tech/index.htm'
DOWNLOAD_TARGET = {
    '1.国内': True,
    '2.外国(相互承認)': False,
}


def get_page(url):
    browser = mechanicalsoup.Browser()

    page = browser.get(url)
    if not page:
        return str(page)
    return Page(page)


class Page(object):
    def __init__(self, page):
        self._page = page

    def _get_h1(self, text):
        for line in text.split('\n'):
            if line.startswith('var titlew'):
                return line

    def _h1str(self, h1):
        value = re.search(r'"(.*)"', h1).group(1)
        return BeautifulSoup(value).text

    def advertised_year(self):
        for script in self._page.soup.select('script'):
            h1 = self._get_h1(script.text)
            if h1:
                return self._h1str(h1)
        return ''

    def get_tables(self):
        for div in self._page.soup.select('.mbtab0'):
            h2 = div.find_previous('h2').string
            h3 = div.find_previous('h3').string
            yield Table(h2, h3, div.find('table'))


class Table(object):
    def __init__(self, h2, h3, tables):
        self._h2 = h2
        self._h3 = h3
        self._tables = tables

    def h2(self):
        return self._h2

    def h3(self):
        return self._h3

    def _get_association(self, a):
        if self._h2.startswith('1'):
            return a.find_previous('p').string
        else:
            return self._h2

    def get_links(self):
        for a in self._tables.select('a'):
            link = a.attrs.get('href')
            if not link.endswith(('xls', 'htm')):
                continue
            association = self._get_association(a)

            yield Link(link, association)


class Link(object):
    def __init__(self, link, association):
        self._link = link
        self._association = association

    def link(self):
        return self._link

    def filename(self):
        return os.path.basename(self._link)

    def association(self):
        return self._association


def download_xls(url, dst):
    r = requests.get(url, stream=True)
    if not r:
        return False

    with open(dst, 'wb') as f:
        for chunk in r.iter_content(chunk_size=1024):
            if chunk:
                f.write(chunk)
                f.flush()
    return True


def download(url, path, download_all=False, update=False):
    logger.info('start download')
    result = []

    # get page
    page = get_page(url)
    if not isinstance(page, Page):
        msg = 'failed to download: {0}'.format(page)
        logger.error(msg)
        return

    advertised_year = page.advertised_year()
    logger.info('title: "{0}"'.format(advertised_year))

    # parse table
    for table in page.get_tables():
        logger.info('table: "{0} {1}"'.format(table.h2(), table.h3()))
        if (not download_all) and (not DOWNLOAD_TARGET.get(table.h2(), False)):
            logger.info('skip')
            continue

        # get link
        for link in table.get_links():
            filename = link.filename()
            association = link.association()
            logger.info(
                'link: "{0}", association: "{1}"'.format(filename, association)
            )
            if filename.endswith('htm'):
                if not download_all:
                    logger.info('skip')
                else:
                    result += download(
                        BASE_URL + link.link(),
                        path,
                        download_all=False,
                        update=update
                    )

                continue

            # download xls
            dst = os.path.join(path, filename)
            if os.path.exists(dst) and (not update):
                logger.info('already exists')
                continue
            if not download_xls(BASE_URL + link.link(), dst):
                logger.error('failed to download xls')
                continue

            logger.info('save to "{0}/{1}"'.format(path, filename))
            result.append([
                filename,
                association,
                str(datetime.datetime.now()),
            ])

    logger.info('done download')
    return result


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
            'This script downloads the xls of '
            'Japan certified radio equipment list.'
        )
    )

    parser.add_argument(
        '-d', '--dst',
        action='store',
        nargs='?',
        type=str,
        help='destination directory of the file to download',
        metavar='PATH',
    )
    parser.add_argument(
        '-A', '--all',
        action='store_true',
        default=False,
        help='download all of the files, including the past',
        dest='download_all',
    )
    parser.add_argument(
        '-U', '--update',
        action='store_true',
        default=False,
        help='overwrite if the file to download is exists',
    )

    args = parser.parse_args()
    logger.debug(args)

    # prepare
    path = args.dst
    if not path:
        path = os.path.join(os.getcwd(), 'downloads')

    if not os.path.exists(path):
        logger.debug('mkdir {0}'.format(path))
        os.mkdir(path)

    # download
    result = download(
        TECH_URL,
        path,
        download_all=args.download_all,
        update=args.update,
    )
    if not result:
        return

    # generate summary csv
    summary_file = os.path.join(
        path,
        'download_summary_{0}.csv'.format(
            datetime.datetime.now().strftime('%Y%m%d')
        )
    )
    with open(summary_file, 'w') as f:
        writer = csv.writer(f)
        for r in result:
            writer.writerow(r)


if __name__ == '__main__':
    main()
