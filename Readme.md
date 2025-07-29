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

<img width="1917" height="859" alt="Screenshot (34)" src="https://github.com/user-attachments/assets/e2e3a256-47e8-4acd-85c2-a66f9b260898" />
<img width="1915" height="844" alt="Screenshot 2025-07-29 214839" src="https://github.com/user-attachments/assets/2294a24a-6cd1-4599-bfbd-79b0ed752207" />


## Setup

```bash
# Clone the repository
git clone https://github.com/amelia2802/book-suggestion-agent.git
cd book-suggestion-agent

# Install dependencies
npm install
npm install -g @hypermode/hyp-cli 
# Login to Hypermode
# Make sure you have a Hypermode account and OpenAI API key, create a .env file in root folder and place your api key there
 ```MODUS_OPENAI_API_KEY= your_openai_api_key```
hyp dev
hyp login
```

## Usage

To run the book suggestion agent, follow these steps:
```bash
# Start the Hypermode development environment  
modus dev
# Go to the following address and check by entering the genre
http://localhost:8686/explorer/
```

## Requirements

- Node.js - v22 or higher
- go version >= 1.20
- TinyGo
- Hypermode CLI
- Hypermode account
- OpenAI API key
- Modus SDK
- Node.js


---

*Created as part of the Hypermode Workshop by Torc.*
