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
imagecacher -input=. -regex='<img.*?src="(?P<url>.*?)"'
```
```html
<img src="data:image/png;base64,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" style="max-width: 100%; height: auto;" alt="PlantUML diagram" onload="doneLoading()">
```

## Installation <a name = "installation"></a>

#### go
`go get github.com/joshcarp/imagecacher`

#### docker

`docker pull joshcarp/imagecacher:latest`

## Usage <a name = "usage"></a>

`imagecacher -input=<dir> -regex=<regex>`

where:
  - `input` is the directory to recursively replace images
  - `regex` is the regex to match to the links, where the `url` named capturing group is used to determine the image url:
    - `'!\[.*?\]\((?P<url>.*?)\)'` for markdown
    - `'<img.*?src="(?P<url>.*?)"'` for html

