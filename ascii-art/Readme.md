# ASCII ART #

Ascii-art is a program which consists in receiving a `string` as an argument and outputting the `string` in a graphic representation using ASCII. Time to write big.  
What we mean by a graphic representation using ASCII, is to write the string received using ASCII characters, as you can see in the example below:

```
@@@@@@BB@@@@``^^``^^``@@BB$$@@BB$$
@@%%$$$$^^^^WW&&8888&&^^""BBBB@@@@
@@@@@@""WW8888&&WW888888WW``@@@@$$
BB$$``&&&&WWWW8888&&&&8888&&``@@@@
$$``&&WW88&&88&&&&8888&&88WW88``$$
@@""&&&&&&&&88888888&&&&&&88&&``$$
``````^^``^^^^^^````""^^``^^``^^``
""WW^^@@@@^^``````^^BB@@^^``^^&&``
^^&&^^@@````^^``&&``@@````^^^^&&``
``WW&&^^""``^^WW&&&&""``^^^^&&88``
^^8888&&&&&&WW88&&88WW&&&&88&&WW``
@@``&&88888888WW&&WW88&&88WW88^^$$
@@""88&&&&&&&&888888&&``^^&&88``$$
@@@@^^&&&&&&""``^^^^^^8888&&^^@@@@
@@@@@@^^888888&&88&&&&MM88^^BB$$$$
@@@@@@BB````&&&&&&&&88""``BB@@BB$$
$$@@$$$$$$$$``````````@@$$@@$$$$$$
```

- This project should handle an input with numbers, letters, spaces, special characters and `\n`.  
- Take a look at the ASCII manual.

## Project Details ##

- The project is written in golang
- Some banner files with a specific graphical template representation using ASCII will be given. The files are formatted in a way that is not necessary to change them.  
  - [standard](standard.txt)
  - [shadow](shadow.txt)
  - [thinkertoy](thinkertoy.txt)

## Banner Format ##

- Each character has a height of 8 lines.
- Characters are separated by a new line \n.
- Here is an example of ' ', '!' and '"'(one dot represents one space) :

```file
......
......
......
......
......
......
......
......

._..
|.|.
|.|.
|.|.
|_|.
(_).
....
....

._._..
(.|.).
.V.V..
......
......
......
......
......

```

## Program Usage ##

Clone and `cd` into the repository  

```bash
git clone https://learn.zone01kisumu.ke/git/aochuka/ascii-art
cd ascii-art
```

while in the root directory, run the program.  
Examples:  
**Run**

```bash
bash $ go run . "" | cat -e
```

**output**

```bash
bash  $
```

**Run**

```bash
bash $ go run . "\n" | cat -e
```

**Output**

```bash  
$
```

**Run**

```bash
bash $ go run . "Hello\n" | cat -e
```

**Output**

```bash
 _    _          _   _          $
| |  | |        | | | |         $
| |__| |   ___  | | | |   ___   $
|  __  |  / _ \ | | | |  / _ \  $
| |  | | |  __/ | | | | | (_) | $
|_|  |_|  \___| |_| |_|  \___/  $
                                $
                                $
$ 
```

**Run**

```bash
bash $ go run . "hello" | cat -e
```

**Output**

```bash
 _              _   _          $
| |            | | | |         $
| |__     ___  | | | |   ___   $
|  _ \   / _ \ | | | |  / _ \  $
| | | | |  __/ | | | | | (_) | $
|_| |_|  \___| |_| |_|  \___/  $
                               $
                               $
```

**Run**

```bash
bash $ go run . "HeLlO" | cat -e
```

**OutPut**

```bash
 _    _          _        _    ____   $
| |  | |        | |      | |  / __ \  $
| |__| |   ___  | |      | | | |  | | $
|  __  |  / _ \ | |      | | | |  | | $
| |  | | |  __/ | |____  | | | |__| | $
|_|  |_|  \___| |______| |_|  \____/  $
                                      $
                                      $
```

```bash
bash $ go run . "Hello There" | cat -e
```

**OutPut**

```bash
 _    _          _   _                 _______   _                           $
| |  | |        | | | |               |__   __| | |                          $
| |__| |   ___  | | | |   ___            | |    | |__     ___   _ __    ___  $
|  __  |  / _ \ | | | |  / _ \           | |    |  _ \   / _ \ | '__|  / _ \ $
| |  | | |  __/ | | | | | (_) |          | |    | | | | |  __/ | |    |  __/ $
|_|  |_|  \___| |_| |_|  \___/           |_|    |_| |_|  \___| |_|     \___| $
                                                                             $
                                                                             $
```

```bash
bash $ go run . "1Hello 2There" | cat -e
```

**OutPut**

```bash
     _    _          _   _                         _______   _                           $
 _  | |  | |        | | | |                ____   |__   __| | |                          $
/ | | |__| |   ___  | | | |   ___         |___ \     | |    | |__     ___   _ __    ___  $
| | |  __  |  / _ \ | | | |  / _ \          __) |    | |    |  _ \   / _ \ | '__|  / _ \ $
| | | |  | | |  __/ | | | | | (_) |        / __/     | |    | | | | |  __/ | |    |  __/ $
|_| |_|  |_|  \___| |_| |_|  \___/        |_____|    |_|    |_| |_|  \___| |_|     \___| $
                                                                                         $
                                                                                         $
```

```bash
bash $ go run . "{Hello There}" | cat -e
```

**OutPut**

```bash
   __  _    _          _   _                 _______   _                           __    $
  / / | |  | |        | | | |               |__   __| | |                          \ \   $
 | |  | |__| |   ___  | | | |   ___            | |    | |__     ___   _ __    ___   | |  $
/ /   |  __  |  / _ \ | | | |  / _ \           | |    |  _ \   / _ \ | '__|  / _ \   \ \ $
\ \   | |  | | |  __/ | | | | | (_) |          | |    | | | | |  __/ | |    |  __/   / / $
 | |  |_|  |_|  \___| |_| |_|  \___/           |_|    |_| |_|  \___| |_|     \___|  | |  $
  \_\                                                                              /_/   $
                                                                                         $
```

```bash
bash $ go run . "Hello\nThere" | cat -e
```

**OutPut**

```bash
 _    _          _   _          $
| |  | |        | | | |         $
| |__| |   ___  | | | |   ___   $
|  __  |  / _ \ | | | |  / _ \  $
| |  | | |  __/ | | | | | (_) | $
|_|  |_|  \___| |_| |_|  \___/  $
                                $
                                $
 _______   _                           $
|__   __| | |                          $
   | |    | |__     ___   _ __    ___  $
   | |    |  _ \   / _ \ | '__|  / _ \ $
   | |    | | | | |  __/ | |    |  __/ $
   |_|    |_| |_|  \___| |_|     \___| $
                                       $
                                       $
```

```bash
bash $ go run . "Hello\n\nThere" | cat -e
```

**OutPut**

```bash
 _    _          _   _          $
| |  | |        | | | |         $
| |__| |   ___  | | | |   ___   $
|  __  |  / _ \ | | | |  / _ \  $
| |  | | |  __/ | | | | | (_) | $
|_|  |_|  \___| |_| |_|  \___/  $
                                $
                                $
$
 _______   _                           $
|__   __| | |                          $
   | |    | |__     ___   _ __    ___  $
   | |    |  _ \   / _ \ | '__|  / _ \ $
   | |    | | | | |  __/ | |    |  __/ $
   |_|    |_| |_|  \___| |_|     \___| $
                                       $
                                       $
```

## Output Function ##

**Purpose**:

Prints words as ASCII art using a map of characters.

How it works, step by step:

- Takes a list of words (input) and a map of letters to 8-line ASCII art (charmap).

- For each word in the list:

- If the word is empty, just prints a blank line. Otherwise, for each of the 8 lines:

Build a line by putting together the corresponding line of each letter in the word.
Add that line to the result.

Prints all the lines, so the word appears in ASCII art.

## TestMapToASCII tests the MapToASCII function ##  

The purpose of MapToASCII is to convert a flat list of ASCII "lines" into a map where each printable ASCII character (' ' to '~') maps to 8 lines of ASCII art.

## Step 1 ##

- Calculate the number of printable ASCII characters  

- ' ' (space) has ASCII code 32  
- '~' (tilde) has ASCII code 126  
- totalChars = 126 - 32 + 1 = 95 printable characters  
- totalChars := int('~' - ' ' + 1)  

## Step 2 ##

- Calculate the total number of lines needed in our input slice  
  Each character in the ASCII art uses 9 lines,but MapToASCII only stores 8 lines per character in the map.  
  The first line (index 0) is skipped (usually header or padding), so we add 1.  
  totalLines := 1 + totalChars*9  

## Step 3 ##  

- Create a slice to simulate the ASCII art input lines  
  Fill it with placeholder strings ("line") so we can test the mapping algorithm.  

## Step 4 ##  

- Call the function we want to test  
  MapToASCII will take the flat 'lines' slice and convert it into a map,where each ASCII character maps to a slice of 8 lines.  

## Step 5 ##  

- Verify that the number of characters in the resulting map is correct  
  There should be exactly totalChars entries (from ' ' to '~')  

## Step 6 ##

- Verify that each character in the map has exactly 8 lines  
  This ensures MapToASCII correctly slices 8 lines per character  
