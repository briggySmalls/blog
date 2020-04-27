---
title: "No Jekyll"
description: "Wrestling with GitHub pages, Jekyll, and folders with underscores."
date: 2020-04-14T21:35:06+01:00
showDate: true
tags: ["hugo", "github-pages"]
---

I had a bit of a head-scratcher yesterday, as I was writing the [Concepts]({{< relref "concepts.md" >}})
post and found that my slickly auto-generated architecture diagrams were missing from my posts. Navigating
to the image source URLs was giving me the big fat 404 on GitHub pages:

![GitHub pages 404](/github-pages-404.jpg)

For sanity, I checked out the deploy branch, `gh-pages`, and sure enough my diagrams were present,
nestled in a folder called `_gen`.

I was immediately suspicious about the leading underscore, and a [quick search][jekyll-underscore] of
"GitHub pages underscore" reminded me that GitHub processes deployments with Jekyll unless instructed
not to. And what does Jekyll do? Removes files with a leading underscore. The smoking gun!

[jekyll-underscore]: https://github.community/t5/GitHub-Pages/Cannot-access-underscore-folder-file-with-nojekyll/td-p/36087

Straightforward stuff, I thought. I've been here before, just whack a [.nojekyll file][no-jekyll] in
the root of the pages branch. So I updated my CircleCI deploy job:

[no-jekyll]: https://github.blog/2009-12-29-bypassing-jekyll-on-github-pages/

```diff
+ - run:
+     name: Disable jekyll processing
+     command: touch ./public/.nojekyll
  - run:
      name: Deploy to gh-pages branch
      command: gh-pages --dist ./public --message "[ci skip] Deploy updates"
```

But hold on. Still a 404? What in God's name have I done to deserve this persecution? Is it there? I
check out the branch, _no_, for crying out loud where's it gone? Does `touch` not work on this damn
docker image?

After losing some time to unhelpful stress[^1] I eventually saw that _now_ the blame probably lay with
the subsequent `gh-pages` deploy step. Having a wee butchers at the tool's options was telling:

[^1]: Note to self, it's sometimes very important to just take some time out and go make a cup of tea

<!-- markdownlint-disable fenced-code-language -->
```
Usage: gh-pages [options]

Options:
  -V, --version            output the version number
...
  -t, --dotfiles           Include dotfiles
...
  -h, --help               output usage information
```
<!-- markdownlint-enable fenced-code-language -->

...Well hello there `--dotfiles` you agent of chaos. I've got you now!

```diff
  - run:
      name: Deploy to gh-pages branch
-     command: gh-pages --dist ./public --message "[ci skip] Deploy updates"
+     command: gh-pages --dotfiles --dist ./public --message "[ci skip] Deploy updates"
```

Praise be. My nightmare ends.
