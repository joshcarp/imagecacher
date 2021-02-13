# imagecacher

## Table of Contents
+ [About](#about)
+ [Installation](#installation)

+ [Usage](#usage)

## About <a name = "about"></a>
a cli tool for caching external image urls to file in html into data urls

#### Example

```html
<img src="http://www.plantuml.com/plantuml/png/SyfFKj2rKt3CoKnELR1Io4ZDoSa70000" style="max-width: 100%; height: auto;" alt="PlantUML diagram" onload="doneLoading()">
```
```bash
imagecacher -regex='<img.*?src="(?P<url>.*?)"' -output=outputdir -prefix=/foo/bar
```
This will output the files to `outputdir` and will prefix the replaced links with `/foo/bar` (useful for putting in static image directories)

## Installation <a name = "installation"></a>

#### go
`go get github.com/joshcarp/imagecacher`

#### docker

`docker pull joshcarp/imagecacher:latest`

## Usage <a name = "usage"></a>

```bash
imagecacher -input=<dir> -regex=<regex>
docker run -v $(pwd):/usr/app joshcarp/imagecacher -input=. -regex='<img.*?src="(?P<url>.*?)"'
```

where:
  - `input` is the directory to recursively replace images
  - `regex` is the regex to match to the links, where the `url` named capturing group is used to determine the image url:
    - `'!\[.*?\]\((?P<url>.*?)\)'` for markdown
    - `'<img.*?src="(?P<url>.*?)"'` for html

