# Book Suggestion Agent

This is a homework project for the [**Hypermode Workshop**](https://hypermode.com/) organized by [**Torc**](https://platform.torc.dev/#/r/ZILxKHb0/cp).

![banner](https://pbs.twimg.com/media/GvApKqNWkAAiXbl?format=jpg&name=900x900)

## Table of Contents

- [Overview](#overview)
- [Setup](#setup)
- [Usage](#usage)
- [Requirements](#requirements)
- [Contributing](#contributing)
- [License](#license)

## Overview

This project is a book suggestion agent that utilizes the Hypermode platform to provide personalized book recommendations based on user preferences. It integrates with the OpenAI API to analyze user input and suggest books accordingly.

## Setup

```bash
# Clone the repository
git clone https://github.com/amelia2802/book-suggestion-agent.git
cd book-suggestion-agent

# Install dependencies
go mod tidy
npm install
# Install Hypermode CLI
scoop install tinygo binaryen
# Install Hypermode CLI and Modus SDK
npm install -g @hypermode/modus-sdk 
npm install -g @hypermode/hyp-cli 
go get github.com/hypermodeinc/modus/sdk/go/pkg/models/openai@v0.18.0
# Login to Hypermode
# Make sure you have a Hypermode account and OpenAI API key 
hyp dev
hyp login
```

## Usage

To run the book suggestion agent, follow these steps:
```bash
# Start the Hypermode development environment  
modus dev
```

## Requirements

- List any prerequisites or dependencies.
- go version >= 1.20
- TinyGo
- Hypermode CLI
- Hypermode account
- OpenAI API key
- Modus SDK
- Node.js


---

*Created as part of the Hypermode Workshop by Torc.*