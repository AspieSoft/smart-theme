# Smart Theme

[![donation link](https://img.shields.io/badge/buy%20me%20a%20coffee-paypal-blue)](https://paypal.me/shaynejrtaylor?country.x=US&locale.x=en_US)

A simple theme that adjusts easily to a users preferences.

## Notice: This Theme Is Currently In Alpha

## Installation

```html
<!-- Config -->
<!-- Note: You should replace these stylesheets with a local copy that you can modify to your preferences -->
<link rel="stylesheet" href="https://cdn.jsdelivr.net/gh/AspieSoft/smart-theme@0.0.1/theme/fonts.swap.css"/>
<link rel="stylesheet" href="https://cdn.jsdelivr.net/gh/AspieSoft/smart-theme@0.0.1/theme/config.css"/>

<!-- Theme Style -->
<link rel="stylesheet" href="https://cdn.jsdelivr.net/gh/AspieSoft/smart-theme@0.0.1/theme/style.norm.min.css"/>
<!-- or for a lighter version without normalize.css -->
<link rel="stylesheet" href="https://cdn.jsdelivr.net/gh/AspieSoft/smart-theme@0.0.1/theme/style.min.css"/>

<!-- Theme Script -->
<script src="https://cdn.jsdelivr.net/gh/AspieSoft/smart-theme@0.0.1/theme/script.min.js" defer></script>
```

## Setup

```html
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8"/>
  <meta name="viewport" content="width=device-width, height=device-height, initial-scale=1.0, minimum-scale=1.0"/>
  
  <link rel="preconnect" href="https://fonts.googleapis.com" crossorigin/>
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin/>

  <meta name="description" content="Your Description"/>
  <title>Your Title</title>

  <link rel="stylesheet" href="/theme/fonts.swap.css"/>
  <link rel="stylesheet" href="/theme/config.css"/>

  <link rel="stylesheet" href="/theme/style.norm.min.css"/>

  <script src="/theme/script.min.js" defer></script>

  <!--? optional: tailwind and htmx recommended -->
  <script src="https://cdn.tailwindcss.com"></script>
  <script src="https://unpkg.com/htmx.org@1.9.10"></script>
</head>
<body>
  <header>
    <div class="header-top">
      <a href="/" class="site-title">Your Logo</a>
      <!-- todo: add nav menu support here -->
    </div>
    <div class="header-image">
      <!--? optional element: nav menu -->
      <nav>
        <ul>
          <li><a href="/">Home</a></li>
          <li><a href="/about">About</a></li>
          <li><a href="/contact">Contact Us</a></li>
        </ul>
      </nav>
      <div class="content">
        <h1>Welcome To Your Title</h1>
      </div>
    </div>
  </header>
  <div id="page-content">
    <!--? optional class: page-bg -->
    <main class="content-grid page-bg">
      <h1>Hello, World!</h1>
      <p>Lorem ipsum dolor sit amet consectetur adipisicing elit. Aspernatur deleniti quidem alias delectus eius, ipsam molestiae? Dignissimos sequi eaque hic accusantium, molestias suscipit molestiae pariatur quisquam incidunt est nihil eveniet.</p>
      <p class="breakout">Lorem ipsum dolor sit amet consectetur adipisicing elit. Molestiae vero, cupiditate nulla esse nihil impedit voluptatum architecto autem incidunt animi? Neque, recusandae ipsum voluptatibus laudantium quod at ullam similique sequi.</p>
      <p>Lorem ipsum dolor sit amet consectetur adipisicing elit. Esse, exercitationem sunt nihil eum dolorum dolor culpa! Accusamus iusto ex nisi et doloremque expedita, voluptatem perferendis soluta iure omnis rerum sit.</p>
      <img src="/assets/background/falcon-desert.webp" alt="millennium falcon"/>
      <p>Lorem ipsum dolor sit amet consectetur adipisicing elit. Aspernatur deleniti quidem alias delectus eius, ipsam molestiae? Dignissimos sequi eaque hic accusantium, molestias suscipit molestiae pariatur quisquam incidunt est nihil eveniet.</p>
      <p class="breakout">Lorem ipsum dolor sit amet consectetur adipisicing elit. Molestiae vero, cupiditate nulla esse nihil impedit voluptatum architecto autem incidunt animi? Neque, recusandae ipsum voluptatibus laudantium quod at ullam similique sequi.</p>
      <p>Lorem ipsum dolor sit amet consectetur adipisicing elit. Esse, exercitationem sunt nihil eum dolorum dolor culpa! Accusamus iusto ex nisi et doloremque expedita, voluptatem perferendis soluta iure omnis rerum sit.</p>
    </main>
  </div>
  <footer>
    <!--? optional element: footer content -->
    <div class="content">
      <h1>Footer Content</h1>
      <p>Lorem ipsum dolor sit amet, consectetur adipisicing elit. Modi, libero animi. Doloremque minus quidem quasi vitae voluptates, porro officiis praesentium cupiditate omnis eligendi assumenda neque non? Impedit eligendi ab aut! Lorem ipsum dolor sit amet consectetur adipisicing elit. Corporis, sapiente animi enim sed tempora itaque eveniet excepturi nisi nobis similique ad voluptatum tempore minima sit fugiat reiciendis sunt quasi ab.</p>
    </div>
    <div class="copyright">
      <!--? you can replace this with your copyright text -->
      Powered by <a href="https://aspiesoft.com">AspieSoft</a>
    </div>
  </footer>
</body>
</html>
```
