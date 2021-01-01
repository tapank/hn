## hn
### terminal based hacker news (https://news.ycombinator.com) news for a lurker

To run: download and execute the binary: ./hn

What does it do?

1. Lists top 30 hn articles in the following format:

`[01/01 08:18 173 thesephist] 01. TabFS: Mount your Browser Tabs as a Filesystem (omar.website)`

The fields from left to write are: post date (mm/dd), post time (HH:mm), points, user id, post rank, title, domain.

2. Presents the following prompt and waits for input:

`enter your choice [<sno> | (m)ore | (t)op | (b)est | (n)ew | (q)uit | (r)efresh]:`

Enter an option and hit return. Options are:

- sno: the post rank. This opens the post in your Firefox browser
- m: next set of 30 posts in the currently chosen category (top posts is default)
- t: list top 30 posts
- b: list top 30 posts in best category
- n: list 30 newest posts
- q: quit the program
- r: refresh the list, while retaining current category and current offset of post rank

3. That's it. There is no other functionality in this tool (upvoting, posting, etc.), and is ideal for a lurker to keep an eye on hn all day long.

### Note:

- The release binary is for MaxOS (arch: amd64). You will need to download source and build your binary on a different platform.
- The posts open in Firefox. If you want to use a different browser, please update source and rebuild to suite your needs.
