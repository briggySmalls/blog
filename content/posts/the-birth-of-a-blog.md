---
title: "The birth of a blog"
description: "How I built and deployed this website, including stack-selection and pipeline pitfalls!"
date: 2020-02-03T14:59:18Z
showDate: true
tags: ["hugo","web"]
---

## Introduction

This post is about how I set up this site!

## Framework

So I wanted a static site generator. You know, the trendy, 'lightning-fast', modern-web-has-gone-full-circle
technology. Because I just want to write some text that people can read, and I can execute a build step
without hurting myself.

> Perfect is the enemy of good

Knowing I'm the kind of person who often lets perfect get the better of me, I didn't want to get stuck
at step one _this_ time.

I took the top three static site frameworks from [stackshare][stackshare-generators]: Jekyll, Gatsby,
Hugo. And whittled them down through opinionated first-impressions and prejudices:

[stackshare-generators]: https://stackshare.io/static-site-generators

1. [Gatsby](https://www.gatsbyjs.org/) is react-based and seemed to be peddling [JAMstack](https://jamstack.org/)
'webapp'. This sounded a lot like 'JavaScript bloat' for a blog use-case.
2. I'd used [Jekyll](https://jekyllrb.com/) once before, and for my chosen theme I had to commit some
file dump and modify the theme files directly. Maybe this isn't always the case, but my experience
scarred me. No theme updates?
3. [Hugo](https://gohugo.io/) is like Jekyll, but _new_, and Go 🥰.

## Setup

Being on mac, installing Hugo with [homebrew](https://brew.sh/) is breezy:

```bash
brew install hugo
```

And following their [quick-start](https://gohugo.io/getting-started/quick-start/) was straightforward
enough.

After my Jekyll experience I was delighted to find themes are tracked as git submodules. Themes are swappable
and updates are easy! 🙌

## Deployment

### Netlify misfire

I had assumed I'd host for free with [netlify](https://www.netlify.com/), another trendy platform where
Hugo is a first-class citizen.

But after linking to my GitHub repo and configuring the build command I was greeted with:

<!-- markdownlint-disable fenced-code-language line-length -->
```
Transformation failed: POSTCSS: failed to transform “css/main.css” (text/css): PostCSS not found; install with “npm install postcss-cli”. See https://gohugo.io/hugo-pipes/postcss/ 2
```
<!-- markdownlint-enable fenced-code-language line-length -->

It turns out my chosen theme ([Call me Sam](https://themes.gohugo.io/hugo-theme-sam/) - what a narcissist!)
depended on the [PostCSS Hugo Pipe](https://gohugo.io/hugo-pipes/postcss/), which meant I needed
some npm packages installed globally.

Whilst I could fix my local build like this easily enough, it was less obvious where to do this with
Netlify. Was node preinstalled? Can I execute multiple build commands?

<!-- markdownlint-disable no-trailing-punctuation -->
### If you want a job done right...
<!-- markdownlint-enable no-trailing-punctuation -->

I decided to take control of the build pipeline myself. Who knows what future dependencies I might want?
If I had to learn some platform-specific configuration, I'd rather use one where I knew I couldn't
hit a brick wall.

So I built the site in a [CircleCI](https://circleci.com/) pipeline and deployed to GitHub pages. I
like CircleCI because you can use docker images for each job, unlike e.g. Travis, which can save some
time installing dependencies.

After following some [instructions][circleci-tutorial-github], the job that pushes to GitHub pages announced:

[circleci-tutorial-github]: https://circleci.com/blog/deploying-documentation-to-github-pages-with-continuous-integration/

<!-- markdownlint-disable fenced-code-language -->
```
ERROR: The key you are authenticating with has been marked as read only.
```
<!-- markdownlint-enable fenced-code-language -->

Clearly I'd lost patience reading, because the 'Provisioning a deploy key' section contained it all.
Once I created a new key and gave the public part to GitHub and the private to CircleCI I was in business!
Seriously, take your time with that article, it has everything you need.

You can see the pipeline config [here][blog-circleci-config] in all its glory.

[blog-circleci-config]: https://github.com/briggySmalls/blog/blob/a5c9745a1134491a0369aee7bf43883e8b045b3d/.circleci/config.yml

## Conclusion

And there we have it! A Hugo site, built on CircleCI, deployed on GitHub pages. The internet has yet
another blog.
