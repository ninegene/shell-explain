---
schema_version: 1
id: fetch-from-the-web
title: Ask a Web Service
subtitle: Send a request with curl and read the answer
summary: Ask a web service a question with curl, read the answer and its status code, and tell a failed request from a successful one. Practice on local files first, then on the public GitHub API.
series: getting-things-done
order: 5
level: beginner
reading_minutes: 8
updated: 2026-11-10
tags: [curl, http, api, network]
related_commands: [curl, jq, cat]
related_guides: [input-output-errors, pipes-and-redirects, exit-codes]
references:
  - https://curl.se/docs/manpage.html
  - https://docs.github.com/en/rest/using-the-rest-api/getting-started-with-the-rest-api
  - https://docs.github.com/en/rest/using-the-rest-api/rate-limits-for-the-rest-api
  - https://developer.mozilla.org/en-US/docs/Web/HTTP/Status
---
## Goal
By the end you can ask a web service for information from the terminal, see whether it said yes or no, keep the answer in a file, and pick one field out of it. This guide is about the layer where your shell talks to another computer.

## The idea in plain words
A **web service** is a program on another computer that answers questions sent over the internet. A question is called a **request** and the reply is a **response**. Both use a set of rules called **HTTP**.

A request has an address, called a **URL**, such as `https://api.github.com/zen`. The response has two parts that matter to you:

1. A **status code**, a three-digit number that says how it went. `200` means success. `404` means the thing you asked for was not found. Codes that start with `4` mean your request was wrong, and codes that start with `5` mean the service had a problem.
2. A **body**, the content you asked for. Web services often send it as **JSON**, a text format made of names and values in curly braces.

A program that lets you send a request from the command line is `curl`. It prints the body on standard output, so everything you learned about pipes and redirects applies to it.

A web service that is meant for programs, not people with a browser, is called an **API**. The GitHub API is a good one to practice on: it is public, it is well documented, and reading from it needs no account. Without signing in you get about 60 requests per hour, which is plenty here.

## A picture
Read the picture from left to right. Your command goes out as a request, and the answer comes back with a status code and a body.

```diagram
   your shell                              the web service
   ----------                              ---------------
   curl URL        -- request ------>      reads the URL
                                           works out an answer
   status + body   <-- response -----      sends it back

   the body goes to stdout, so you can pipe it or save it:
   curl URL | jq -r .name        curl -o answer.json URL
```

Diagram: on the left is your shell and on the right is the web service. An arrow labelled request goes from the curl command to the service. An arrow labelled response comes back carrying a status code and a body. Below, two lines remind you that the body goes to standard output, so it can be piped into jq or saved into a file with curl dash o.

## Another way to see it
Think of ordering at a counter. The URL is what you ask for. The status code is the first thing the person says: "here you go" (200), "we do not have that" (404) or "the kitchen is broken" (500). The body is the food itself.

The table lists the status codes you will meet most.

| Code | Meaning | What to do |
| --- | --- | --- |
| 200 | It worked | Read the body |
| 404 | Not found | Check the address |
| 403 | Refused, often a limit | Wait, or sign in |
| 500 | Service problem | Try again later |

## Walk through it
Here is what happens, step by step, when you run `curl -s https://api.github.com/repos/octocat/Hello-World | jq -r .full_name`.

1. `curl` opens a connection to `api.github.com` and sends a request for the path `/repos/octocat/Hello-World`.
2. The `-s` option means silent. Without it, curl prints a progress meter on standard error, which gets in the way when you read the answer.
3. The service finds the repository and sends back a status of `200` and a body of JSON, about a hundred lines long.
4. `curl` writes the body on standard output. The pipe `|` hands it to `jq`.
5. `jq -r .full_name` reads the JSON and prints only the field `full_name`. The `-r` option prints it as plain text, without quotation marks. You see `octocat/Hello-World`.

If the repository did not exist, step 3 would send the status `404` and a short JSON body that says `Not Found`. Note that curl itself counts this as a success, because the conversation worked: the service answered. To see the status code, you have to ask for it, as the next section shows.

## Try it
You should not need the internet to learn the habits, so these commands ask for a file on your own computer. A URL that starts with `file://` tells curl to read a local file instead of calling a service. The rest is the same as for a web address. The block that follows uses `$PWD`, the folder you are in, to build the address.

Make a small file, then fetch it. The answer is the body, printed on standard output:

```sh
printf 'hello\n' > f.txt
curl -s "file://$PWD/f.txt"
```

You should see this:

```text
hello
```

Save the answer into a file with `-o`, and print the number of bytes that arrived with `-w`. The `-w` option writes a message after the transfer, and `%{size_download}` stands for the size:

```sh
printf 'hello\n' > f.txt
curl -s -o saved.txt -w '%{size_download}\n' "file://$PWD/f.txt"
cat saved.txt
```

You should see this:

```text
6
hello
```

Ask for something that is not there, and look at the exit status. Curl uses exit status 37 when it cannot read a local file. The `-s` option hides its error message, so the status is the only sign:

```sh
curl -s "file://$PWD/missing.txt"
echo "exit status $?"
```

You should see this:

```text
exit status 37
```

### On the real internet
The commands below call the public GitHub API. The tests for this guide never send them, because a network can be down and a limit can be reached, but they were run by hand on 2026-11-10. Each one only reads.

Ask for a short saying. The answer changes from one request to the next:

```bash
curl -s https://api.github.com/zen
```

Read one field out of a JSON answer:

```bash
curl -s https://api.github.com/repos/octocat/Hello-World | jq -r .full_name
```

You should see this:

```text
octocat/Hello-World
```

Show only the status line of the response. `-I` asks for the headers, not the body. The exact words depend on the HTTP version, so it may say `HTTP/1.1 200 OK` instead:

```bash
curl -sI https://api.github.com/repos/octocat/Hello-World | head -n 1
```

You should see something like this:

```text
HTTP/2 200
```

Ask for a repository that does not exist, throw the body away with `-o /dev/null`, and print only the status code:

```bash
curl -s -o /dev/null -w '%{http_code}\n' \
  https://api.github.com/repos/octocat/nope
```

You should see this:

```text
404
```

Look at the body of that same failed request. The service explains the problem in JSON:

```bash
curl -s https://api.github.com/repos/octocat/nope | jq -r .message
```

You should see this:

```text
Not Found
```

## Common confusions
- **A finished curl means the request worked.** Curl exits with status 0 whenever the service answers, even with `404`. Print the status with `-w '%{http_code}\n'` when it matters.
- **The `-s` option makes curl quiet about errors too.** It hides the progress meter and the error messages. When a command fails and you cannot see why, run it again without `-s`.
- **A web address in a script is always safe to run.** Reading is harmless, but a request can also send data or change something on the service. Run only requests whose purpose you understand, and never paste a command that contains someone else's secret key.
- **An API key is just a password for programs.** Treat it as one. Do not put it in a guide, a shared note or a command you share.
- **The pipe waits for the answer.** `curl ... | jq ...` shows nothing until the whole answer has arrived. A slow network looks like a hang.
- **A limit is a bug.** Public services allow only so many requests per hour. If you see a `403` after many tries, wait before you try again.

## Check yourself
1. What are the two parts of a response that matter most? Answer: the status code and the body.
2. What does the status code `404` mean? Answer: the thing you asked for was not found.
3. Why does `curl` exit with status 0 on a `404`? Answer: the conversation worked and the service answered, so curl counts it as a success. Ask for the status with `-w '%{http_code}\n'`.
4. How do you keep an answer in a file? Answer: `curl -s -o answer.json URL`, or redirect with `>`.

## Where next
This ends the goal-first guides: you can find a word, count and sort, tidy files, make a backup and ask a web service, and each one uses ideas from the earlier series.

For the details of each command, see `explain curl`, `explain jq` and `explain cat`.
