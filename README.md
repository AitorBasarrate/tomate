# tomate
Tomate is a CLI pomodoro app written in GO.
## Options
- Focus
    - `-fm {number}` to specify the number of minutes you want to focus
    - `-fs {number}` to specify the number of seconds you want to focus
- Rest
    - `-dm {number}` to specify the number of minutes you want to rest
    - `-ds {number}` to specify the number of seconds you want to rest
Both focus and rest can mix minutes and second to be more specific. Also you can use only seconds and it will be transformed into minutes and seconds.
- Repetition
    - `-r {number}` to specify the number of times you want the cicle to repeat
