# IMGDD

## Introduction

## Deployment
### Config

## FAQ:
<details>
  <summary>
    <b>How do I change the footer text:</b>
  </summary>
IMGDD provides a slot-filler based system to allow users to change some text on the web client. The footer text uses this system and has the slot-id `footer-content`. You can register a plugin to change it to the content you want.

For example you can use the following snipet:
```javascript
  registerPlugin({textSlots: {
    "footer-content": function(React) {
      return React.createElement("div", {className: "text-center text-sm text-neutral-500 dark:text-neutral-400"}, "Powered by ", React.createElement("a", {href: "https://mywebsite.home.arpa", className: "text-neutral-500 dark:text-neutral-400"}, "My website"));
    }
  }});
```
</details>
