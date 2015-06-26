# TheMatrix style Digital Rain
Recreate The Matrix digital rain effect in Golang.

This was my first musings in Google's go language. It turned out well, it took me about a week off-and-on to get through this for a demo who I am going to dedicate this to.

# Generally, how it works

The script leverages the "termbox-go" package to control the terminal/console to render the different characters.  A static utf-8 string of japanese characters is used to randomly fill the screen with characters.  The app leverages goroutines to "produce" and "consume" cells.  The producers generate the cells and the consumer renders the cells in the channel/queue.

# Dedication

This little demo app is dedicated to Daniel Philips, I hope it inspires him to be a "coder". The USA could use more native engineers.
