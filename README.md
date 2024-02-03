![Build and Lint CLI](https://github.com/runvelocity/cli/actions/workflows/build-cli.yaml/badge.svg)

<div align="center">
  <!-- <a href="https://github.com/othneildrew/Best-README-Template">
    <img src="images/logo.png" alt="Logo" width="80" height="80">
  </a> -->

  <h3 align="center">CLI</h3>

  <p align="center">
    The Velocity CLI allows you to create and run serverless Velocity functions from your terminal.
    <br />
    <a href="https://docs.runvelocity.dev"><strong>Explore the docs »</strong></a>
    <br />
  </p>
</div>

<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about">About</a>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#installation">Installation</a></li>
      </ul>
    </li>
    <li><a href="#usage">Usage</a></li>
    <li><a href="#roadmap">Roadmap</a></li>
    <li><a href="#contact">Contact</a></li>
  </ol>
</details>

## About

The CLI is written in Go using the Cobra framework and Bubbletea as a TUI. It allows you to create, delete and invoke Velocity functions.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- GETTING STARTED -->
## Getting Started

This section outlines how to install the CLI to get started creating functions.

### Installation

Download the precompiled binary for your operating system on the [release page](https://github.com/runvelocity/cli/releases)

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- USAGE EXAMPLES -->
## Usage

### Creating a function

Create a file called function.js and input the following code

```javascript
module.exports = (args) => {
    return {
        hello: "world",
    }
}
```

Next, you have to create a zip file containing this code. If you are on windows with powershell, use the following command

```shell
Compress-Archive -LiteralPath .\function.js -DestinationPath  code.zip
```

On Linux/Mac, use the following command

```shell
zip code.zip function.js
```

Next, create a function by running the following command

```
.\cli.exe create --name demo-func --file-path code.zip --handler function
```

### Invoking a function

To invoke a function, run the following command

```shell
.\cli.exe invoke --name demo-func
```

<!-- ROADMAP -->
## Roadmap

- [ ] Add support for updating functions
- [ ] Pass in parameters from the command line

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- CONTACT -->
## Contact

Utibeabasi Umanah - [@utibeumanah_](https://twitter.com/utibeumanah_) - utibeabasiumanah6@gmail.com

Project Link: [https://github.com/runvelocity](https://github.com/runvelocity)

<p align="right">(<a href="#readme-top">back to top</a>)</p>