Hey!

So here are the details about my Brick Recon tool.  The main aim is to help me find lego kits which provide parts I need to make my own projects.

Projects are what I call a single lego model I want to build. `model` was a too overloaded term in dev!
Projects can be read in from Stud.io, LDraw or WantedList TSVs.
Projects can have bricks marked as owned (inventory)
Projects can have Kits applied to them
Projects can have parts lists updated by re-importing
Projects can have BrickLink XML generated for buying remaining parts

Kits are official lego models
Kits are imported from a 3rd party api, for which you have to ask for access to nicely by email (BrickOwl)

When a Kit or Project is added, part data is scraped (from BrickLink, BrickOwl, and import source) to get images etc.
When a Kit or Project are added, I calculate the intersection of parts between all kits and projects, to see what kits are the best to add to a project

## Tech Stuff:

I wrote this to solve the main aim above, but also gave myself some extra constraints: no js, server side rendered, eventsouced.

- written in go
- datastore is a filesystem based eventsourcing lib (written in the project too)
- http "framework" is called Preen, and is written in project.  Provides controllers and rendering, and some data mapping.
- ui is just html files and go templates (which are nasty, but I haven't gotten round to writing a template engine yet...)
- frontend has little to no js.  What js there is degrades gracefully
- cli for doing lots of things, like diffing project parts, updating projects etc.


## Future Things

- Price scraping of parts (see what's expensive)
- Call APIs to find what kits contain a part, make kit suggestions
- Maybe image recognition and some way to process parts automagically? e.g. webcam, conveyor belt, image recogntion
-


## Screenshots:

br-kit.png:
- the view of an imported Lego Kit

br-project.png:
- a small project model I use to test with
- parts filterable by name, all/needed/owned (i.e. inventory >= quantity), colour
- kits which have parts that could be used
- events (changelog)

pr-project-kit.png:
- a kit is selected to see what it would add
- clicking apply increases the inventory of each part affected

pr-project-edit.png:
- editing the inventory
- you can see in the events the kit was added too

