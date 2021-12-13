# SimplyRunes

CLI Program for setting League Of Legends runes automatically. Written in Golang.

I made this, because I was tired of other apps like Porofessor or Blitz using way more RAM then they need to, and showing multiple ads.

# Usage

Simply just call the program from your terminal

```
> ./simplyrunes.exe
```

after that you should get messege

```
> Waiting for champion select.
```

then if you lock any champion it will set your runes based on most common ones (Platinium+) and after champ select is over, you should get most common item build on that champion

# Known issues

- SimplyRunes will delete any other rune page.
- Simplyrunes will set your runes based on most popular role for that champion (ex. if you wanted to go Lucian TOP it wil set your runes to Lucian ADC)
- Program exists after the game has started
- It shows summoner spells after champ select is over

Currently im working on a complete revamp on the entire app, but if you think that SimplyRunes will fit your needs right now, feel free to use it.
