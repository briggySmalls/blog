---
title: "CircleCI workspaces"
description: "Correctly attaching CircleCI workspaces"
date: 2020-02-09T23:30:00Z
showDate: true
tags: ["circle-ci"]
---

## The problem

When I created my CircleCI pipeline in my [first post]({{< relref "the-birth-of-a-blog.md" >}}) I came
unstuck passing the built site from one job to another as a _workspace_.

This manifested in two ways, where I'd made the same mistake in two different jobs.

### Deployment issues

The deploy job appeared to work correctly, but when I navigated to the site I was greeted with a 404
error:

![GitHub pages 404](/github-pages-404.jpg)

### HTMLProofer internal link checks

And when I ran the `htmlproofer` tool on my generated site I was greeted with a slew of errors:

<!-- markdownlint-disable fenced-code-language -->
```
- ./public/public/about/index.html
  *  internally linking to /about, which does not exist (line 0)
     <a href="/about">about</a>
```
<!-- markdownlint-enable fenced-code-language -->

In fact, _every_ internal link was flagged as having a target that didn't exist. Fishy...

## The solution

In both cases I had managed to point the respective tools to the wrong directory, based on a misunderstanding
on how workspaces are persisted.

My build job included the following configuration:

```yaml
  - persist_to_workspace:
      root: .
      paths: public
```

Which would store a workspace with the structure:

<!-- markdownlint-disable fenced-code-language -->
```
workspace
└── public
    ├── about
    ├── categories
    ...
```
<!-- markdownlint-enable fenced-code-language -->

And I'd then attach and point at them with:

```yaml
  - attach_workspace:
      at: ./public
  - run:
      name: test HTML files
      command: htmlproofer ./public --allow-hash-href --check-html
```

...Can you see my error?

After attaching the workspace I had a job containing:

<!-- markdownlint-disable fenced-code-language -->
```
public
└── public
    ├── about
    ├── categories
    ...
```
<!-- markdownlint-enable fenced-code-language -->

Meaning that in both jobs I was pointing to a directory `public` that contained one child directory
`public`. A sort-of off-by-one directory error.

Updating the config in both cases to mount the workspace to the working directory fixed both issues:

```diff
  - attach_workspace:
-     at: ./public
+     at: .
  - run:
      name: test HTML files
      command: htmlproofer ./public --allow-hash-href --check-html
```

## Conclusion

> Measure twice, cut once

As always, the issues that I seem to spend most time on are stupid errors due to not reading
documentation carefully enough!
