# Datebook
> A small CLI utility for keeping a "datebook"

## What is this?

`datebook` is a CLI utility for keeping track of your day-to-day. Each `datebook` entry corresponds to single day which is a place to store notes, links, thoughts, scratch work and anything else you can put into a markdown file for the day or week.

Here's what a `datebook` entry _could_ look like:

~~~ markdown
# Sunday October 18

## todo

* finish the first version of `datebook`
* write some tests for `datebook`

## week

* work on more functional date parsing for `datebook`

~~~

## Installation

**TODO** add installation directions

## Usage

`datebook` uses a basic **NLP** based library to parse dates. For instance, the following commands could be used to pull up any assortment of `datebook` entries:

* `$ datebook tomorrow`
* `$ datebook today`
* `$ datebook june 1st`
* `$ datebook 4 months ago`
* `$ datebook tuesday`

## Template File

When a new entry is created in your `datebook` you can specify a starting template. By default, `datebook` expects this file to exist at `$HOME/.datebook.md`. To specify a custom template file, simply call `datebook` with a template argument.

~~~ bash
$ datebook -template /my_template_path.md
~~~

The template file allows some rudimentary string interpolation and currently supports the following variables.

* year
* month
* day
* weekday


### Sample

~~~ markdown

# %weekday% %month% %day%

## todo

## links

~~~

This template would create `datebook` entries like the following:

~~~ markdown
# Sunday October 18

## todo

## links

~~~

## Weeks 

Your `datebook` is a place to keep track of things. Sometimes it is helpful to keep a few things around through the course of a week. Things such as checklists, reminders for later in the week and anything else that might not correspond to just a single day can be passed between `datebook` entries.

By adding a `## week` block in any entry, each `datebook` entry that week will share said block.


