# cmus-lyric
[cmus](https://cmus.github.io/) lyrics [viewer](https://asciinema.org/a/69vGAibC1uvkDBR7WuSecbydd)


# Like
![](https://i.imgur.com/WNxuUZ7.png)

With tmux 
![](https://i.imgur.com/wL3tPZa.png)

With comments
![](https://i.imgur.com/UUUf9lZ.png)

help
```bash
usage:

 q or <C-c>: quit
 m         : view comments
 y         : view lyrics
 ?         : help


```


# Install
Linux

`curl -L https://github.com/rockagen/cmus-lyric/raw/master/lyrics -o lyrics`

MacOS

`curl -L https://github.com/rockagen/cmus-lyric/raw/master/lyrics_osx -o lyrics`


`chmod u+x lyrics`


# How
Check cmus current file exist lyric,fetch from music.163.com if not found

# Requirements
`go` compile 

`termui` term ui


# Build
Install lyrics
```bash
make install
```


# Run
`./lyrics`

type `q` to quit


happy enjoy!
