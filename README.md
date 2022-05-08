# Devoxx UK - Question 1

This is question 1 of a game intended to demonstrate serverless and event-driven features on Kubernetes, using Knative and CloudEvents.

* For more information on the game, visit this [page](https://github.com/salaboy/from-monolith-to-k8s/tree/main/game).
* For instructions on deploying the entire system, visit this [page](https://github.com/ThomasVitale/eventing-game).

Question 1 is a Go project relying on the CloudEvents GO SDK. The project has been initialized using
the Knative [func](https://github.com/knative-sandbox/kn-plugin-func) plugin.

## Usage

```shell
$ http <url> player="jon-snow" sessionId="game-blahblah" optionA=true optionB=false optionC=false optionD=false remainingTime=13

HTTP/1.1 200 OK
Content-Length: 98
Content-Type: application/json
accept-encoding: gzip, deflate
connection: keep-alive
user-agent: HTTPie/3.1.0

{
    "player": "jon-snow"
    "level": "devoxxuk-question-1",
    "levelScore": 18,
    "sessionId": "game-blahblah",
    "gameTime": "2022-04-19T11:40:46.04108"
}
```
