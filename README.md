# CI sitemap-checker

## What is it about?
This utility will help you check if all pages from your sitemap are accessible

## Example usage:
``go run ./... -sitemap https://example.com/sitemap.xml``

## Example success output:
All pages are success checked.
Terminal status code - 0.

## Example when having error:
``Page: http://localhost/about. Status code: 500``<br>
We have page that has error and given status code.
Terminal status code - 1.